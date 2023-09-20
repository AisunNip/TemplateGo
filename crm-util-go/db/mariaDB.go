package db

import (
	"context"
	"crm-util-go/errorcode"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var lockInitMariaDBPool sync.Mutex

// Database connection pool
var mariaPool *sql.DB
var countFailMaria int

func (dbPool DBPool) initMariaDBPool(transID string) *sql.DB {

	lockInitMariaDBPool.Lock()
	defer lockInitMariaDBPool.Unlock()

	var err error

	if mariaPool != nil {
		if countFailMaria < DbMaxFailTimes {
			return mariaPool
		}

		err = mariaPool.PingContext(context.Background())

		if err != nil {
			mariaPool.Close()
		} else {
			countFailMaria = 0
			return mariaPool
		}
	}

	mariaPool, err = sql.Open("mysql", dbPool.DataSourceName)

	if err != nil {
		dbPool.Logger.Error(transID, "Can not initial Maria connection pool", err)
	}

	if mariaPool != nil {
		mariaPool.SetMaxOpenConns(dbPool.MaxOpenConns)
		mariaPool.SetMaxIdleConns(dbPool.MaxIdleConns)
		mariaPool.SetConnMaxLifetime(dbPool.MaxLifetime)
		countFailMaria = 0
		dbPool.Logger.Info(transID, "Init MariaDBPool success")
	}

	return mariaPool
}

func (dbPool DBPool) GetMariaDBPool(transID string) (*sql.DB, error) {
	var err error

	if mariaPool == nil {
		mariaPool = dbPool.initMariaDBPool(transID)
	}

	dbPool.Logger.Debug(transID, fmt.Sprintf("Maria DB Stat: %+v", mariaPool.Stats()))

	err = mariaPool.PingContext(context.Background())

	if err != nil {
		dbPool.Logger.Error(transID, "Can not verify a connection to Maria DB because "+err.Error(), err)

		countFailMaria++

		if countFailMaria > DbMaxFailTimes {
			mariaPool = dbPool.initMariaDBPool(transID)
		}
	} else {
		countFailMaria = 0
	}

	return mariaPool, err
}

func (dbPool DBPool) CloseMariaDBPool() error {
	if mariaPool != nil {
		return mariaPool.Close()
	}

	return nil
}

func (dbPool DBPool) CreatePreparedStatementMaria(transID string, sql string, errCode errorcode.CrmErrorCode) (*sql.Stmt, *errorcode.CrmErrorCodeResp) {
	oraclePool, err := dbPool.GetMariaDBPool(transID)

	if err != nil {
		crmErrorCodeResp := errCode.GenerateAppError(dbPool.Logger.ApplicationName, "Get Maria database pool error: "+err.Error())
		return nil, &crmErrorCodeResp
	}

	prepare, err := oraclePool.Prepare(sql)
	if err != nil {
		crmErrorCodeResp := errCode.GenerateAppError(dbPool.Logger.ApplicationName, "Create a prepared statement error: "+err.Error())
		return nil, &crmErrorCodeResp
	}

	return prepare, nil
}
