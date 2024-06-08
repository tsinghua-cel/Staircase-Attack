package strategy

import "github.com/tsinghua-cel/attacker-service/types"

type internalValidatorStrategy struct {
	ValidatorIndex    int `json:"validator_index"`
	AttackerStartSlot int `json:"attacker_start_slot"`
	AttackerEndSlot   int `json:"attacker_end_slot"`
}

func parseToInternalValidatorsStrategy(strategy []types.ValidatorStrategy) []internalValidatorStrategy {
	is := make([]internalValidatorStrategy, len(strategy))
	for i, s := range strategy {
		is[i] = internalValidatorStrategy{
			ValidatorIndex:    s.ValidatorIndex,
			AttackerStartSlot: s.AttackerStartSlot,
			AttackerEndSlot:   s.AttackerEndSlot,
		}
	}
	return is
}
