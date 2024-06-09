package apis

type SlotI interface {
	SlotsPerEpoch() int
}

type SlotTool struct {
	sloti SlotI
}

func (s SlotTool) SlotToEpoch(slot int) int {
	return slot / s.sloti.SlotsPerEpoch()
}

func (s SlotTool) EpochEnd(epoch int) int {
	return (epoch+1)*s.sloti.SlotsPerEpoch() - 1
}

func (s SlotTool) EpochStart(epoch int) int {
	return epoch * s.sloti.SlotsPerEpoch()
}
