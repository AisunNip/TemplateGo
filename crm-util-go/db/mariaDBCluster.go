package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/benthor/clustersql"
	"github.com/go-sql-driver/mysql"
	"sync"
)

var lockInitMariaDBPoolCluster sync.Mutex

// Database connection pool
var mariaPoolCluster *sql.DB
var countFailMariaCluster int

func (dbPool DBPoolCluster) initMariaDBPoolCluster(transID string) *sql.DB {

	lockInitMariaDBPoolCluster.Lock()
	defer lockInitMariaDBPoolCluster.Unlock()

	var err error

	if mariaPoolCluster != nil {
		if countFailMariaCluster < DbMaxFailTimes {
			return mariaPoolCluster
		}

		err = mariaPoolCluster.PingContext(context.Background())

		if err != nil {
			mariaPoolCluster.Close()
		} else {
			countFailMariaCluster = 0
			return mariaPoolCluster
		}
	}

	mysqlDriver := mysql.MySQLDriver{}

	clusterDriver := clustersql.NewDriver(mysqlDriver)

	// clusterDriver.AddNode("galera1", "user:password@tcp(dbhost1:3306)/db_name")
	for _, dbNode := range dbPool.DBNode {
		clusterDriver.AddNode(dbNode.NodeName, dbNode.DataSourceName)
	}

	sql.Register(dbPool.RegisterName, clusterDriver)

	mariaPoolCluster, err = sql.Open(dbPool.RegisterName, "NoDSN")

	if err != nil {
		dbPool.Logger.Error(transID, "Can not initial Maria connection pool", err)
	}

	if mariaPoolCluster != nil {
		mariaPoolCluster.SetMaxOpenConns(dbPool.MaxOpenConns)
		mariaPoolCluster.SetMaxIdleConns(dbPool.MaxIdleConns)
		mariaPoolCluster.SetConnMaxLifetime(dbPool.MaxLifetime)
		countFailMariaCluster = 0
		dbPool.Logger.Info(transID, "Init MariaDBPoolCluster success")
	}

	return mariaPoolCluster
}

func (dbPool DBPoolCluster) GetMariaDBPoolCluster(transID string) (*sql.DB, error) {
	var err error

	if mariaPoolCluster == nil {
		mariaPoolCluster = dbPool.initMariaDBPoolCluster(transID)
	}

	dbPool.Logger.Debug(transID, fmt.Sprintf("Maria DB Stat: %+v", mariaPoolCluster.Stats()))

	err = mariaPoolCluster.PingContext(context.Background())

	if err != nil {
		dbPool.Logger.Error(transID, "Can not verify a connection to Maria DB because "+err.Error(), err)

		countFailMariaCluster++

		if countFailMariaCluster > DbMaxFailTimes {
			mariaPoolCluster = dbPool.initMariaDBPoolCluster(transID)
		}
	} else {
		countFailMariaCluster = 0
	}

	return mariaPoolCluster, err
}
