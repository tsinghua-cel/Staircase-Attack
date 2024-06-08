package reward

import (
	"encoding/csv"
	"errors"
	"github.com/astaxie/beego/orm"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/beaconapi"
	"github.com/tsinghua-cel/attacker-service/dbmodel"
	"os"
	"strconv"
)

func GetRewardsToMysql(gwEndpoint string) error {
	client := beaconapi.NewBeaconGwClient(gwEndpoint)
	slots_per_epoch, err := client.GetIntConfig(beaconapi.SLOTS_PER_EPOCH)
	if err != nil {
		log.WithError(err).Error("GetRewardsToMysql get chain config failed")
		return err
	}
	latestHeader, err := client.GetLatestBeaconHeader()
	if err != nil {
		return err
	}

	latestSlot, _ := strconv.ParseInt(latestHeader.Header.Message.Slot, 10, 64)
	latestEpoch := latestSlot / int64(slots_per_epoch)

	curMaxEpoch := dbmodel.GetMaxEpoch()
	epochNumber := curMaxEpoch + 1
	if curMaxEpoch < 0 {
		epochNumber = 0
	}
	o := orm.NewOrm()

	//  开始事务
	if err = o.Begin(); err != nil {
		log.WithError(err).Error("GetRewardsToMysql orm begin failed")
		return err
	}
	repo := dbmodel.NewBlockRewardRepository(o)
	log.WithFields(log.Fields{
		"epochNumber": epochNumber,
		"latestEpoch": latestEpoch,
	}).Debug("GetRewardsToMysql")

	for epochNumber <= (latestEpoch - 2) {
		totalRewards, err := client.GetAllValReward(int(epochNumber))
		if err != nil {
			return err
		}

		for _, totalReward := range totalRewards {
			valIdx, _ := strconv.ParseInt(totalReward.ValidatorIndex, 10, 64)
			headAmount, _ := strconv.ParseInt(totalReward.Head, 10, 64)
			targetAmount, _ := strconv.ParseInt(totalReward.Target, 10, 64)
			record := &dbmodel.BlockReward{
				Epoch:          epochNumber,
				ValidatorIndex: int(valIdx),
				HeadAmount:     headAmount,
				TargetAmount:   targetAmount,
			}
			if err = repo.Create(record); err != nil {
				o.Rollback()
				return errors.New("insert failed")
			}
		}
		epochNumber++
	}
	if err = o.Commit(); err != nil {
		return errors.New("commit failed")
	}
	return nil
}

func GetRewards(gwEndpoint string, output string) error {
	bakfile := output + ".bak"
	file, err := os.Create(bakfile)
	if err != nil {
		return err
	}
	succeed := false
	defer func() {
		file.Close()
		if succeed {
			os.Rename(bakfile, output)
		}
	}()
	client := beaconapi.NewBeaconGwClient(gwEndpoint)

	slots_per_epoch, err := client.GetIntConfig(beaconapi.SLOTS_PER_EPOCH)
	if err != nil {
		// todo: add log
		return err
	}
	latestHeader, err := client.GetLatestBeaconHeader()
	if err != nil {
		return err
	}

	latestSlot, _ := strconv.ParseInt(latestHeader.Header.Message.Slot, 10, 64)
	latestEpoch := latestSlot / int64(slots_per_epoch)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Epoch", "Validator Index", "Head", "Target", "Source", "Inclusion Delay", "Inactivity"})

	epochNumber := int64(0)

	for epochNumber <= (latestEpoch - 2) {
		totalRewards, err := client.GetAllValReward(int(epochNumber))
		if err != nil {
			return err
		}
		for _, totalReward := range totalRewards {
			writer.Write([]string{strconv.FormatInt(epochNumber, 10), totalReward.ValidatorIndex, totalReward.Head, totalReward.Target, totalReward.Source, totalReward.InclusionDelay, totalReward.Inactivity})
		}

		epochNumber++

	}
	succeed = true
	return err
}
