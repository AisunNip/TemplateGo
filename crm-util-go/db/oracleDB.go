package db

import (
	"context"
	"crm-util-go/common"
	"crm-util-go/errorcode"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"strings"
	"sync"
)

// [username/[password]@]host[:port][/service_name][?param1=value1&...&paramN=valueN]
// Supported parameters are:
//
// loc - the time location for reading timestamp (without time zone). Defaults to UTC
// Note that writing a timestamp (without time zone) just truncates the time zone.
//
// isolation - the isolation level that can be set to: READONLY, SERIALIZABLE, or DEFAULT
//
// prefetch_rows - the number of top level rows to be prefetched. Defaults to 0. A 0 means unlimited rows.
//
// prefetch_memory - the max memory for top level rows to be prefetched. Defaults to 4096. A 0 means unlimited memory.
//
// questionph - when true, enables question mark placeholders. Defaults to false. (uses strconv.ParseBool to check for true)

var lockInitOracleDBPool sync.Mutex
var oraclePool *sql.DB
var countFailOracle int

func (dbPool DBPool) initOracleDBPool(transID string) *sql.DB {

	lockInitOracleDBPool.Lock()
	defer lockInitOracleDBPool.Unlock()

	var err error

	if oraclePool != nil {
		if countFailOracle < DbMaxFailTimes {
			return oraclePool
		}

		err = oraclePool.PingContext(context.Background())

		if err != nil {
			oraclePool.Close()
		} else {
			countFailOracle = 0
			return oraclePool
		}
	}

	oraclePool, err = sql.Open("oci8", dbPool.DataSourceName)

	if err != nil {
		dbPool.Logger.Error(transID, "Can not initial Oracle connection pool", err)
	}

	if oraclePool != nil {
		oraclePool.SetMaxOpenConns(dbPool.MaxOpenConns)
		oraclePool.SetMaxIdleConns(dbPool.MaxIdleConns)
		countFailOracle = 0
		dbPool.Logger.Info(transID, "Init OracleDBPool success")
	}

	return oraclePool
}

func (dbPool DBPool) GetOracleDBPool(transID string) (*sql.DB, error) {
	var err error

	if oraclePool == nil {
		oraclePool = dbPool.initOracleDBPool(transID)
	}

	dbPool.Logger.Debug(transID, fmt.Sprintf("Oracle DB Stat: %+v", oraclePool.Stats()))

	err = oraclePool.PingContext(context.Background())

	if err != nil {
		dbPool.Logger.Error(transID, "Can not verify a connection to Oracle DB because "+err.Error(), err)

		countFailOracle++

		if countFailOracle > DbMaxFailTimes {
			oraclePool = dbPool.initOracleDBPool(transID)
		}
	} else {
		countFailOracle = 0
	}

	return oraclePool, err
}

func (dbPool DBPool) CloseOracleDBPool() error {
	if oraclePool != nil {
		return oraclePool.Close()
	}

	return nil
}

func (dbPool DBPool) BeginTransactionOracle(transID string, errCode errorcode.CrmErrorCode) (*sql.Tx, *errorcode.CrmErrorCodeResp) {
	oraclePool, err := dbPool.GetOracleDBPool(transID)

	if err != nil {
		crmErrorCodeResp := errCode.GenerateAppError(dbPool.Logger.ApplicationName, "Get Oracle database pool error: "+err.Error())
		return nil, &crmErrorCodeResp
	}

	tx, err := oraclePool.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {
		crmErrorCodeResp := errCode.GenerateAppError(dbPool.Logger.ApplicationName, "Oracle start a transaction error: "+err.Error())
		return nil, &crmErrorCodeResp
	}

	/*
		_, execErr := tx.Exec(sql, bindingValue)
		tx.Commit()
		tx.Rollback()
	*/

	return tx, nil
}

func (dbPool DBPool) CreatePreparedStatementOracle(transID string, sql string, errCode errorcode.CrmErrorCode) (*sql.Stmt, *errorcode.CrmErrorCodeResp) {
	oraclePool, err := dbPool.GetOracleDBPool(transID)

	if err != nil {
		crmErrorCodeResp := errCode.GenerateAppError(dbPool.Logger.ApplicationName, "Get Oracle database pool error: "+err.Error())
		return nil, &crmErrorCodeResp
	}

	prepare, err := oraclePool.Prepare(sql)
	if err != nil {
		crmErrorCodeResp := errCode.GenerateAppError(dbPool.Logger.ApplicationName, "Create a prepared statement error: "+err.Error())
		return nil, &crmErrorCodeResp
	}

	return prepare, nil
}

func (dbPool DBPool) GenerateSQLPagingOracle(sql string, condition string, pageNo int, pageSize int) string {
	startRecord := (pageNo-1)*pageSize + 1
	endRecord := pageNo * pageSize

	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("SELECT DATA2.* ")
	sqlBuilder.WriteString("FROM (SELECT ROWNUM MYNUM, DATA1.* ")
	sqlBuilder.WriteString("FROM (")
	sqlBuilder.WriteString(sql)
	sqlBuilder.WriteString(") DATA1 ")
	sqlBuilder.WriteString("WHERE ")

	if len(condition) > 0 {
		sqlBuilder.WriteString(condition)
		sqlBuilder.WriteString(" AND ROWNUM <= ")
		sqlBuilder.WriteString(common.IntToString(endRecord))
	} else {
		sqlBuilder.WriteString("ROWNUM <= ")
		sqlBuilder.WriteString(common.IntToString(endRecord))
	}

	sqlBuilder.WriteString(") DATA2 WHERE MYNUM >= ")
	sqlBuilder.WriteString(common.IntToString(startRecord))

	return sqlBuilder.String()
}
