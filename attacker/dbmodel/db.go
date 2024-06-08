package dbmodel

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/config"
)

func DbInit(dbconf config.MysqlConfig) {
	// Set up database
	datasource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", dbconf.User, dbconf.Passwd, dbconf.Host, dbconf.Port, dbconf.DbName)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", datasource)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}
	orm.RegisterModel(new(BlockReward))
	orm.RegisterModel(new(ChainReorg))
	orm.RunSyncdb("default", true, true)
}
