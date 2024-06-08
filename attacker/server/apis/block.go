package apis

import (
	"errors"
	"github.com/prysmaticlabs/prysm/v4/cache/lru"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/plugins"
	"github.com/tsinghua-cel/attacker-service/types"
	"time"
)

var (
	ErrNilObject              = errors.New("nil object")
	ErrUnsupportedBeaconBlock = errors.New("unsupported beacon block")
	blockCacheContent         = lru.New(1000)
)

// BlockAPI offers and API for block operations.
type BlockAPI struct {
	b      Backend
	plugin plugins.AttackerPlugin
}

// NewBlockAPI creates a new tx pool service that gives information about the transaction pool.
func NewBlockAPI(b Backend, plugin plugins.AttackerPlugin) *BlockAPI {
	return &BlockAPI{b, plugin}
}

func (s *BlockAPI) GetNewParentRoot(slot uint64, pubkey string, parentRoot string) types.AttackerResponse {
	result := types.AttackerResponse{
		Cmd:    types.CMD_NULL,
		Result: parentRoot,
	}
	if t, find := findMaxLevelStrategy(s.b.GetInternalSlotStrategy(), int64(slot)); find {
		action := t.Actions["BlockGetNewParentRoot"]
		if action != nil {
			r := action.RunAction(s.b, int64(slot), pubkey, parentRoot)
			result.Cmd = r.Cmd
			if r.Result != nil {
				result.Result = r.Result.(string)
			}
		}
	}
	log.WithFields(log.Fields{
		"cmd":    result.Cmd,
		"slot":   slot,
		"action": "BlockGetNewParentRoot",
	}).Info("exit GetNewParentRoot")

	return result
}

func (s *BlockAPI) BroadCastDelay(slot uint64) types.AttackerResponse {
	return s.todoActionsWithSlot(slot, "BlockDelayForBroadCast")
}

func (s *BlockAPI) DelayForReceiveBlock(slot uint64) types.AttackerResponse {
	s.b.SetSlotStartTime(int(slot), time.Now().Unix())
	return s.todoActionsWithSlot(slot, "BlockDelayForReceiveBlock")
}

func (s *BlockAPI) BeforeBroadCast(slot uint64) types.AttackerResponse {
	return s.todoActionsWithSlot(slot, "BlockBeforeBroadCast")
}

func (s *BlockAPI) AfterBroadCast(slot uint64) types.AttackerResponse {
	return s.todoActionsWithSlot(slot, "BlockAfterBroadCast")
}

func (s *BlockAPI) BeforeSign(slot uint64, pubkey string, signedBlockDataBase64 string) types.AttackerResponse {
	return s.todoActionsWithSignedBlock(slot, pubkey, signedBlockDataBase64, "BlockBeforeSign")
}

func (s *BlockAPI) AfterSign(slot uint64, pubkey string, signedBlockDataBase64 string) types.AttackerResponse {
	return s.todoActionsWithSignedBlock(slot, pubkey, signedBlockDataBase64, "BlockAfterSign")
}

func (s *BlockAPI) BeforePropose(slot uint64, pubkey string, signedBlockDataBase64 string) types.AttackerResponse {
	return s.todoActionsWithSignedBlock(slot, pubkey, signedBlockDataBase64, "BlockBeforePropose")
}

func (s *BlockAPI) AfterPropose(slot uint64, pubkey string, signedBlockDataBase64 string) types.AttackerResponse {
	return s.todoActionsWithSignedBlock(slot, pubkey, signedBlockDataBase64, "BlockAfterPropose")
}

func (s *BlockAPI) todoActionsWithSlot(slot uint64, name string) types.AttackerResponse {
	result := types.AttackerResponse{
		Cmd: types.CMD_NULL,
	}

	if t, find := findMaxLevelStrategy(s.b.GetInternalSlotStrategy(), int64(slot)); find {
		action := t.Actions[name]
		if action != nil {
			r := action.RunAction(s.b, int64(slot), "")
			result.Cmd = r.Cmd
		}
	}
	log.WithFields(log.Fields{
		"cmd":    result.Cmd,
		"slot":   slot,
		"action": name,
	}).Info("exit todoActionsWithSlot")

	return result
}

func (s *BlockAPI) todoActionsWithSignedBlock(slot uint64, pubkey string, signedBlockDataBase64 string, name string) types.AttackerResponse {
	signedBlock, err := common.Base64ToSignedCapellaBlock(signedBlockDataBase64)
	if err != nil {
		return types.AttackerResponse{
			Cmd:    types.CMD_NULL,
			Result: signedBlockDataBase64,
		}
	}
	result := types.AttackerResponse{
		Cmd:    types.CMD_NULL,
		Result: signedBlockDataBase64,
	}

	if t, find := findMaxLevelStrategy(s.b.GetInternalSlotStrategy(), int64(slot)); find {
		action := t.Actions[name]
		if action != nil {
			r := action.RunAction(s.b, int64(slot), pubkey, signedBlock)
			result.Cmd = r.Cmd
			if newBlockBase64, err := common.SignedCapellaBlockToBase64(signedBlock); err != nil {
				log.WithError(err).WithFields(log.Fields{
					"slot":   slot,
					"action": name,
				}).Error("marshal to block failed")
			} else {
				result.Result = newBlockBase64
			}
		}
	}
	log.WithFields(log.Fields{
		"cmd":    result.Cmd,
		"slot":   slot,
		"action": name,
	}).Info("exit todoActionsWithBlock")

	return result
}
