package slotstrategy

import (
	"errors"
	"fmt"
	"github.com/tsinghua-cel/attacker-service/plugins"
	"github.com/tsinghua-cel/attacker-service/types"
	"strconv"
)

type SlotIns interface {
	// if slotIns < slot, return -1
	// if slotIns == slot, return 0
	// if slotIns > slot, return 1
	Compare(slot int64) int
}

type ActionIns interface {
	RunAction(backend types.ServiceBackend, slot int64, pubkey string, params ...interface{}) plugins.PluginResponse
}

type InternalSlotStrategy struct {
	Slot    SlotIns              `json:"slot"`
	Level   int                  `json:"level"`
	Actions map[string]ActionIns `json:"actions"`
}

func ParseToInternalSlotStrategy(backend types.ServiceBackend, strategy []types.SlotStrategy) ([]InternalSlotStrategy, error) {
	is := make([]InternalSlotStrategy, len(strategy))
	for i, s := range strategy {
		is[i].Level = s.Level
		if n, err := strconv.ParseInt(s.Slot, 10, 64); err == nil {
			is[i].Slot = NumberSlot(n)
		} else {
			calc, err := GetFunctionSlot(backend, s.Slot)
			if err != nil {
				return nil, err
			}
			is[i].Slot = FunctionSlot{calcFunc: calc}
		}
		is[i].Actions = make(map[string]ActionIns)
		for point, action := range s.Actions {
			if types.CheckActionPointExist(point) == false {
				return nil, errors.New(fmt.Sprintf("action point %s not exist", point))
			}
			actionDo, err := GetFunctionAction(backend, action)
			if err != nil {
				return nil, err
			}
			is[i].Actions[point] = FunctionAction{doFunc: actionDo}
		}
	}
	//log.Printf("parsed internal slot strategy is %v\n", is)
	return is, nil
}
