package apis

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/prysmaticlabs/prysm/v4/cache/lru"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	attaggregation "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1/attestation/aggregation/attestations"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/strategy"
	"github.com/tsinghua-cel/attacker-service/types"
	"google.golang.org/protobuf/proto"
)

var (
	ErrNilObject              = errors.New("nil object")
	ErrUnsupportedBeaconBlock = errors.New("unsupported beacon block")
	blockCacheContent         = lru.New(1000)
)

// BlockAPI offers and API for block operations.
type BlockAPI struct {
	b Backend
}

// NewBlockAPI creates a new tx pool service that gives information about the transaction pool.
func NewBlockAPI(b Backend) *BlockAPI {
	return &BlockAPI{b}
}

func (s *BlockAPI) GetStrategy(cliInfo string) []byte {
	d, _ := json.Marshal(s.b.GetStrategy().Block)
	return d
}

func (s *BlockAPI) UpdateStrategy(data []byte) error {
	var blockStrategy strategy.BlockStrategy
	if err := json.Unmarshal(data, &blockStrategy); err != nil {
		return err
	}
	s.b.GetStrategy().Block = blockStrategy
	log.Infof("block strategy updated to %v\n", blockStrategy)
	return nil
}

func (s *BlockAPI) BroadCastDelay() types.AttackerResponse {
	bs := s.b.GetStrategy().Block
	if !bs.DelayEnable {
		return types.AttackerResponse{
			Cmd: types.CMD_NULL,
		}
	}
	time.Sleep(time.Millisecond * time.Duration(s.b.GetStrategy().Block.BroadCastDelay))
	return types.AttackerResponse{
		Cmd: types.CMD_NULL,
	}
}

