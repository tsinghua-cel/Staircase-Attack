package server

import (
	"context"
	"errors"
	"fmt"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	ethtype "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/groupcache/lru"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/beaconapi"
	"github.com/tsinghua-cel/attacker-service/config"
	"github.com/tsinghua-cel/attacker-service/dbmodel"
	"github.com/tsinghua-cel/attacker-service/openapi"
	"github.com/tsinghua-cel/attacker-service/plugins"
	"github.com/tsinghua-cel/attacker-service/rpc"
	"github.com/tsinghua-cel/attacker-service/server/apis"
	"github.com/tsinghua-cel/attacker-service/strategy"
	"github.com/tsinghua-cel/attacker-service/strategy/slotstrategy"
	"github.com/tsinghua-cel/attacker-service/types"
	"math/big"
	"strconv"
	"time"
)

type Server struct {
	config       *config.Config
	rpcAPIs      []rpc.API   // List of APIs currently provided by the node
	http         *httpServer //
	strategy     *types.Strategy
	internal     []slotstrategy.InternalSlotStrategy
	execClient   *ethclient.Client
	beaconClient *beaconapi.BeaconGwClient

	validatorSetInfo *types.ValidatorDataSet
	openApi          *openapi.OpenAPI
	cache            *lru.Cache
}

func NewServer(conf *config.Config, plugin plugins.AttackerPlugin) *Server {
	s := &Server{}
	s.cache = lru.New(10000)
	s.config = conf
	s.rpcAPIs = apis.GetAPIs(s, plugin)
	client, err := ethclient.Dial(conf.ExecuteRpc)
	if err != nil {
		panic(fmt.Sprintf("dial execute failed with err:%v", err))
	}
	s.execClient = client
	s.beaconClient = beaconapi.NewBeaconGwClient(conf.BeaconRpc)
	s.http = newHTTPServer(log.WithField("module", "server"), rpc.DefaultHTTPTimeouts)
	s.strategy = strategy.ParseStrategy(s, conf.Strategy)
	s.validatorSetInfo = types.NewValidatorSet()
	s.openApi = openapi.NewOpenAPI(s, conf)
	return s
}

