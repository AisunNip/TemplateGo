package db

import (
	"context"
	"crm-util-go/common"
	"crm-util-go/errorcode"
	"crm-util-go/logging"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"
)

type CampaignTransReqBean struct {
	Status      string      `json:"status,omitempty"`
	SubStatus   string      `json:"subStatus,omitempty"`
	CreatedDate CrmDateTime `json:"createdDate,omitempty"`
}

type CampaignTransBean struct {
	CampTransID *string      `json:"campTransID,omitempty"`
	CampID      *string      `json:"campID,omitempty"`
	Status      *string      `json:"status,omitempty"`
	SubStatus   *string      `json:"subStatus,omitempty"`
	CallStatus  *string      `json:"callStatus,omitempty"`
	PinReqDate  *CrmDateTime `json:"pinReqDate,omitempty"`
	CreatedBy   *string      `json:"createdBy,omitempty"`
	CreatedDate *CrmDateTime `json:"createdDate,omitempty"`
}

type CampaignTransRespBean struct {
	Code                  string              `json:"code,omitempty"`
	Msg                   string              `json:"msg,omitempty"`
	TransID               string              `json:"transID,omitempty"`
	CampaignTransBeanList []CampaignTransBean `json:"campaignTransBeanList,omitempty"`
}

var daoLogger *logging.PatternLogger
var daoErrCode errorcode.CrmErrorCode
var dbPool DBPool

func initMaria() {
	daoLogger = logging.InitInboundLogger("crm-util-go", logging.CrmDatabase)
	daoLogger.Level = logging.LEVEL_ALL

	// UAT: 172.19.208.111
	// PRD report: 172.19.249.173

	// Initial database pool
	// MySQL Note: struct time.Time --> &parseTime=true&loc=Asia%2FBangkok
	// dbPool.DataSourceName = "crmadm:crmadm_001@tcp(172.19.249.173:3306)/CRMX?charset=utf8&checkConnLiveness=true&timeout=5s&readTimeout=60s&writeTimeout=60s&parseTime=true&loc=Asia%2FBangkok"

	dbPool.DataSourceName = "crmapp:crmapp2020@tcp(172.19.208.111:3306)/CRMX2?charset=utf8&checkConnLiveness=true&timeout=5s&readTimeout=30s&writeTimeout=30s&parseTime=true&loc=Asia%2FBangkok"
	dbPool.MaxOpenConns = 50
	dbPool.MaxIdleConns = 5
	dbPool.MaxLifetime = 3 * time.Minute
	dbPool.Logger = daoLogger

	// Init Error Code
	errorcode.InitConfig("../config")
	daoErrCode.SystemCode = "CIB"
	daoErrCode.ModuleCode = "CN"
}

func TestQueryMaria(t *testing.T) {
	transID := common.NewUUID()

	initMaria()
	defer dbPool.CloseMariaDBPool()

	campReqBean := new(CampaignTransReqBean)
	campReqBean.Status = "Accept"

	campTransRespBean := getCampTransListDAO(transID, campReqBean)

	fmt.Println(fmt.Sprintf("TransId: %s, CampTransRespBean: %#v", transID, campTransRespBean))

	if campTransRespBean.Code != "0" {
		t.Errorf("expecting result to be Code=0 but %s", campTransRespBean.Code)
	}
}

func generateSQLGetCampTrans(request *CampaignTransReqBean) (sql string, bindValue []interface{}) {
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("SELECT CAMP_TRANS_ID, CAMP_ID, STATUS, SUB_STATUS, CALL_STATUS, PIN_REQUEST_DATE, CREATED_BY, CREATED_DATE FROM CAMPAIGN_TRANS WHERE 1=1")

	if len(request.Status) > 0 {
		sqlBuilder.WriteString(" AND STATUS = ?")
		bindValue = append(bindValue, request.Status)
	}

	if len(request.SubStatus) > 0 {
		sqlBuilder.WriteString(" AND SUB_STATUS = ?")
		bindValue = append(bindValue, request.SubStatus)
	}

	if !request.CreatedDate.Time().IsZero() {
		sqlBuilder.WriteString(" AND CREATED_DATE >= ?")
		bindValue = append(bindValue, request.CreatedDate.Time())
	}

	sqlBuilder.WriteString(" LIMIT 50")

	sql = sqlBuilder.String()

	return
}

func getCampTransListDAO(transID string, request *CampaignTransReqBean) CampaignTransRespBean {
	var campTransRespBean CampaignTransRespBean
	var campTransBeanList []CampaignTransBean

	startDT := daoLogger.LogRequestDBClient(transID)

	campTransRespBean.TransID = transID

	sqlStmt, bindValue := generateSQLGetCampTrans(request)
	daoLogger.Debug(transID, sqlStmt)

	prepare, crmErrorCodeResp := dbPool.CreatePreparedStatementMaria(transID, sqlStmt, daoErrCode)

	if crmErrorCodeResp != nil {
		daoLogger.LogResponseDBClient(transID, crmErrorCodeResp.ErrorCode, startDT)

		campTransRespBean.Code = crmErrorCodeResp.ErrorCode
		campTransRespBean.Msg = crmErrorCodeResp.ErrorMessage

		return campTransRespBean
	}
	defer prepare.Close()

	var rows *sql.Rows
	// rows, err = prepare.Query(bindValue...)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	rows, err := prepare.QueryContext(ctx, bindValue...)

	if err != nil {
		daoLogger.Error(transID, "Query DB Error: "+err.Error())
		crmErrorCodeResp := daoErrCode.GenerateAppError("crm-util-go", "Query DB Error: "+err.Error())

		campTransRespBean.Code = crmErrorCodeResp.ErrorCode
		campTransRespBean.Msg = crmErrorCodeResp.ErrorMessage

		return campTransRespBean
	}

	defer rows.Close()

	for rows.Next() {
		var campaignTransBean CampaignTransBean

		err = rows.Scan(&campaignTransBean.CampTransID,
			&campaignTransBean.CampID,
			&campaignTransBean.Status,
			&campaignTransBean.SubStatus,
			&campaignTransBean.CallStatus,
			&campaignTransBean.PinReqDate,
			&campaignTransBean.CreatedBy,
			&campaignTransBean.CreatedDate)

		if err != nil {
			crmErrorCodeResp := daoErrCode.GenerateAppError("crm-util-go",
				"Can not read row from database. "+err.Error())

			campTransRespBean.Code = crmErrorCodeResp.ErrorCode
			campTransRespBean.Msg = crmErrorCodeResp.ErrorMessage

			return campTransRespBean
		}

		campTransBeanList = append(campTransBeanList, campaignTransBean)
	}

	if len(campTransBeanList) > 0 {
		campTransRespBean.Code = "0"
		campTransRespBean.Msg = "Success"
		campTransRespBean.CampaignTransBeanList = campTransBeanList
	} else {
		crmErrorCodeResp := daoErrCode.GenerateDataNotFound("CampaignTrans", "Campaign")

		campTransRespBean.Code = crmErrorCodeResp.ErrorCode
		campTransRespBean.Msg = crmErrorCodeResp.ErrorMessage
	}

	daoLogger.LogResponseDBClient(transID, campTransRespBean.Code, startDT)
	return campTransRespBean
}
