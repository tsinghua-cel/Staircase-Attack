package types

import "encoding/json"

type AttackerResponse struct {
	Cmd    AttackerCommand `json:"cmd"`
	Result string          `json:"result"`
}

type ClientInfo struct {
	UUID           string `json:"uuid"`
	ValidatorIndex int    `json:"validatorIndex"`
}

func ToClientInfo(cliInfo string) ClientInfo {
	var cinfo ClientInfo
	json.Unmarshal([]byte(cliInfo), &cinfo)
	return cinfo
}

type SlotStateRoot struct {
	Root string `json:"root"`
}