// startRPC is a helper method to configure all the various RPC endpoints during node
// startup. It's not meant to be called at any time afterwards as it makes certain
// assumptions about the state of the node.
func (n *Server) startRPC() error {
	// Filter out personal api
	var (
		servers []*httpServer
	)

	rpcConfig := rpcEndpointConfig{
		batchItemLimit:         config.APIBatchItemLimit,
		batchResponseSizeLimit: config.APIBatchResponseSizeLimit,
	}

	initHttp := func(server *httpServer, port int) error {
		if err := server.setListenAddr(n.config.HttpHost, port); err != nil {
			return err
		}
		if err := server.enableRPC(n.rpcAPIs, httpConfig{
			CorsAllowedOrigins: config.DefaultCors,
			Vhosts:             config.DefaultVhosts,
			Modules:            config.DefaultModules,
			prefix:             config.DefaultPrefix,
			rpcEndpointConfig:  rpcConfig,
		}); err != nil {
			return err
		}
		servers = append(servers, server)
		return nil
	}

	// Set up HTTP.
	// Configure legacy unauthenticated HTTP.
	if err := initHttp(n.http, n.config.HttpPort); err != nil {
		return err
	}

	// Start the servers
	for _, server := range servers {
		if err := server.start(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) monitorEvent() {
	ticker := time.NewTicker(time.Minute * 2)
	defer ticker.Stop()

	handler := func(ch chan *apiv1.ChainReorgEvent) {
		for {
			select {
			case reorg := <-ch:
				log.WithFields(log.Fields{
					"slot": reorg.Slot,
				}).Info("reorg event")
				ev := types.ReorgEvent{
					Epoch:        int64(reorg.Epoch),
					Slot:         int64(reorg.Slot),
					Depth:        int64(reorg.Depth),
					OldHeadState: reorg.OldHeadState.String(),
					NewHeadState: reorg.NewHeadState.String(),
				}
				if oldHeader, err := s.beaconClient.GetBlockHeaderById(reorg.OldHeadBlock.String()); err == nil {
					ev.OldBlockSlot = int64(oldHeader.Header.Message.Slot)
					ev.OldBlockProposerIndex = int64(oldHeader.Header.Message.ProposerIndex)
				}
				if newHeader, err := s.beaconClient.GetBlockHeaderById(reorg.NewHeadBlock.String()); err == nil {
					ev.NewBlockSlot = int64(newHeader.Header.Message.Slot)
					ev.NewBlockProposerIndex = int64(newHeader.Header.Message.ProposerIndex)
				}
				dbmodel.InsertNewReorg(ev)
			}
		}
	}

	for {
		select {
		case <-ticker.C:
			eventCh := s.beaconClient.MonitorReorgEvent()
			if eventCh != nil {
				go handler(eventCh)
				ticker.Reset(time.Hour * 256)
			} else {
				ticker.Reset(time.Minute)
			}
		}
	}

}

func (s *Server) monitorDuties() {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	dutyTicker := time.NewTicker(time.Minute)
	defer dutyTicker.Stop()

	slotsPerEpoch := int64(32)
	dumped := make(map[int64]bool)

	for {
		select {

		case <-dutyTicker.C:
			header, err := s.beaconClient.GetLatestBeaconHeader()
			if err != nil {
				log.WithError(err).Debug("duty ticker get latest beacon header failed")
				continue
			}
			curSlot, _ := strconv.ParseInt(header.Header.Message.Slot, 10, 64)
			curEpoch := curSlot / slotsPerEpoch
			nextEpoch := curEpoch + 1
			if curEpoch == 0 && dumped[curEpoch] == false {
				//
				if err := s.dumpDuties(curEpoch); err == nil {
					dumped[curEpoch] = true
				}
			}
			if dumped[nextEpoch] == false {
				if err := s.dumpDuties(nextEpoch); err == nil {
					dumped[nextEpoch] = true
				}
			}

		case <-ticker.C:
			duties, err := s.beaconClient.GetCurrentEpochAttestDuties()
			if err != nil {
				continue
			}
			for _, duty := range duties {
				idx, _ := strconv.Atoi(duty.ValidatorIndex)
				s.validatorSetInfo.AddValidator(idx, duty.Pubkey)
			}

			ticker.Reset(time.Second * 2)

		}
	}
}

func (s *Server) Start() {
	// start RPC endpoints
	err := s.startRPC()
	if err != nil {
		s.stopRPC()
	}
	s.openApi.Start()
	// start collect duties info.
	go s.monitorDuties()
	go s.monitorEvent()
}

func (s *Server) stopRPC() {
	s.http.stop()
}

// implement backend
func (s *Server) SomeNeedBackend() bool {
	return true
}

func (s *Server) GetBlockHeight() (uint64, error) {
	return s.execClient.BlockNumber(context.Background())
}

func (s *Server) GetBlockByNumber(number *big.Int) (*ethtype.Block, error) {
	return s.execClient.BlockByNumber(context.Background(), number)
}

func (s *Server) GetHeightByNumber(number *big.Int) (*ethtype.Header, error) {
	return s.execClient.HeaderByNumber(context.Background(), number)
}

func (s *Server) GetStrategy() *types.Strategy {
	return s.strategy
}

func (s *Server) GetValidatorRoleByPubkey(slot int, pubkey string) types.RoleType {
	if val := s.validatorSetInfo.GetValidatorByPubkey(pubkey); val != nil {
		return s.GetValidatorRole(slot, int(val.Index))
	} else {
		return types.NormalRole
	}
}

func (s *Server) GetCurrentEpochProposeDuties() ([]types.ProposerDuty, error) {
	return s.beaconClient.GetCurrentEpochProposerDuties()
}

func (s *Server) GetCurrentEpochAttestDuties() ([]types.AttestDuty, error) {
	return s.beaconClient.GetCurrentEpochAttestDuties()
}

func (s *Server) GetSlotsPerEpoch() int {
	count, err := s.beaconClient.GetIntConfig(beaconapi.SLOTS_PER_EPOCH)
	if err != nil {
		return 6
	}
	return count
}

func (s *Server) GetIntervalPerSlot() int {
	interval, err := s.beaconClient.GetIntConfig(beaconapi.SECONDS_PER_SLOT)
	if err != nil {
		return 12
	}
	return interval
}

func (s *Server) AddSignedAttestation(slot uint64, pubkey string, attestation *ethpb.Attestation) {
	s.validatorSetInfo.AddSignedAttestation(slot, pubkey, attestation)
}

func (s *Server) AddSignedBlock(slot uint64, pubkey string, block *ethpb.GenericSignedBeaconBlock) {
	s.validatorSetInfo.AddSignedBlock(slot, pubkey, block)
}

func (s *Server) GetAttestSet(slot uint64) *types.SlotAttestSet {
	return s.validatorSetInfo.GetAttestSet(slot)
}

func (s *Server) GetBlockSet(slot uint64) *types.SlotBlockSet {
	return s.validatorSetInfo.GetBlockSet(slot)
}

func (s *Server) GetValidatorDataSet() *types.ValidatorDataSet {
	return s.validatorSetInfo
}

func (s *Server) GetValidatorByProposeSlot(slot uint64) (int, error) {
	epochPerSlot := uint64(s.GetSlotsPerEpoch())
	epoch := slot / epochPerSlot
	duties, err := s.beaconClient.GetProposerDuties(int(epoch))
	if err != nil {
		return 0, err
	}
	for _, duty := range duties {
		dutySlot, _ := strconv.ParseInt(duty.Slot, 10, 64)
		if uint64(dutySlot) == slot {
			idx, _ := strconv.Atoi(duty.ValidatorIndex)
			return idx, nil
		}
	}
	return 0, errors.New("not found")
}

func (s *Server) GetProposeDuties(epoch int) ([]types.ProposerDuty, error) {
	return s.beaconClient.GetProposerDuties(epoch)
}

func (s *Server) SlotsPerEpoch() int {
	return s.GetSlotsPerEpoch()
}

func (s *Server) GetValidatorRole(slot int, valIdx int) types.RoleType {
	if slot < 0 {
		header, err := s.beaconClient.GetLatestBeaconHeader()
		if err != nil {
			return types.NormalRole
		}
		slot, _ = strconv.Atoi(header.Header.Message.Slot)
	}
	return s.strategy.GetValidatorRole(valIdx, int64(slot))
}

func (s *Server) GetInternalSlotStrategy() []slotstrategy.InternalSlotStrategy {
	var err error
	if len(s.internal) == 0 {
		s.internal, err = slotstrategy.ParseToInternalSlotStrategy(s, s.strategy.Slots)
		if err != nil {
			log.WithError(err).Error("parse strategy failed")
			return nil
		}
	}
	return s.internal
}
func (s *Server) GetSlotRoot(slot int64) (string, error) {
	return s.beaconClient.GetSlotRoot(slot)
}

func (s *Server) dumpDuties(epoch int64) error {
	duties, err := s.GetProposeDuties(int(epoch))
	if err != nil {
		return err
	}
	for _, duty := range duties {
		log.WithFields(log.Fields{
			"epoch":     epoch,
			"slot":      duty.Slot,
			"validator": duty.ValidatorIndex,
		}).Info("epoch duty")
	}
	return nil
}

func (s *Server) UpdateStrategy(strategy *types.Strategy) error {
	parsed, err := slotstrategy.ParseToInternalSlotStrategy(s, strategy.Slots)
	if err != nil {
		return err
	}
	s.strategy = strategy
	s.internal = parsed
	return nil
}

func (s *Server) GetSlotStartTime(slot int) (int64, bool) {
	key := fmt.Sprintf("slot_start_time_%d", slot)
	if v, ok := s.cache.Get(key); ok {
		return v.(int64), true
	}
	return 0, false
}

func (s *Server) SetSlotStartTime(slot int, time int64) {
	key := fmt.Sprintf("slot_start_time_%d", slot)
	if _, ok := s.cache.Get(key); !ok {
		s.cache.Add(key, time)
	}
}
