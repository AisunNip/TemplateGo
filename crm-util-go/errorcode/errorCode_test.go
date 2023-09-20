package errorcode

import (
	"crm-util-go/httpclient"
	"fmt"
	"testing"
)

func TestErrorCode(t *testing.T) {
	InitConfig("../config")

	var errCode CrmErrorCode
	errCode.SystemCode = "CIB"
	errCode.ModuleCode = "AC"

	var crmErrorCodeResp CrmErrorCodeResp

	// Parameter Require
	crmErrorCodeResp = errCode.GenerateParameterRequire("IDNo")
	fmt.Printf("ErrorCode: %s, ErrorMessage: %s\n", crmErrorCodeResp.ErrorCode, crmErrorCodeResp.ErrorMessage)
	fmt.Println("############################################")

	// Parameter Invalid
	crmErrorCodeResp = errCode.GenerateParameterInvalid("MobileNo")
	fmt.Printf("ErrorCode: %s, ErrorMessage: %s\n", crmErrorCodeResp.ErrorCode, crmErrorCodeResp.ErrorMessage)
	fmt.Println("############################################")

	// Data Not Found
	crmErrorCodeResp = errCode.GenerateDataNotFound("Account", "CRM")
	fmt.Println(errCode.IsDataNotFound(crmErrorCodeResp))
	fmt.Printf("ErrorCode: %s, ErrorMessage: %s\n", crmErrorCodeResp.ErrorCode, crmErrorCodeResp.ErrorMessage)
	fmt.Println("############################################")

	// Business Logic error
	crmErrorCodeResp = errCode.GenerateBusinessLogic("Age >= 20")
	fmt.Printf("ErrorCode: %s, ErrorMessage: %s\n", crmErrorCodeResp.ErrorCode, crmErrorCodeResp.ErrorMessage)
	fmt.Println("############################################")

	// Application Error
	crmErrorCodeResp = errCode.GenerateAppError("CRMIAccountGo", "error message !!")
	fmt.Printf("ErrorCode: %s, ErrorMessage: %s\n", crmErrorCodeResp.ErrorCode, crmErrorCodeResp.ErrorMessage)
	fmt.Println("############################################")

	// Call API error
	var backendResp BackendResp
	backendResp.Url = "http://cfm.true.th/xxx"
	backendResp.MethodName = "createTT"
	backendResp.ErrorCode = "CFM123"
	backendResp.ErrorMessage = "error message from CFM system"

	crmErrorCodeResp = errCode.GenerateByAPI("900003", backendResp)
	fmt.Printf("ErrorCode: %s, ErrorMessage: %s\n", crmErrorCodeResp.ErrorCode, crmErrorCodeResp.ErrorMessage)
	fmt.Println("############################################")

	var httpResp httpclient.HttpResponse
	httpResp.HttpStatusCode = 404
	httpResp.HttpStatusMsg = "Not Found"
	reqURL := "http://intx.true.th/xxx"
	action := "getBilling"
	errorCode := "900019xx"
	crmErrorCodeResp = errCode.GenerateByAPIHttpError(reqURL, action, errorCode, httpResp)
	fmt.Printf("ErrorCode: %s, ErrorMessage: %s\n", crmErrorCodeResp.ErrorCode, crmErrorCodeResp.ErrorMessage)
	fmt.Println("############################################")
}
