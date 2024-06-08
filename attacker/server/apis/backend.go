package apis

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/plugins"
	"github.com/tsinghua-cel/attacker-service/rpc"
	"github.com/tsinghua-cel/attacker-service/strategy/slotstrategy"
	"github.com/tsinghua-cel/attacker-service/types"
)

// Backend interface provides the common API services (that are provided by
// both full and light clients) with access to necessary functions.
type Backend interface {
	SomeNeedBackend() bool
	GetInternalSlotStrategy() []slotstrategy.InternalSlotStrategy
	types.ExecuteBackend
	types.BeaconBackend
	types.StrategyBackend
	types.CacheBackend
}

func GetAPIs(apiBackend Backend, plugin plugins.AttackerPlugin) []rpc.API {
	return []rpc.API{
		{
			Namespace: "admin",
			Service:   NewAdminAPI(apiBackend, plugin),
		},
		{
			Namespace: "block",
			Service:   NewBlockAPI(apiBackend, plugin),
		},
		{
			Namespace: "attest",
			Service:   NewAttestAPI(apiBackend, plugin),
		},
	}
}

func pluginContext(backend types.ServiceBackend) plugins.PluginContext {
	return plugins.PluginContext{
		Backend: backend,
		Context: context.Background(),
		Logger:  log.WithField("module", "attacker-service"),
	}
}
