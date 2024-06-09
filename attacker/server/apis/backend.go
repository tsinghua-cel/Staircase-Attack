package apis

import (
	"github.com/ethereum/go-ethereum/core/types"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/tsinghua-cel/attacker-service/beaconapi"
	"github.com/tsinghua-cel/attacker-service/rpc"
	"github.com/tsinghua-cel/attacker-service/strategy"
	types2 "github.com/tsinghua-cel/attacker-service/types"
	"github.com/tsinghua-cel/attacker-service/validatorSet"
	"math/big"
)

// Backend interface provides the common API services (that are provided by
// both full and light clients) with access to necessary functions.
type Backend interface {
	SomeNeedBackend() bool
	// update strategy
	GetStrategy() *strategy.Strategy
	UpdateBlockBroadDelay(milliSecond int64) error
	UpdateAttestBroadDelay(milliSecond int64) error

	// get data from execute node.
	GetBlockHeight() (uint64, error)
	GetBlockByNumber(number *big.Int) (*types.Block, error)
	GetHeightByNumber(number *big.Int) (*types.Header, error)

	GetValidatorRole(slot int, valIdx int) types2.RoleType
	GetValidatorRoleByPubkey(slot int, pubkey string) types2.RoleType
	GetCurrentEpochProposeDuties() ([]beaconapi.ProposerDuty, error)
	GetSlotsPerEpoch() int
	SlotsPerEpoch() int
	GetIntervalPerSlot() int
	AddSignedAttestation(slot uint64, pubkey string, attestation *ethpb.Attestation)
	AddSignedBlock(slot uint64, pubkey string, block *ethpb.GenericSignedBeaconBlock)
	GetAttestSet(slot uint64) *validatorSet.SlotAttestSet
	GetBlockSet(slot uint64) *validatorSet.SlotBlockSet
	GetValidatorDataSet() *validatorSet.ValidatorDataSet
	GetValidatorByProposeSlot(slot uint64) (int, error)
	GetProposeDuties(epoch int) ([]beaconapi.ProposerDuty, error)
}

func GetAPIs(apiBackend Backend) []rpc.API {
	return []rpc.API{
		{
			Namespace: "admin",
			Service:   NewAdminAPI(apiBackend),
		},
		{
			Namespace: "block",
			Service:   NewBlockAPI(apiBackend),
		},
		{
			Namespace: "attest",
			Service:   NewAttestAPI(apiBackend),
		},
	}
}
