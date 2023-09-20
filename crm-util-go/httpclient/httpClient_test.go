package httpclient_test

import (
	"crm-util-go/common"
	"crm-util-go/errorcode"
	"crm-util-go/httpclient"
	"crm-util-go/logging"
	"errors"
	"strconv"
	"testing"
)

func TestPostLineNotify(t *testing.T) {
	var httpResp httpclient.HttpResponse
	var httpErr error

	var logger = logging.InitInboundLogger("crm-util-go", logging.CrmOutbound)
	logger.Level = logging.LEVEL_ALL

	transID := common.NewUUID()
	lineToken := "MvQBqhT4UMp6GRENMsAxwplNskKYr6fQskouZnB1KGA"
	message := "!! ทดสอบ test 5555"

	httpClient := httpclient.NewHttpClient()
	httpClient.Logger = logger

	targetURL := "https://notify-api.line.me/api/notify"
	action := "PostLineNotify"

	requestDateTime := logger.LogRequestFormClient(transID, targetURL, action)
	httpResp, err := httpClient.PostLineNotify(transID, lineToken, message)
	logger.LogResponseFormClient(transID, targetURL, action,
		strconv.FormatInt(int64(httpResp.HttpStatusCode), 10), requestDateTime)

	if err != nil || httpResp.HttpStatusCode != 200 {
		// Custom Error
		if err != nil {
			httpErr = errorcode.HttpError{
				HttpStatusCode: httpResp.HttpStatusCode,
				Err:            err,
			}
		} else {
			httpErr = errorcode.HttpError{
				HttpStatusCode: httpResp.HttpStatusCode,
				Err:            errors.New(httpResp.HttpStatusMsg),
			}
		}
	}

	if httpErr != nil {
		t.Errorf("TestPostLineNotify Error %s", httpErr.Error())
	} else {
		logger.Info(transID, "HttpStatusCode:", httpResp.HttpStatusCode,
			", HttpStatusMsg:", httpResp.HttpStatusMsg, ", IsRedirect:", httpResp.IsRedirect)
		logger.Info(transID, "ResponseMsg:", httpResp.ResponseMsg)
	}
}

func TestBasicAuthorization(t *testing.T) {
	var logger = logging.InitInboundLogger("crm-util-go", logging.CrmOutbound)
	logger.Level = logging.LEVEL_ALL

	transID := common.NewUUID()

	targetURL := "http://crmapigw-uat4.true.th/CRMIAsset/getLatestAssetRoot"
	jsonReq := `{ "serviceID" : "9600000005" }`
	action := "getLatestAssetRoot"

	httpClient := httpclient.NewHttpClient()
	httpClient.Logger = logger

	// Basic Authorization
	userName := "CRMCID"
	password := "ps&D=4%"
	httpHeaderMap := httpClient.GenerateBasicAuthorization(userName, password)

	requestDateTime := logger.LogRequestRESTClient(transID, targetURL, action)
	httpResp, err := httpClient.PostJson(transID, targetURL, jsonReq, httpHeaderMap)
	logger.LogResponseRESTClient(transID, targetURL, action,
		strconv.FormatInt(int64(httpResp.HttpStatusCode), 10), requestDateTime)

	if err != nil {
		t.Errorf("TestBasicAuthorization Error %s", err.Error())
	} else {
		logger.Info(transID, "HttpStatusCode:", httpResp.HttpStatusCode,
			", HttpStatusMsg:", httpResp.HttpStatusMsg, ", IsRedirect:", httpResp.IsRedirect)
		logger.Info(transID, "ResponseMsg:", httpResp.ResponseMsg)
	}
}

func TestTrustCertFile(t *testing.T) {
	var logger = logging.InitInboundLogger("crm-util-go", logging.CrmOutbound)
	logger.Level = logging.LEVEL_ALL

	transID := common.NewUUID()

	// TDAA : Post "https://10.95.108.180:9200/elf-tx-dev-tdaa-*/_search": x509: certificate signed by unknown authority
	targetURL := "https://10.95.108.180:9200/elf-tx-dev-tdaa-*/_search"

	jsonReq := `{
		"query":{"bool":{"must":[{"match":{"product_id":"TDAA2021011212481325913"}},
		{"match":{"app":"update-loan"}},
		{"match":{"app":"insert-loan"}}]}},
		"sort":[{"@timestamp":"desc"}]
	}`

	httpClient := httpclient.NewHttpClient()
	httpClient.Logger = logger
	/* Trust certificate file */
	httpClient.CertSkipVerify = false
	httpClient.CertServerName = "es.tybdev.tyb.rft"
	httpClient.CertPEMFileName = "../certfile/ca-true-engineer.pem"
	httpClient.BasicAuthen.UserName = "crmapi"
	httpClient.BasicAuthen.Password = "Crm@pi#132"

	action := "searchTDAA"
	requestDateTime := logger.LogRequestRESTClient(transID, targetURL, action)
	httpResp, err := httpClient.PostJson(transID, targetURL, jsonReq, nil)
	logger.LogResponseRESTClient(transID, targetURL, action,
		strconv.FormatInt(int64(httpResp.HttpStatusCode), 10), requestDateTime)

	if err != nil {
		t.Errorf("TestTrustCertFile Error %s", err.Error())
	} else {
		logger.Info(transID, "HttpStatusCode:", httpResp.HttpStatusCode,
			", HttpStatusMsg:", httpResp.HttpStatusMsg, ", IsRedirect:", httpResp.IsRedirect)
		logger.Info(transID, "ResponseMsg:", httpResp.ResponseMsg)
	}
}
