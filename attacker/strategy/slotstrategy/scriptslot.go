package slotstrategy

import "strings"

// todo:implement script config in strategy, parse it to a function and can run.
type Runable interface {
	Run() int64
}

type ScriptSlot struct {
	toRun Runable
}

func (s ScriptSlot) Compare(slot int64) int {
	cSlot := int64(0)
	if s.toRun != nil {
		cSlot = s.toRun.Run()
	}
	if cSlot > slot {
		return 1
	}
	if cSlot < slot {
		return -1
	}
	return 0
}

func parseStringToRunable(script string) Runable {
	if strings.HasPrefix(script, "go:") {
		// todo: implement go script parse.
	} else if strings.HasPrefix(script, "lua:") {
		// todo: implement lua script parse.
	} else {
		// unsupported
	}
	return nil

}
