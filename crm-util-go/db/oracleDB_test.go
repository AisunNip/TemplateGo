package db

import (
	"context"
	"crm-util-go/common"
	"crm-util-go/errorcode"
	"crm-util-go/logging"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

type RespLatestAssetRootDB struct {
	Code           string  `json:"code,omitempty"`
	Msg            string  `json:"msg,omitempty"`
	TransID        string  `json:"transID,omitempty"`
	BillPayFlag    *string `json:"billPayFlag,omitempty"`
	IDNo           *string `json:"idNo"`
	BusIDNo        *string `json:"busIdNo"`
	CaAccountType  *string `json:"caAccountType"`
	StatusPriority *string `json:"statusPriority"`
	Age            *int    `json:"age"`
}

var assetLogger *logging.PatternLogger
var assetErrCode errorcode.CrmErrorCode
var oracleDBPool DBPool

func initOracle() {
	assetLogger = logging.InitInboundLogger("crm-util-go", logging.CrmDatabase)
	// assetLogger.EnableFileLogger("D:/", "AssetDao")

	// https://godror.github.io/godror/doc/connection.html
	// Initial database pool
	/*
			oracleDBPool.DataSourceName = `user="ccbcdv" password="ccbcdv#123"
		connectString="(DESCRIPTION=(ADDRESS_LIST = (ADDRESS = (PROTOCOL = TCP)(HOST = 172.19.216.26)(PORT = 1534)))(CONNECT_DATA = (SERVICE_NAME = CRMSIT02)))"
		poolMaxSessions=100 poolMinSessions=10 poolSessionMaxLifetime=10s poolWaitTimeout=40s standaloneConnection=1
		`
	*/

	// https://github.com/mattn/go-oci8
	// DataSourceName = user/password@host:port/sid

	// oracleDBPool.DataSourceName = "ccbcdv/ccbcdv#234@(DESCRIPTION=(ADDRESS_LIST=(LOAD_BALANCE=on)(FAILOVER=ON)(ADDRESS=(PROTOCOL=TCP)(HOST=172.19.190.148)(PORT=1555))(ADDRESS=(PROTOCOL=TCP)(HOST=172.19.190.157)(PORT=1555)))(CONNECT_DATA=(SERVER=DEDICATED)(SERVICE_NAME=CRMOLPRD)))"
	// oracleDBPool.DataSourceName = "ccbcdv/ccbcdv#123@172.19.216.26:1534/CRMSIT02"
	// oracleDBPool.DataSourceName = "ccbcdv/ccbcdv#123@(DESCRIPTION=(ADDRESS_LIST = (ADDRESS = (PROTOCOL = TCP)(HOST = 172.19.216.26)(PORT = 1534)))(CONNECT_DATA = (SERVICE_NAME = CRMSIT02)))"
	// oracleDBPool.DataSourceName = "TRUAPP9/TRUAPP9@ccbdbts1:1565/TEST02"

	oracleDBPool.DataSourceName = "ccbcdv/ccbcdv#123@(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=172.19.192.73)(PORT=1532))(CONNECT_DATA=(SERVER=DEDICATED)(SERVICE_NAME=CRMUAT01)))"
	oracleDBPool.MaxOpenConns = 50
	oracleDBPool.MaxIdleConns = 5
	oracleDBPool.Logger = assetLogger

	// Init Error Code
	errorcode.InitConfig("../config")
	assetErrCode.SystemCode = "CIB"
	assetErrCode.ModuleCode = "AS"
}

func TestGenerateSQLPagingOracle(t *testing.T) {
	sqlPaging := oracleDBPool.GenerateSQLPagingOracle("SELECT SYSDATE FROM DUAL", "", 1, 10)
	fmt.Println(sqlPaging)
}

func TestQueryOracle(t *testing.T) {
	transID := common.NewUUID()

	initOracle()
	defer oracleDBPool.CloseOracleDBPool()

	serviceID := "9166003560"
	assetRootResp := getLatestAssetRootDAO(transID, serviceID)
	fmt.Println(fmt.Sprintf("TransId: %s, AssetRootResp: %#v", transID, assetRootResp))

	if assetRootResp.Code != "0" {
		t.Errorf("expecting result to be Code=0 but %s", assetRootResp.Code)
	}
}

var sqlQueryLatestAssetRoot = `SELECT DATA1.* FROM
  (SELECT BILLPAY_FLG, CONID.ATTRIB_03 AS ID_NUM, CA.X_ID_NUM AS BUS_ID_NUM, CA.OU_TYPE_CD AS CA_ACCOUNT_TYPE,
    ROW_NUMBER() OVER(PARTITION BY AST.SERIAL_NUM ORDER BY DECODE(AST.STATUS_CD, 'Inactive', 2, 'Cancelled', 2, 'Terminate', 2, 1) ASC, AST.X_START_DT DESC) AS STATUS_PRIORITY
  FROM SIEBEL.S_ASSET AST
  LEFT JOIN SIEBEL.S_ORG_EXT CA ON AST.OWNER_ACCNT_ID = CA.ROW_ID
  LEFT JOIN SIEBEL.S_CONTACT CON ON CA.LOC = CON.PERSON_UID
  LEFT JOIN SIEBEL.S_CONTACT_XM CONID ON CON.X_PR_ID = CONID.ROW_ID
  WHERE AST.PAR_ASSET_ID IS NULL
  AND AST.SERIAL_NUM = :1) DATA1
WHERE STATUS_PRIORITY = 1`

func getLatestAssetRootDAO(transID string, serviceID string) RespLatestAssetRootDB {
	startDT := assetLogger.LogRequestDBClient(transID)

	var resp RespLatestAssetRootDB
	resp.TransID = transID

	prepare, crmErrorCodeResp := oracleDBPool.CreatePreparedStatementOracle(transID, sqlQueryLatestAssetRoot, assetErrCode)

	if crmErrorCodeResp != nil {
		assetLogger.LogResponseDBClient(transID, crmErrorCodeResp.ErrorCode, startDT)

		resp.Code = crmErrorCodeResp.ErrorCode
		resp.Msg = crmErrorCodeResp.ErrorMessage
		return resp
	}

	defer prepare.Close()

	var rows *sql.Rows
	// rows, err = prepare.Query(serviceID)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelFunc()
	rows, err := prepare.QueryContext(ctx, serviceID)

	if err != nil {
		assetLogger.Error(transID, "Query DB Error: "+err.Error())
		crmErrorCodeResp := assetErrCode.GenerateAppError("crm-util-go", "Query DB Error: "+err.Error())

		resp.Code = crmErrorCodeResp.ErrorCode
		resp.Msg = crmErrorCodeResp.ErrorMessage

		assetLogger.LogResponseDBClient(transID, resp.Code, startDT)
		return resp
	}
	defer rows.Close()

	hasResult := false

	for rows.Next() {
		rows.Scan(
			&resp.BillPayFlag,
			&resp.IDNo,
			&resp.BusIDNo,
			&resp.CaAccountType,
			&resp.StatusPriority,
		)
		hasResult = true
	}

	if hasResult {
		resp.Code = "0"
		resp.Msg = "Success"
	} else {
		crmErrorCodeResp := assetErrCode.GenerateDataNotFound("AssetRoot", "CRM")

		resp.Code = crmErrorCodeResp.ErrorCode
		resp.Msg = crmErrorCodeResp.ErrorMessage
	}

	assetLogger.LogResponseDBClient(transID, resp.Code, startDT)
	return resp
}