func (s *BlockAPI) modifyBlock(slot uint64, pubkey string, blockDataBase64 string) types.AttackerResponse {
	// 1. 只有每个epoch最后一个出块的恶意节点出块，其他节点不出快
	valIdx, err := s.b.GetValidatorByProposeSlot(slot)
	if err != nil {
		val := s.b.GetValidatorDataSet().GetValidatorByPubkey(pubkey)
		if val == nil {
			return types.AttackerResponse{
				Cmd:    types.CMD_NULL,
				Result: blockDataBase64,
			}
		}
		valIdx = int(val.Index)
	}
	role := s.b.GetValidatorRole(int(slot), valIdx)
	log.WithFields(log.Fields{
		"slot":   slot,
		"valIdx": valIdx,
		"role":   role,
	}).Info("in modify block, get validator by propose slot")

	if role != types.AttackerRole {
		return types.AttackerResponse{
			Cmd:    types.CMD_NULL,
			Result: blockDataBase64,
		}
	}
	epoch := SlotTool{s.b}.SlotToEpoch(int(slot))

	duties, err := s.b.GetProposeDuties(int(epoch))
	if err != nil {
		return types.AttackerResponse{
			Cmd:    types.CMD_NULL,
			Result: blockDataBase64,
		}
	}

	latestSlotWithAttacker := int64(-1)
	for _, duty := range duties {
		dutySlot, _ := strconv.ParseInt(duty.Slot, 10, 64)
		dutyValIdx, _ := strconv.Atoi(duty.ValidatorIndex)
		if s.b.GetValidatorRole(int(slot), dutyValIdx) == types.AttackerRole && dutySlot > latestSlotWithAttacker {
			latestSlotWithAttacker = dutySlot
		}
	}
	log.WithFields(log.Fields{
		"slot":               slot,
		"latestAttackerSlot": latestSlotWithAttacker,
	}).Info("modify block")

	if slot != uint64(latestSlotWithAttacker) {
		// 不是最后一个出块的恶意节点，不出块
		return types.AttackerResponse{
			Cmd:    types.CMD_RETURN,
			Result: blockDataBase64,
		}
	}

	genericBlock, err := s.getGenericSignedBlockFromData(blockDataBase64)
	if err != nil {
		log.WithError(err).Error("get block from data failed")
		return types.AttackerResponse{
			Cmd:    types.CMD_NULL,
			Result: blockDataBase64,
		}
	}
	block, err := s.getCapellaBlockFromGenericSigned(genericBlock)
	if err != nil {
		log.WithError(err).Error("get block from data failed")
		return types.AttackerResponse{
			Cmd:    types.CMD_NULL,
			Result: blockDataBase64,
		}
	}

	// 3.出的块的一个字段attestation要包含其他恶意节点的attestation。
	startEpoch := SlotTool{s.b}.EpochStart(epoch)
	endEpoch := SlotTool{s.b}.EpochEnd(epoch)
	attackerAttestations := make([]*ethpb.Attestation, 0)
	validatorSet := s.b.GetValidatorDataSet()
	for i := startEpoch; i <= endEpoch; i++ {
		allSlotAttest := s.b.GetAttestSet(uint64(i))
		if allSlotAttest == nil {
			continue
		}

		for publicKey, att := range allSlotAttest.Attestations {
			val := validatorSet.GetValidatorByPubkey(publicKey)
			valRole := s.b.GetValidatorRole(int(i), int(val.Index))
			if val != nil && valRole == types.AttackerRole {
				log.WithField("pubkey", publicKey).Debug("add attacker attestation to block")
				attackerAttestations = append(attackerAttestations, att)
			}
		}
	}

	allAtt := append(block.Capella.Block.Body.Attestations, attackerAttestations...)
	{
		// Remove duplicates from both aggregated/unaggregated attestations. This
		// prevents inefficient aggregates being created.
		atts, _ := types.ProposerAtts(allAtt).Dedup()
		attsByDataRoot := make(map[[32]byte][]*ethpb.Attestation, len(atts))
		for _, att := range atts {
			attDataRoot, err := att.Data.HashTreeRoot()
			if err != nil {
			}
			attsByDataRoot[attDataRoot] = append(attsByDataRoot[attDataRoot], att)
		}

		attsForInclusion := types.ProposerAtts(make([]*ethpb.Attestation, 0))
		for _, as := range attsByDataRoot {
			as, err := attaggregation.Aggregate(as)
			if err != nil {
				continue
			}
			attsForInclusion = append(attsForInclusion, as...)
		}
		deduped, _ := attsForInclusion.Dedup()
		sorted, err := deduped.SortByProfitability()
		if err != nil {
			log.WithError(err).Error("sort attestation failed")
		} else {
			atts = sorted.LimitToMaxAttestations()
		}
		allAtt = atts
	}

	block.Capella.Block.Body.Attestations = allAtt

	// 4. encode to base64.
	genericBlock.Block = block

	resBlockBase64, err := s.genericSignedBlockToBase64(genericBlock)
	if err != nil {
		return types.AttackerResponse{
			Cmd:    types.CMD_NULL,
			Result: blockDataBase64,
		}
	}
	return types.AttackerResponse{
		Cmd:    types.CMD_NULL,
		Result: resBlockBase64,
	}
}

func (s *BlockAPI) getCapellaBlockFromGenericSigned(block *ethpb.GenericSignedBeaconBlock) (*ethpb.GenericSignedBeaconBlock_Capella, error) {
	switch b := block.Block.(type) {
	case nil:
		return nil, ErrNilObject
	case *ethpb.GenericSignedBeaconBlock_Phase0:
		return nil, ErrUnsupportedBeaconBlock
	case *ethpb.GenericSignedBeaconBlock_Altair:
		return nil, ErrUnsupportedBeaconBlock
	case *ethpb.GenericSignedBeaconBlock_Bellatrix:
		return nil, ErrUnsupportedBeaconBlock
	case *ethpb.GenericSignedBeaconBlock_BlindedBellatrix:
		return nil, ErrUnsupportedBeaconBlock
	case *ethpb.GenericSignedBeaconBlock_Capella:
		return b, nil
	case *ethpb.GenericSignedBeaconBlock_BlindedCapella:
		return nil, ErrUnsupportedBeaconBlock
	case *ethpb.GenericSignedBeaconBlock_Deneb:
		return nil, ErrUnsupportedBeaconBlock
	case *ethpb.GenericSignedBeaconBlock_BlindedDeneb:
		return nil, ErrUnsupportedBeaconBlock
	default:
		log.WithError(ErrUnsupportedBeaconBlock).Errorf("unsupported beacon block from type %T", b)
		return nil, ErrUnsupportedBeaconBlock
	}
}

