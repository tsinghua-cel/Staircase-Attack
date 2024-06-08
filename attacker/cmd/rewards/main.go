package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/reward"
)

var (
	noderpc    = flag.String("node", "127.0.0.1:8545", "execute node rpc addr")
	rewardfile = flag.String("output", "reward.csv", "output file for reward.")
)

func main() {
	flag.Parse()
	initLog()
	log.WithFields(log.Fields{
		"node":   *noderpc,
		"output": *rewardfile,
	}).Info("start get reward")

	err := reward.GetRewards(*noderpc, *rewardfile)
	if err != nil {
		log.WithError(err).Error("get reward failed")
	} else {
		log.Info("finish get reward")
	}
}

func initLog() {
	// standard setting
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05.000"})
}
