package types

type ActionPoint string

var allActionPoints = make(map[ActionPoint]bool)

const (
	AttestBeforeBroadCast     ActionPoint = "AttestBeforeBroadCast"
	AttestAfterBroadCast      ActionPoint = "AttestAfterBroadCast"
	AttestBeforeSign          ActionPoint = "AttestBeforeSign"
	AttestAfterSign           ActionPoint = "AttestAfterSign"
	AttestBeforePropose       ActionPoint = "AttestBeforePropose"
	AttestAfterPropose        ActionPoint = "AttestAfterPropose"
	BlockDelayForReceiveBlock ActionPoint = "BlockDelayForReceiveBlock"
	BlockBeforeBroadCast      ActionPoint = "BlockBeforeBroadCast"
	BlockAfterBroadCast       ActionPoint = "BlockAfterBroadCast"
	BlockBeforeSign           ActionPoint = "BlockBeforeSign"
	BlockAfterSign            ActionPoint = "BlockAfterSign"
	BlockBeforePropose        ActionPoint = "BlockBeforePropose"
	BlockAfterPropose         ActionPoint = "BlockAfterPropose"
)

func init() {
	allActionPoints[AttestBeforeBroadCast] = true
	allActionPoints[AttestAfterBroadCast] = true
	allActionPoints[AttestBeforeSign] = true
	allActionPoints[AttestAfterSign] = true
	allActionPoints[AttestBeforePropose] = true
	allActionPoints[AttestAfterPropose] = true
	allActionPoints[BlockDelayForReceiveBlock] = true
	allActionPoints[BlockBeforeBroadCast] = true
	allActionPoints[BlockAfterBroadCast] = true
	allActionPoints[BlockBeforeSign] = true
	allActionPoints[BlockAfterSign] = true
	allActionPoints[BlockBeforePropose] = true
	allActionPoints[BlockAfterPropose] = true
}
func CheckActionPointExist(action string) bool {
	_, ok := allActionPoints[ActionPoint(action)]
	return ok
}
