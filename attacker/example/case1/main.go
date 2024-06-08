package main

import (
	"github.com/tsinghua-cel/attacker-service/config"
	"github.com/tsinghua-cel/attacker-service/server"
	"sync"
)

func main() {
	rpcServer := server.NewServer(config.GetConfig(), NewPluginCaseV1())
	rpcServer.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)

	wg.Wait()
}
