package types

import "encoding/json"

type TotalReward struct {
	ValidatorIndex string `json:"validator_index"`
	Head           string `json:"head"`
	Target         string `json:"target"`
	Source         string `json:"source"`
	InclusionDelay string `json:"inclusion_delay"`
	Inactivity     string `json:"inactivity"`
}

type RewardInfo struct {
	TotalRewards []TotalReward `json:"total_rewards"`
}

type BeaconHeaderInfo struct {
	Header struct {
		Message struct {
			Slot          string `json:"slot"`
			ProposerIndex string `json:"proposer_index"`
			ParentRoot    string `json:"parent_root"`
			StateRoot     string `json:"state_root"`
			BodyRoot      string `json:"body_root"`
		} `json:"message"`
		Signature string `json:"signature"`
	} `json:"header"`
	Root      string `json:"root"`
	Canonical bool   `json:"canonical"`
}

type ProposerDuty struct {
	Pubkey         string `json:"pubkey"`
	ValidatorIndex string `json:"validator_index"`
	Slot           string `json:"slot"`
}

type AttestDuty struct {
	Pubkey                  string `json:"pubkey"`
	ValidatorIndex          string `json:"validator_index"`
	CommitteeIndex          string `json:"committee_index"`
	CommitteeLength         string `json:"committee_length"`
	CommitteesAtSlot        string `json:"committees_at_slot"`
	ValidatorCommitteeIndex string `json:"validator_committee_index"`
	Slot                    string `json:"slot"`
}

type BeaconResponse struct {
	Data json.RawMessage `json:"data"`
}

type ReorgEvent struct {
	Epoch                 int64  `json:"epoch"`
	Slot                  int64  `json:"slot"`
	Depth                 int64  `json:"depth"`
	OldBlockSlot          int64  `json:"old_block_slot"`
	NewBlockSlot          int64  `json:"new_block_slot"`
	OldBlockProposerIndex int64  `json:"old_block_proposer_index"`
	NewBlockProposerIndex int64  `json:"new_block_proposer_index"`
	OldHeadState          string `json:"old_head_state"`
	NewHeadState          string `json:"new_head_state"`
}
