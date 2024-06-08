package common

type SlotTool struct {
	SlotsPerEpoch int
}

func (s SlotTool) SlotToEpoch(slot int64) int64 {
	return slot / int64(s.SlotsPerEpoch)
}

func (s SlotTool) EpochEnd(epoch int64) int64 {
	return (epoch+1)*int64(s.SlotsPerEpoch) - 1
}

func (s SlotTool) EpochStart(epoch int64) int64 {
	return epoch * int64(s.SlotsPerEpoch)
}
