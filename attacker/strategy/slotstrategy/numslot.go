package slotstrategy

type NumberSlot int64

func (n NumberSlot) Compare(slot int64) int {
	if int64(n) > slot {
		return 1
	}
	if int64(n) < slot {
		return -1
	}
	return 0
}