func (s *BlockAPI) getCapellaBlockFromGeneric(block *ethpb.GenericBeaconBlock) (*ethpb.GenericBeaconBlock_Capella, error) {
	switch b := block.Block.(type) {
	case nil:
		return nil, ErrNilObject
	case *ethpb.GenericBeaconBlock_Phase0:
		return nil, ErrUnsupportedBeaconBlock
	case *ethpb.GenericBeaconBlock_Altair:
		return nil, ErrUnsupportedBeaconBlock
	case *ethpb.GenericBeaconBlock_Bellatrix:
		return nil, ErrUnsupportedBeaconBlock
	case *ethpb.GenericBeaconBlock_BlindedBellatrix:
		return nil, ErrUnsupportedBeaconBlock
	case *ethpb.GenericBeaconBlock_Capella:
		return b, nil
	case *ethpb.GenericBeaconBlock_BlindedCapella:
		return nil, ErrUnsupportedBeaconBlock
	case *ethpb.GenericBeaconBlock_Deneb:
		return nil, ErrUnsupportedBeaconBlock
	case *ethpb.GenericBeaconBlock_BlindedDeneb:
		return nil, ErrUnsupportedBeaconBlock
	default:
		log.WithError(ErrUnsupportedBeaconBlock).Errorf("unsupported beacon block from type %T", b)
		return nil, ErrUnsupportedBeaconBlock
	}
}

func (s *BlockAPI) getGenericBlockFromData(blockDataBase64 string) (*ethpb.GenericBeaconBlock, error) {
	blockData, err := base64.StdEncoding.DecodeString(blockDataBase64)
	if err != nil {
		log.WithError(err).Error("base64 decode block data failed")
		return nil, err
	}
	var block = new(ethpb.GenericBeaconBlock)
	if err := proto.Unmarshal(blockData, block); err != nil {
		log.WithError(err).Error("unmarshal block data failed")
		return nil, err
	}
	return block, nil
}

func (s *BlockAPI) getGenericSignedBlockFromData(signedBlockDataBase64 string) (*ethpb.GenericSignedBeaconBlock, error) {
	blockData, err := base64.StdEncoding.DecodeString(signedBlockDataBase64)
	if err != nil {
		log.WithError(err).Error("base64 decode block data failed")
		return nil, err
	}
	var block = new(ethpb.GenericSignedBeaconBlock)
	if err := proto.Unmarshal(blockData, block); err != nil {
		log.WithError(err).Error("unmarshal block data failed")
		return nil, err
	}
	return block, nil
}

