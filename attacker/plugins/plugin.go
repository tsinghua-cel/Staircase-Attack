package plugins

import (
	"context"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/types"
)

type PluginContext struct {
	Context context.Context
	Backend types.ServiceBackend
	Logger  *logrus.Entry
}

type PluginResponse struct {
	Cmd    types.AttackerCommand
	Result interface{}
}

type AttackerPlugin interface {
	AttestBeforeBroadCast(ctx PluginContext, slot uint64) PluginResponse
	AttestAfterBroadCast(ctx PluginContext, slot uint64) PluginResponse
	AttestBeforeSign(ctx PluginContext, slot uint64, pubkey string, attestData *ethpb.AttestationData) PluginResponse
	AttestAfterSign(ctx PluginContext, slot uint64, pubkey string, attest *ethpb.Attestation) PluginResponse
	AttestBeforePropose(ctx PluginContext, slot uint64, pubkey string, attest *ethpb.Attestation) PluginResponse
	AttestAfterPropose(ctx PluginContext, slot uint64, pubkey string, attest *ethpb.Attestation) PluginResponse

	BlockDelayForReceiveBlock(ctx PluginContext, slot uint64) PluginResponse
	BlockBeforeBroadCast(ctx PluginContext, slot uint64) PluginResponse
	BlockAfterBroadCast(ctx PluginContext, slot uint64) PluginResponse
	BlockBeforeSign(ctx PluginContext, slot uint64, pubkey string, block *ethpb.SignedBeaconBlockCapella) PluginResponse
	BlockAfterSign(ctx PluginContext, slot uint64, pubkey string, block *ethpb.SignedBeaconBlockCapella) PluginResponse
	BlockBeforePropose(ctx PluginContext, slot uint64, pubkey string, block *ethpb.SignedBeaconBlockCapella) PluginResponse
	BlockAfterPropose(ctx PluginContext, slot uint64, pubkey string, block *ethpb.SignedBeaconBlockCapella) PluginResponse
}
