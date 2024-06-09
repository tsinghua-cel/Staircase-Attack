package apis

import (
	"encoding/base64"
	"encoding/json"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/strategy"
	"github.com/tsinghua-cel/attacker-service/types"
	"google.golang.org/protobuf/proto"
)

// AttestAPI offers and API for attestation operations.
type AttestAPI struct {
	b Backend
}

// NewBlockAPI creates a new tx pool service that gives information about the transaction pool.
func NewAttestAPI(b Backend) *AttestAPI {
	return &AttestAPI{b}
}

func (s *AttestAPI) GetStrategy() []byte {
	d, _ := json.Marshal(s.b.GetStrategy().Attest)
	return d
}

func (s *AttestAPI) UpdateStrategy(data []byte) error {
	var attestStrategy strategy.AttestStrategy
	if err := json.Unmarshal(data, &attestStrategy); err != nil {
		return err
	}
	s.b.GetStrategy().Attest = attestStrategy
	log.Infof("attest strategy updated to %v\n", attestStrategy)
	return nil
}

func (s *AttestAPI) BeforeBroadCast(slot uint64) types.AttackerResponse {
	return types.AttackerResponse{
		Cmd: types.CMD_NULL,
	}
}

func (s *AttestAPI) AfterBroadCast(slot uint64) types.AttackerResponse {
	return types.AttackerResponse{
		Cmd: types.CMD_NULL,
	}
}

func (s *AttestAPI) BeforeSign(slot uint64, pubkey string, attestDataBase64 string) types.AttackerResponse {
	return types.AttackerResponse{
		Cmd:    types.CMD_NULL,
		Result: attestDataBase64,
	}
}

func (s *AttestAPI) AfterSign(slot uint64, pubkey string, signedAttestDataBase64 string) types.AttackerResponse {
	signedAttestData, err := base64.StdEncoding.DecodeString(signedAttestDataBase64)
	if err != nil {
		log.WithError(err).Error("base64 decode attest data failed")
		return types.AttackerResponse{
			Cmd:    types.CMD_NULL,
			Result: signedAttestDataBase64,
		}
	}
	if role := s.b.GetValidatorRoleByPubkey(int(slot), pubkey); role == types.NormalRole {
		return types.AttackerResponse{
			Cmd:    types.CMD_NULL,
			Result: signedAttestDataBase64,
		}
	}
	var attest = new(ethpb.Attestation)
	if err := proto.Unmarshal(signedAttestData, attest); err != nil {
		log.WithError(err).Error("unmarshal attest data failed")
		return types.AttackerResponse{
			Cmd:    types.CMD_NULL,
			Result: signedAttestDataBase64,
		}
	}
	log.WithFields(log.Fields{
		"slot":   slot,
		"pubkey": pubkey,
	}).Debug("receive signed attest")
	s.b.AddSignedAttestation(slot, pubkey, attest)

	return types.AttackerResponse{
		Cmd:    types.CMD_NULL,
		Result: signedAttestDataBase64,
	}
}

func (s *AttestAPI) BeforePropose(slot uint64, pubkey string, signedAttestDataBase64 string) types.AttackerResponse {
	isAttacker := false
	if s.b.GetValidatorRoleByPubkey(int(slot), pubkey) == types.AttackerRole {
		isAttacker = true
	}

	if isAttacker { // 所有的恶意节点不广播Attestation.
		log.WithFields(log.Fields{}).Debug("this is attacker, not broadcast attest")
		return types.AttackerResponse{
			Cmd: types.CMD_RETURN,
		}
	} else {
		return types.AttackerResponse{
			Cmd: types.CMD_NULL,
		}
	}
}

func (s *AttestAPI) AfterPropose(slot uint64, pubkey string, signedAttestDataBase64 string) types.AttackerResponse {
	return types.AttackerResponse{
		Cmd:    types.CMD_NULL,
		Result: signedAttestDataBase64,
	}
}
