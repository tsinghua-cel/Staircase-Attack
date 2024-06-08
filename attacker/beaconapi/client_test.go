package beaconapi

import (
	"fmt"
	"testing"
)

func TestBeaconGwClient_MonitorReorgEvent(t *testing.T) {
	beaconGwClient := NewBeaconGwClient("13.41.176.56:15000")
	//ch := beaconGwClient.MonitorReorgEvent()
	//for {
	//	select {
	//	case reorgEvent := <-ch:
	//		fmt.Sprintf("reorg event: %v", reorgEvent)
	//		t.Log(reorgEvent)
	//	}
	//}
	header, err := beaconGwClient.GetBlockHeaderById("0x410152ba2011e946c3d305d38d842548d5f68281c07e30a18b2e8927db4346a2")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("header slot and proposer is ", header.Header.Message.Slot, header.Header.Message.ProposerIndex)

}
