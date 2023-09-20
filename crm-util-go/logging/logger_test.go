package logging

import (
	"crm-util-go/common"
	"testing"
)

func TestLoggerAPIProviderMonitor(t *testing.T) {
	var appName = "crm-util-go"
	transID := common.NewUUID()

	// Ex. Rest API
	logger := InitInboundLogger(appName, CrmInbound)
	logger.Level = LEVEL_ALL
	//logger.EnableFileLogger("D:/", appName)

	startREST := logger.LogRequestRESTProvider(transID)
	logger.Info(transID, "TODO: API logic")
	restRespCode := "0"
	logger.LogResponseRESTProvider(transID, restRespCode, startREST)
}

func TestLoggerAPIClientMonitor(t *testing.T) {
	var appName = "crm-util-go"
	transID := common.NewUUID()

	// Ex. Rest client
	logger := InitOutboundLogger(appName, OmxMF)
	logger.Level = LEVEL_ALL

	omxMFUrl := "http://omx-mf.true.th"
	action := "submitOrder"
	omxMFRespCode := "0"

	startOmxMF := logger.LogRequestRESTClient(transID, omxMFUrl, action)
	logger.Info(transID, "TODO: rest client logic")
	logger.LogResponseRESTClient(transID, omxMFUrl, action, omxMFRespCode, startOmxMF)
}

func TestLoggerDBClientMonitor(t *testing.T) {
	var appName = "crm-util-go"
	transID := common.NewUUID()

	// Ex. DAO
	logger := InitInboundLogger(appName, CrmDatabase)
	logger.Level = LEVEL_ALL

	startDAO := logger.LogRequestDBClient(transID)
	logger.Info(transID, "TODO: DAO logic")
	daoRespCode := "0"
	logger.LogResponseDBClient(transID, daoRespCode, startDAO)
}

type TestSecurityAuditBean struct {
	RowID     string `json:"rowID"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func TestLoggerSecurityAuditView(t *testing.T) {
	var appName = "crm-util-go"
	transID := common.NewUUID()

	logger := InitInboundLogger(appName, CrmInbound)
	logger.Level = LEVEL_ALL

	// echo framework
	// Ex. clientIPAddr := c.RealIP()
	clientIPAddr := "192.168.1.1"
	employeeID := "01018298"
	objectName := "BillingAccount"

	request := TestSecurityAuditBean{
		FirstName: "abc",
	}

	response := TestSecurityAuditBean{
		RowID:     "1-1",
		FirstName: "abc",
		LastName:  "xyz",
	}

	isSuccess := true
	remark := "SELECT * FROM TABLE_NAME WHERE FIRST_NAME = ?"

	logger.SecurityAuditView(transID, clientIPAddr, employeeID, objectName,
		request, response, isSuccess, remark)
}

func TestLoggerSecurityAuditModify(t *testing.T) {
	var appName = "crm-util-go"
	transID := common.NewUUID()

	logger := InitInboundLogger(appName, CrmInbound)
	logger.Level = LEVEL_ALL

	// echo framework
	// Ex. clientIPAddr := c.RealIP()
	clientIPAddr := "192.168.1.1"
	employeeID := "01018298"
	objectName := "BillingAccount"

	request := TestSecurityAuditBean{
		RowID:     "1-1",
		FirstName: "abc",
	}

	oldValue := TestSecurityAuditBean{
		RowID:     "1-1",
		FirstName: "aaa",
	}

	response := TestSecurityAuditBean{
		RowID:     "1-1",
		FirstName: "abc",
	}

	isSuccess := true
	remark := "UPDATE TABLE_NAME SET FIRST_NAME = ? WHERE ROW_ID = ?"

	logger.SecurityAuditModify(transID, clientIPAddr, employeeID, objectName,
		request, oldValue, response, isSuccess, remark)
}

func TestLoggerSecurityAuditCreate(t *testing.T) {
	var appName = "crm-util-go"
	transID := common.NewUUID()

	logger := InitInboundLogger(appName, CrmInbound)
	logger.Level = LEVEL_ALL

	// echo framework
	// Ex. clientIPAddr := c.RealIP()
	clientIPAddr := "192.168.1.1"
	employeeID := "01018298"
	objectName := "BillingAccount"

	request := TestSecurityAuditBean{
		RowID:     "1-1",
		FirstName: "aaa",
		LastName:  "bbb",
	}

	response := TestSecurityAuditBean{
		RowID:     "1-1",
		FirstName: "aaa",
		LastName:  "bbb",
	}

	isSuccess := true
	remark := "INSERT INTO TABLE_NAME (ROW_ID, FIRST_NAME, LAST_NAME) VALUES (?,?,?)"

	logger.SecurityAuditCreate(transID, clientIPAddr, employeeID, objectName,
		request, response, isSuccess, remark)
}

func TestLoggerSecurityAuditDelete(t *testing.T) {
	var appName = "crm-util-go"
	transID := common.NewUUID()

	logger := InitInboundLogger(appName, CrmInbound)
	logger.Level = LEVEL_ALL

	// echo framework
	// Ex. clientIPAddr := c.RealIP()
	clientIPAddr := "192.168.1.1"
	employeeID := "01018298"
	objectName := "BillingAccount"

	request := TestSecurityAuditBean{
		RowID: "1-1",
	}

	oldValue := TestSecurityAuditBean{
		RowID:     "1-1",
		FirstName: "aaa",
		LastName:  "bbb",
	}

	response := TestSecurityAuditBean{
		RowID:     "1-1",
		FirstName: "aaa",
		LastName:  "bbb",
	}

	isSuccess := true
	remark := "DELETE FROM TABLE_NAME WHERE ROW_ID = ?"

	logger.SecurityAuditDelete(transID, clientIPAddr, employeeID, objectName,
		request, oldValue, response, isSuccess, remark)
}

func TestLoggerLevel(t *testing.T) {
	var appName = "crm-util-go"
	transID := common.NewUUID()
	/*
		ALL   = 7
		TRACE = 6
		DEBUG = 5
		INFO  = 4
		WARN  = 3
		ERROR = 2
		FATAL = 1
		OFF   = 0
	*/
	logger := InitInboundLogger(appName, CrmInbound)
	logger.Level = LEVEL_ERROR

	logger.Trace(transID, "Trace Level")
	logger.Debug(transID, "Debug Level")
	logger.Info(transID, "Info Level")
	logger.Warn(transID, "Warn Level")
	logger.Error(transID, "Error Level")
	logger.Fatal(transID, "Fatal Level")
}
