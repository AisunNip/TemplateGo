package test

import (
	"crm-util-go/common"
	"crm-util-go/db"
	"crm-util-go/logging"
	config "github.com/spf13/viper"
	"time"
)

var csdLogger *logging.PatternLogger

func initCassandra() {
	csdLogger = logging.InitInboundLogger("crm-util-go", logging.CrmDatabase)
	csdLogger.Level = logging.LEVEL_ALL

	// ############ Init Config ############
	config.SetConfigName("prod1")
	config.SetConfigType("yaml")
	config.AddConfigPath("./config")
	err := config.ReadInConfig()

	if err != nil {
		panic("Error load a configuration file")
	}
}

func TestCassandraProduction() {
	transID := common.NewUUID()

	initCassandra()

	hosts := config.GetStringSlice("crm.cassandra.cluster")
	username := config.GetString("crm.cassandra.username")
	password := config.GetString("crm.cassandra.password")
	keyspace := config.GetString("crm.cassandra.keyspace")

	csdPool := db.NewCassandraPool(hosts, username, password, keyspace, csdLogger)

	for true {
		val, err := csdPool.GetConfigVal(transID,"siebel.header.user")

		if err != nil {
			csdLogger.Info(transID, "Error " + err.Error())
		} else {
			csdLogger.Info(transID, "Success " + val)
		}

		time.Sleep(1 * time.Second)
	}
}