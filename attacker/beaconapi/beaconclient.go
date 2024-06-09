package beaconapi

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	log "github.com/sirupsen/logrus"
	"strconv"
)

const (
	SLOTS_PER_EPOCH  = "SLOTS_PER_EPOCH"
	SECONDS_PER_SLOT = "SECONDS_PER_SLOT"
)

type BeaconGwClient struct {
	endpoint string
	config   map[string]string
}

func NewBeaconGwClient(endpoint string) *BeaconGwClient {
	return &BeaconGwClient{
		endpoint: endpoint,
		config:   make(map[string]string),
	}
}

func (b *BeaconGwClient) GetIntConfig(key string) (int, error) {
	config := b.GetBeaconConfig()
	if v, exist := config[key]; !exist {
		return 0, nil
	} else {
		return strconv.Atoi(v)
	}
}

func (b *BeaconGwClient) doGet(url string) (BeaconResponse, error) {
	resp, err := httplib.Get(url).Response()
	if err != nil {
		return BeaconResponse{}, err
	}
	defer resp.Body.Close()

	var response BeaconResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.WithError(err).Error("Error decoding response")
	}
	return response, nil
}

func (b *BeaconGwClient) doPost(url string, data []byte) (BeaconResponse, error) {
	resp, err := httplib.Post(url).Body(data).Response()
	if err != nil {
		return BeaconResponse{}, err
	}
	defer resp.Body.Close()

	var response BeaconResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.WithError(err).Error("Error decoding response")
	}
	return response, nil
}

func (b *BeaconGwClient) getBeaconConfig() (map[string]interface{}, error) {
	response, err := b.doGet(fmt.Sprintf("http://%s/eth/v1/config/spec", b.endpoint))

	config := make(map[string]interface{})
	err = json.Unmarshal(response.Data, &config)
	if err != nil {
		log.WithError(err).Error("unmarshal config data failed")
	}
	return config, nil
}

func (b *BeaconGwClient) GetBeaconConfig() map[string]string {
	if len(b.config) == 0 {
		config, err := b.getBeaconConfig()
		if err != nil {
			// todo: add log
			return nil
		}
		b.config = make(map[string]string)
		for key, v := range config {
			b.config[key] = v.(string)
		}
	}
	return b.config
}

func (b *BeaconGwClient) GetLatestBeaconHeader() (BeaconHeaderInfo, error) {
	response, err := b.doGet(fmt.Sprintf("http://%s/eth/v1/beacon/headers", b.endpoint))
	var headers = make([]BeaconHeaderInfo, 0)
	err = json.Unmarshal(response.Data, &headers)
	if err != nil {
		// todo: add log.
		return BeaconHeaderInfo{}, err
	}

	return headers[0], nil
}

// default grpc-gateway port is 3500
func (b *BeaconGwClient) GetAllValReward(epoch int) ([]TotalReward, error) {
	url := fmt.Sprintf("http://%s/eth/v1/beacon/rewards/attestations/%d", b.endpoint, epoch)
	response, err := b.doPost(url, []byte("[]"))
	var rewardInfo RewardInfo
	err = json.Unmarshal(response.Data, &rewardInfo)
	if err != nil {
		log.WithError(err).Error("unmarshal reward data failed")
		return nil, err
	}
	return rewardInfo.TotalRewards, err
}

func (b *BeaconGwClient) GetValReward(epoch int, valIdxs []int) (BeaconResponse, error) {
	url := fmt.Sprintf("http://%s/eth/v1/beacon/rewards/attestations/%d", b.endpoint, epoch)
	vals := make([]string, len(valIdxs))
	for i := 0; i < len(valIdxs); i++ {
		vals[i] = strconv.FormatInt(int64(valIdxs[i]), 10)
	}
	d, err := json.Marshal(vals)
	if err != nil {
		log.WithError(err).Error("get reward failed when marshal vals")
		return BeaconResponse{}, err
	}
	response, err := b.doPost(url, d)
	return response, err
}

// /eth/v1/validator/duties/proposer/:epoch
func (b *BeaconGwClient) GetProposerDuties(epoch int) ([]ProposerDuty, error) {
	url := fmt.Sprintf("http://%s/eth/v1/validator/duties/proposer/%d", b.endpoint, epoch)
	var duties = make([]ProposerDuty, 0)

	response, err := b.doGet(url)
	err = json.Unmarshal(response.Data, &duties)
	if err != nil {
		return []ProposerDuty{}, err
	}

	return duties, err
}

// POST /eth/v1/validator/duties/attester/:epoch
func (b *BeaconGwClient) GetAttesterDuties(epoch int, vals []int) ([]AttestDuty, error) {
	url := fmt.Sprintf("http://%s/eth/v1/validator/duties/attester/%d", b.endpoint, epoch)
	param := make([]string, len(vals))
	for i := 0; i < len(vals); i++ {
		param[i] = strconv.FormatInt(int64(vals[i]), 10)
	}
	paramData, _ := json.Marshal(param)
	var duties = make([]AttestDuty, 0)

	response, err := b.doPost(url, paramData)
	err = json.Unmarshal(response.Data, &duties)
	if err != nil {
		return []AttestDuty{}, err
	}
	return duties, err
}

func (b *BeaconGwClient) GetNextEpochProposerDuties() ([]ProposerDuty, error) {
	latestHeader, err := b.GetLatestBeaconHeader()
	if err != nil {
		return nil, err
	}
	slotPerEpoch, _ := b.GetIntConfig(SLOTS_PER_EPOCH)
	curSlot, _ := strconv.Atoi(latestHeader.Header.Message.Slot)
	epoch := curSlot / slotPerEpoch
	return b.GetProposerDuties(epoch + 1)
}

func (b *BeaconGwClient) GetCurrentEpochProposerDuties() ([]ProposerDuty, error) {
	latestHeader, err := b.GetLatestBeaconHeader()
	if err != nil {
		return nil, err
	}
	slotPerEpoch, _ := b.GetIntConfig(SLOTS_PER_EPOCH)
	curSlot, _ := strconv.Atoi(latestHeader.Header.Message.Slot)
	epoch := curSlot / slotPerEpoch
	return b.GetProposerDuties(epoch)
}

func (b *BeaconGwClient) GetCurrentEpochAttestDuties() ([]AttestDuty, error) {
	latestHeader, err := b.GetLatestBeaconHeader()
	if err != nil {
		return nil, err
	}
	slotPerEpoch, _ := b.GetIntConfig(SLOTS_PER_EPOCH)
	curSlot, _ := strconv.Atoi(latestHeader.Header.Message.Slot)
	epoch := curSlot / slotPerEpoch
	vals := make([]int, 64)
	for i := 0; i < len(vals); i++ {
		vals[i] = i
	}
	return b.GetAttesterDuties(epoch, vals)
}

func (b *BeaconGwClient) GetNextEpochAttestDuties() ([]AttestDuty, error) {
	latestHeader, err := b.GetLatestBeaconHeader()
	if err != nil {
		return nil, err
	}
	slotPerEpoch, _ := b.GetIntConfig(SLOTS_PER_EPOCH)
	curSlot, _ := strconv.Atoi(latestHeader.Header.Message.Slot)
	epoch := curSlot / slotPerEpoch
	vals := make([]int, 64)
	for i := 0; i < len(vals); i++ {
		vals[i] = i
	}
	return b.GetAttesterDuties(epoch+1, vals)
}
