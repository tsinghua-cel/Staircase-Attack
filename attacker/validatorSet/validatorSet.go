package validatorSet

import (
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"strings"
	"sync"
)

type ValidatorInfo struct {
	Index  int64  `json:"index"`
	Pubkey string `json:"pubkey"`
	//Role   types.RoleType `json:"role"`
	//Attests ValidatorAttestSet `json:"attests"`
	//Blocks  ValidatorBlockSet  `json:"blocks"`
}

type ValidatorDataSet struct {
	ValidatorByIndex  sync.Map                  //map[int]*ValidatorInfo
	ValidatorByPubkey sync.Map                  //map[string]*ValidatorInfo
	AttestSet         map[uint64]*SlotAttestSet // epoch -> attestation
	BlockSet          map[uint64]*SlotBlockSet  // epoch -> block
	lock              sync.RWMutex
}

func NewValidatorSet() *ValidatorDataSet {
	return &ValidatorDataSet{
		AttestSet: make(map[uint64]*SlotAttestSet),
		BlockSet:  make(map[uint64]*SlotBlockSet),
	}
}

func padPubkey(p string) string {
	if strings.HasPrefix(p, "0x") {
		return p
	}
	return "0x" + p
}

func (vs *ValidatorDataSet) AddValidator(index int, pubkey string) {
	pubkey = padPubkey(pubkey)
	vs.lock.Lock()
	defer vs.lock.Unlock()
	v := &ValidatorInfo{
		Index:  int64(index),
		Pubkey: pubkey,
	}
	vs.ValidatorByIndex.Store(index, v)
	vs.ValidatorByPubkey.Store(pubkey, v)
}

func (vs *ValidatorDataSet) GetValidatorByIndex(index int) *ValidatorInfo {
	vs.lock.RLock()
	defer vs.lock.RUnlock()
	if v, exist := vs.ValidatorByIndex.Load(index); !exist {
		return nil
	} else {
		return v.(*ValidatorInfo)
	}
}

func (vs *ValidatorDataSet) GetValidatorByPubkey(pubkey string) *ValidatorInfo {
	pubkey = padPubkey(pubkey)
	vs.lock.RLock()
	defer vs.lock.RUnlock()
	if v, exist := vs.ValidatorByPubkey.Load(pubkey); !exist {
		return nil
	} else {
		return v.(*ValidatorInfo)
	}
}

func (vs *ValidatorDataSet) GetAttestSet(slot uint64) *SlotAttestSet {
	vs.lock.RLock()
	defer vs.lock.RUnlock()
	if v, exist := vs.AttestSet[slot]; !exist {
		return nil
	} else {
		return v
	}
}

func (vs *ValidatorDataSet) GetBlockSet(slot uint64) *SlotBlockSet {
	vs.lock.RLock()
	defer vs.lock.RUnlock()
	if v, exist := vs.BlockSet[slot]; !exist {
		return nil
	} else {
		return v
	}
}

func (vs *ValidatorDataSet) AddSignedAttestation(slot uint64, pubkey string, attestation *ethpb.Attestation) {
	pubkey = padPubkey(pubkey)
	vs.lock.Lock()
	defer vs.lock.Unlock()

	if _, exist := vs.AttestSet[slot]; !exist {
		vs.AttestSet[slot] = &SlotAttestSet{
			Attestations: make(map[string]*ethpb.Attestation),
		}
	}
	vs.AttestSet[slot].Attestations[pubkey] = attestation
}

func (vs *ValidatorDataSet) AddSignedBlock(slot uint64, pubkey string, block *ethpb.GenericSignedBeaconBlock) {
	pubkey = padPubkey(pubkey)
	vs.lock.Lock()
	defer vs.lock.Unlock()

	if _, exist := vs.BlockSet[slot]; !exist {
		vs.BlockSet[slot] = &SlotBlockSet{
			Blocks: make(map[string]*ethpb.GenericSignedBeaconBlock),
		}
	}
	vs.BlockSet[slot].Blocks[pubkey] = block
}

type SlotAttestSet struct {
	Attestations map[string]*ethpb.Attestation
}

type SlotBlockSet struct {
	Blocks map[string]*ethpb.GenericSignedBeaconBlock
}

type ValidatorAttestSet struct {
	Attestations map[uint64]*ethpb.Attestation
}

type ValidatorBlockSet struct {
	Blocks map[uint64]ethpb.GenericBeaconBlock
}
