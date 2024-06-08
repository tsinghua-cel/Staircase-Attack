package apis

import (
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/plugins"
	"github.com/tsinghua-cel/attacker-service/strategy/slotstrategy"
	"github.com/tsinghua-cel/attacker-service/types"
)

// AttestAPI offers and API for attestation operations.
type AttestAPI struct {
	b      Backend
	plugin plugins.AttackerPlugin
}

// NewAttestAPI creates a new tx pool service that gives information about the transaction pool.
func NewAttestAPI(b Backend, plugin plugins.AttackerPlugin) *AttestAPI {
	return &AttestAPI{b, plugin}
}

func findMaxLevelStrategy(is []slotstrategy.InternalSlotStrategy, slot int64) (slotstrategy.InternalSlotStrategy, bool) {
	if len(is) == 0 {
		return slotstrategy.InternalSlotStrategy{}, false
	}
	last := is[0]
	for _, s := range is {
		if s.Slot.Compare(slot) == 0 && s.Level > last.Level {
			last = s
		}
	}
	return last, last.Slot.Compare(slot) == 0
}

func (s *AttestAPI) BeforeBroadCast(slot uint64) types.AttackerResponse {
	result := types.AttackerResponse{
		Cmd: types.CMD_NULL,
	}

	if st, find := findMaxLevelStrategy(s.b.GetInternalSlotStrategy(), int64(slot)); find {
		action := st.Actions["AttestBeforeBroadCast"]
		if action != nil {
			r := action.RunAction(s.b, int64(slot), "")
			result.Cmd = r.Cmd
		}
	}
	log.WithFields(log.Fields{
		"cmd":  result.Cmd,
		"slot": slot,
	}).Info("exit AttestBeforeBroadCast")

	return result
}

func (s *AttestAPI) AfterBroadCast(slot uint64) types.AttackerResponse {
	result := types.AttackerResponse{
		Cmd: types.CMD_NULL,
	}
	if st, find := findMaxLevelStrategy(s.b.GetInternalSlotStrategy(), int64(slot)); find {
		action := st.Actions["AttestAfterBroadCast"]
		if action != nil {
			r := action.RunAction(s.b, int64(slot), "")
			result.Cmd = r.Cmd
		}
	}
	log.WithFields(log.Fields{
		"cmd":  result.Cmd,
		"slot": slot,
	}).Info("exit AttestAfterBroadCast")

	return result
}

func (s *AttestAPI) BeforeSign(slot uint64, pubkey string, attestDataBase64 string) types.AttackerResponse {
	result := types.AttackerResponse{
		Cmd:    types.CMD_NULL,
		Result: attestDataBase64,
	}

	attestation, err := common.Base64ToAttestationData(attestDataBase64)
	if err != nil {
		return types.AttackerResponse{
			Cmd:    types.CMD_NULL,
			Result: attestDataBase64,
		}
	}

	if st, find := findMaxLevelStrategy(s.b.GetInternalSlotStrategy(), int64(slot)); find {
		action := st.Actions["AttestBeforeSign"]
		if action != nil {
			r := action.RunAction(s.b, int64(slot), pubkey, attestation)
			result.Cmd = r.Cmd
			newAttestation, ok := r.Result.(*ethpb.AttestationData)
			if ok {
				newData, _ := common.AttestationDataToBase64(newAttestation)
				result.Result = newData
			}
		}
	}
	log.WithFields(log.Fields{
		"cmd":  result.Cmd,
		"slot": slot,
	}).Info("exit AttestBeforeSign")
	return result
}

func (s *AttestAPI) AfterSign(slot uint64, pubkey string, signedAttestDataBase64 string) types.AttackerResponse {
	signedAttestData, err := common.Base64ToSignedAttestation(signedAttestDataBase64)
	if err != nil {
		return types.AttackerResponse{
			Cmd:    types.CMD_NULL,
			Result: signedAttestDataBase64,
		}
	}
	result := types.AttackerResponse{
		Cmd:    types.CMD_NULL,
		Result: signedAttestDataBase64,
	}

	if t, find := findMaxLevelStrategy(s.b.GetInternalSlotStrategy(), int64(slot)); find {
		action := t.Actions["AttestAfterSign"]
		if action != nil {
			r := action.RunAction(s.b, int64(slot), pubkey, signedAttestData)
			result.Cmd = r.Cmd
			newAttestation, ok := r.Result.(*ethpb.Attestation)
			if ok {
				newData, _ := common.SignedAttestationToBase64(newAttestation)
				result.Result = newData
			}
		}
	}
	log.WithFields(log.Fields{
		"cmd":  result.Cmd,
		"slot": slot,
	}).Info("exit AttestAfterSign")
	return result
}

func (s *AttestAPI) BeforePropose(slot uint64, pubkey string, signedAttestDataBase64 string) types.AttackerResponse {
	signedAttest, err := common.Base64ToSignedAttestation(signedAttestDataBase64)
	if err != nil {
		return types.AttackerResponse{
			Cmd:    types.CMD_NULL,
			Result: signedAttestDataBase64,
		}
	}
	result := types.AttackerResponse{
		Cmd:    types.CMD_NULL,
		Result: signedAttestDataBase64,
	}

	if t, find := findMaxLevelStrategy(s.b.GetInternalSlotStrategy(), int64(slot)); find {
		action := t.Actions["AttestBeforePropose"]
		if action != nil {
			r := action.RunAction(s.b, int64(slot), pubkey, signedAttest)
			result.Cmd = r.Cmd
			newAttestation, ok := r.Result.(*ethpb.Attestation)
			if ok {
				newData, _ := common.SignedAttestationToBase64(newAttestation)
				result.Result = newData
			}
		}
	}
	log.WithFields(log.Fields{
		"cmd":  result.Cmd,
		"slot": slot,
	}).Info("exit AttestBeforePropose")
	return result
}

func (s *AttestAPI) AfterPropose(slot uint64, pubkey string, signedAttestDataBase64 string) types.AttackerResponse {
	signedAttest, err := common.Base64ToSignedAttestation(signedAttestDataBase64)
	if err != nil {
		return types.AttackerResponse{
			Cmd:    types.CMD_NULL,
			Result: signedAttestDataBase64,
		}
	}
	result := types.AttackerResponse{
		Cmd:    types.CMD_NULL,
		Result: signedAttestDataBase64,
	}

	if t, find := findMaxLevelStrategy(s.b.GetInternalSlotStrategy(), int64(slot)); find {
		action := t.Actions["AttestAfterPropose"]
		if action != nil {
			r := action.RunAction(s.b, int64(slot), pubkey, signedAttest)
			result.Cmd = r.Cmd
			newAttestation, ok := r.Result.(*ethpb.Attestation)
			if ok {
				newData, _ := common.SignedAttestationToBase64(newAttestation)
				result.Result = newData
			}
		}
	}

	log.WithFields(log.Fields{
		"cmd":  result.Cmd,
		"slot": slot,
	}).Info("exit AttestAfterPropose")

	return result
}