func (s *BlockAPI) genericBlockToBase64(block *ethpb.GenericBeaconBlock) (string, error) {
	data, err := proto.Marshal(block)
	if err != nil {
		log.WithError(err).Error("marshal block data failed")
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (s *BlockAPI) genericSignedBlockToBase64(block *ethpb.GenericSignedBeaconBlock) (string, error) {
	data, err := proto.Marshal(block)
	if err != nil {
		log.WithError(err).Error("marshal block data failed")
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (s *BlockAPI) DelayForReceiveBlock(slot uint64) types.AttackerResponse {
	valIdx, err := s.b.GetValidatorByProposeSlot(slot)
	if err != nil {
		return types.AttackerResponse{
			Cmd: types.CMD_NULL,
		}
	}
	{
		valRole := s.b.GetValidatorRole(int(slot), valIdx)
		//val := s.b.GetValidatorDataSet().GetValidatorByIndex(valIdx)
		if valRole != types.AttackerRole {
			return types.AttackerResponse{
				Cmd: types.CMD_NULL,
			}
		}

		duties, err := s.b.GetCurrentEpochProposeDuties()
		if err != nil {
			return types.AttackerResponse{
				Cmd: types.CMD_NULL,
			}
		}

		latestAttackerVal := int64(-1)
		for _, duty := range duties {
			dutySlot, _ := strconv.Atoi(duty.Slot)
			dutyValIdx, _ := strconv.Atoi(duty.ValidatorIndex)
			if s.b.GetValidatorRole(dutySlot, dutyValIdx) == types.AttackerRole {
				latestAttackerVal = int64(dutyValIdx)
			}
		}
		if valIdx != int(latestAttackerVal) {
			// 不是最后一个出块的恶意节点，不出块
			return types.AttackerResponse{
				Cmd: types.CMD_RETURN,
			}
		}
	}
	// 当前是最后一个出块的恶意节点，进行延时

	epochSlots := s.b.GetSlotsPerEpoch()
	seconds := s.b.GetIntervalPerSlot()
	delay := (epochSlots - int(slot%uint64(epochSlots))) * seconds
	time.Sleep(time.Second * time.Duration(delay))
	key := fmt.Sprintf("delay_%d_%d", slot, valIdx)
	blockCacheContent.Add(key, delay)
	log.WithFields(log.Fields{
		"slot":     slot,
		"validx":   valIdx,
		"duration": delay,
	}).Info("delay for receive block")

	return types.AttackerResponse{
		Cmd: types.CMD_UPDATE_STATE,
	}
}

func (s *BlockAPI) BeforeBroadCast(slot uint64) types.AttackerResponse {
	valIdx, err := s.b.GetValidatorByProposeSlot(slot)
	if err != nil {
		return types.AttackerResponse{
			Cmd: types.CMD_NULL,
		}
	}
	{
		valRole := s.b.GetValidatorRole(int(slot), valIdx)
		//val := s.b.GetValidatorDataSet().GetValidatorByIndex(valIdx)
		if valRole != types.AttackerRole {
			return types.AttackerResponse{
				Cmd: types.CMD_NULL,
			}
		}

		duties, err := s.b.GetCurrentEpochProposeDuties()
		if err != nil {
			return types.AttackerResponse{
				Cmd: types.CMD_NULL,
			}
		}

		latestAttackerVal := int64(-1)
		for _, duty := range duties {
			dutySlot, _ := strconv.Atoi(duty.Slot)
			dutyValIdx, _ := strconv.Atoi(duty.ValidatorIndex)
			if s.b.GetValidatorRole(dutySlot, dutyValIdx) == types.AttackerRole {
				latestAttackerVal = int64(dutyValIdx)
			}
		}
		if valIdx != int(latestAttackerVal) {
			// 不是最后一个出块的恶意节点，不出块
			return types.AttackerResponse{
				Cmd: types.CMD_RETURN,
			}
		}
	}
	// 当前是最后一个出块的恶意节点，进行延时
	key := fmt.Sprintf("delay_%d_%d", slot, valIdx)
	lastDelay := 0
	if t, exist := blockCacheContent.Get(key); exist {
		lastDelay = t.(int)
	}
	seconds := s.b.GetIntervalPerSlot()
	n2delay := 12 * seconds
	total := n2delay + lastDelay
	time.Sleep(time.Second * time.Duration(total))
	log.WithFields(log.Fields{
		"slot":     slot,
		"validx":   valIdx,
		"duration": total,
	}).Info("delay for beforeBroadcastBlock")

	return types.AttackerResponse{
		Cmd: types.CMD_NULL,
	}
}

func (s *BlockAPI) AfterBroadCast(slot uint64) types.AttackerResponse {
	return types.AttackerResponse{
		Cmd: types.CMD_NULL,
	}
}

func (s *BlockAPI) BeforeMakeBlock(slot uint64, pubkey string) types.AttackerResponse {
	// 1. 只有每个epoch最后一个出块的恶意节点出块，其他节点不出快
	valIdx, err := s.b.GetValidatorByProposeSlot(slot)
	if err != nil {
		val := s.b.GetValidatorDataSet().GetValidatorByPubkey(pubkey)
		if val == nil {
			return types.AttackerResponse{
				Cmd: types.CMD_NULL,
			}
		}
		valIdx = int(val.Index)
	}
	role := s.b.GetValidatorRole(int(slot), valIdx)
	log.WithFields(log.Fields{
		"slot":   slot,
		"valIdx": valIdx,
		"role":   role,
	}).Info("in BeforeMakeBlock, get validator by propose slot")

	if role != types.AttackerRole {
		return types.AttackerResponse{
			Cmd: types.CMD_NULL,
		}
	}
	epoch := SlotTool{s.b}.SlotToEpoch(int(slot))

	duties, err := s.b.GetProposeDuties(int(epoch))
	if err != nil {
		return types.AttackerResponse{
			Cmd: types.CMD_NULL,
		}
	}

	latestSlotWithAttacker := int64(-1)
	for _, duty := range duties {
		dutySlot, _ := strconv.ParseInt(duty.Slot, 10, 64)
		dutyValIdx, _ := strconv.Atoi(duty.ValidatorIndex)
		log.WithFields(log.Fields{
			"slot":   dutySlot,
			"valIdx": dutyValIdx,
		}).Debug("duty slot")
		if s.b.GetValidatorRole(int(slot), dutyValIdx) == types.AttackerRole && dutySlot > latestSlotWithAttacker {
			latestSlotWithAttacker = dutySlot
			log.WithField("latestSlotWithAttacker", latestSlotWithAttacker).Debug("update latestSlotWithAttacker")
		}
	}
	log.WithFields(log.Fields{
		"slot":               slot,
		"latestAttackerSlot": latestSlotWithAttacker,
	}).Info("modify block")

	if slot != uint64(latestSlotWithAttacker) {
		// 不是最后一个恶意的出块，不出块
		return types.AttackerResponse{
			Cmd: types.CMD_RETURN,
		}
	}

	return types.AttackerResponse{
		Cmd: types.CMD_NULL,
	}
}

func (s *BlockAPI) BeforeSign(slot uint64, pubkey string, blockDataBase64 string) types.AttackerResponse {
	modifyBlockRes := s.modifyBlock(slot, pubkey, blockDataBase64)
	return modifyBlockRes
}

func (s *BlockAPI) AfterSign(slot uint64, pubkey string, signedBlockDataBase64 string) types.AttackerResponse {
	valIdx, err := s.b.GetValidatorByProposeSlot(slot)
	if err != nil {
		val := s.b.GetValidatorDataSet().GetValidatorByPubkey(pubkey)
		if val == nil {
			return types.AttackerResponse{
				Cmd: types.CMD_NULL,
			}
		}
		valIdx = int(val.Index)
	}
	role := s.b.GetValidatorRole(int(slot), valIdx)
	log.WithFields(log.Fields{
		"slot":   slot,
		"valIdx": valIdx,
		"role":   role,
	}).Info("in AfterSign, get validator by propose slot")

	if role != types.AttackerRole {
		return types.AttackerResponse{
			Cmd: types.CMD_NULL,
		}
	}
	epoch := SlotTool{s.b}.SlotToEpoch(int(slot))

	duties, err := s.b.GetProposeDuties(int(epoch))
	if err != nil {
		return types.AttackerResponse{
			Cmd: types.CMD_NULL,
		}
	}

	latestSlotWithAttacker := int64(-1)
	for _, duty := range duties {
		dutySlot, _ := strconv.ParseInt(duty.Slot, 10, 64)
		dutyValIdx, _ := strconv.Atoi(duty.ValidatorIndex)
		log.WithFields(log.Fields{
			"slot":   dutySlot,
			"valIdx": dutyValIdx,
		}).Debug("duty slot")
		if s.b.GetValidatorRole(int(slot), dutyValIdx) == types.AttackerRole && dutySlot > latestSlotWithAttacker {
			latestSlotWithAttacker = dutySlot
			log.WithField("latestSlotWithAttacker", latestSlotWithAttacker).Debug("update latestSlotWithAttacker")
		}
	}

	if slot != uint64(latestSlotWithAttacker) {
		// 不是最后一个恶意的出块，不出块
		return types.AttackerResponse{
			Cmd: types.CMD_RETURN,
		}
	}

	return types.AttackerResponse{
		Cmd:    types.CMD_NULL,
		Result: signedBlockDataBase64,
	}
}

func (s *BlockAPI) BeforePropose(slot uint64, pubkey string, signedBlockDataBase64 string) types.AttackerResponse {
	return types.AttackerResponse{
		Cmd:    types.CMD_NULL,
		Result: signedBlockDataBase64,
	}
}

func (s *BlockAPI) AfterPropose(slot uint64, pubkey string, signedBlockDataBase64 string) types.AttackerResponse {
	return types.AttackerResponse{
		Cmd:    types.CMD_NULL,
		Result: signedBlockDataBase64,
	}

}
