package types

import "encoding/json"

type AttackerCommand int

const (
	CMD_NULL AttackerCommand = iota
	CMD_CONTINUE
	CMD_RETURN
	CMD_ABORT
	CMD_SKIP
	CMD_ROLE_TO_NORMAL   // 角色转换为普通节点
	CMD_ROLE_TO_ATTACKER // 角色转换为攻击者
	CMD_EXIT
	CMD_UPDATE_STATE
)

type RoleType int

const (
	NormalRole RoleType = iota
	AttackerRole
)

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
