package beaconapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
)

// test GetValidators
func TestGetValidators(t *testing.T) {
	endpoint := "52.221.177.10:34000" // grpc endpoint
	pubks, err := GetValidators(endpoint)
	if err != nil {
		t.Fatalf("get validators failed err:%s", err)
	}
	fmt.Printf("get validators %v\n", pubks)
}

func TestGetReward(t *testing.T) {
	endpoint := "52.221.177.10:33500" // grpc gateway endpoint
	valIdxs := []int{1, 2, 3, 4, 5}
	client := NewBeaconGwClient(endpoint)
	res, err := client.GetValReward(1, valIdxs)
	if err != nil {
		t.Fatalf("get reward failed err:%s", err)
	}
	fmt.Printf("get specific reward res:%s\n", res)
}

func TestGetAllReward(t *testing.T) {
	endpoint := "52.221.177.10:33500" // grpc gateway endpoint
	client := NewBeaconGwClient(endpoint)
	res, err := client.GetAllValReward(1)
	if err != nil {
		t.Fatalf("get reward failed err:%s", err)
	}
	fmt.Printf("get all reward res:%s\n", res)
}

func TestGetConfig(t *testing.T) {
	endpoint := "52.221.177.10:33500" // grpc gateway endpoint
	client := NewBeaconGwClient(endpoint)
	epoch, err := client.GetIntConfig(SLOTS_PER_EPOCH)
	if err != nil {
		t.Fatalf("get epoch config failed err:%s", err)
	}
	fmt.Printf("get epoch :%d\n", epoch)
}

func TestGetLatestBeaconHeader(t *testing.T) {
	endpoint := "52.221.177.10:33500" // grpc gateway endpoint
	client := NewBeaconGwClient(endpoint)

	header, err := client.GetLatestBeaconHeader()
	if err != nil {
		t.Fatalf("get latest header failed err:%s", err)
	}
	fmt.Printf("get latest header.slot :%s\n", header.Header.Message.Slot)

}

func TestGetAllAttestDuties(t *testing.T) {
	endpoint := "52.221.177.10:14000" // grpc gateway endpoint
	client := NewBeaconGwClient(endpoint)
	duties, err := client.GetProposerDuties(1)
	//duties, err := client.GetCurrentEpochProposerDuties()
	if err != nil {
		t.Fatalf("get proposer duties failed err:%s", err)
	}

	latestSlotWithAttacker := int64(-1)
	for _, duty := range duties {
		dutySlot, _ := strconv.ParseInt(duty.Slot, 10, 64)
		dutyValIdx, _ := strconv.Atoi(duty.ValidatorIndex)
		fmt.Printf("slot=%d, validx =%d\n", dutySlot, dutyValIdx)

		if dutyValIdx <= 31 && dutySlot > latestSlotWithAttacker {
			latestSlotWithAttacker = dutySlot
			fmt.Printf("update latestSlotWithAttacker=%d,\n", dutySlot)
		}
	}

	for _, duty := range duties {
		d, _ := json.Marshal(duty)
		fmt.Printf("get attest duty :%s\n", string(d))
	}
}
