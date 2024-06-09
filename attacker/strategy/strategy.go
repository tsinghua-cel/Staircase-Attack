package strategy

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/types"
	"os"
)

type ValidatorStrategy struct {
	ValidatorIndex    int `json:"validator_index"`
	AttackerStartSlot int `json:"attacker_start_slot"`
	AttackerEndSlot   int `json:"attacker_end_slot"`
}

type BlockStrategy struct {
	DelayEnable    bool  `json:"delay_enable"`
	BroadCastDelay int64 `json:"broad_cast_delay"` // unit millisecond
	ModifyEnable   bool  `json:"modify_enable"`
}

type AttestStrategy struct {
	DelayEnable    bool  `json:"delay_enable"`
	BroadCastDelay int64 `json:"broad_cast_delay"` // unit millisecond
	ModifyEnable   bool  `json:"modify_enable"`
	//lua scripts  => modify attest
}

var (
	defaultValidators    = []ValidatorStrategy{}
	defaultBlockStrategy = BlockStrategy{
		DelayEnable:    false,
		BroadCastDelay: 3000, // 3s
		ModifyEnable:   false,
	}
	defaultAttestStrategy = AttestStrategy{
		DelayEnable:    false,
		BroadCastDelay: 3000, // 3s
		ModifyEnable:   false,
	}
)

type Strategy struct {
	Validators []ValidatorStrategy `json:"validator"`
	Block      BlockStrategy       `json:"block"`
	Attest     AttestStrategy      `json:"attest"`
}

func (s *Strategy) GetValidatorRole(valIdx int, slot int64) types.RoleType {
	for _, v := range s.Validators {
		if v.ValidatorIndex == valIdx {
			if slot >= int64(v.AttackerStartSlot) && slot <= int64(v.AttackerEndSlot) {
				return types.AttackerRole
			}
		}
	}
	return types.NormalRole
}

func ParseStrategy(file string) *Strategy {
	var defautConfig = &Strategy{
		Block:  defaultBlockStrategy,
		Attest: defaultAttestStrategy,
	}
	var s Strategy
	d, err := os.ReadFile(file)
	if err != nil {
		log.WithError(err).Error("read strategy failed, use default config")
		return defautConfig
	}
	if err = json.Unmarshal(d, &s); err != nil {
		log.WithError(err).Error("unmarshal strategy failed, use default config")
		return defautConfig
	}
	return &s
}
