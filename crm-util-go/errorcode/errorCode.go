package errorcode

import (
	"crm-util-go/common"
	"crm-util-go/httpclient"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type CrmErrorCode struct {
	SystemCode string
	ModuleCode string
}

type CrmErrorCodeResp struct {
	ErrorCode    string
	ErrorMessage string
}

type BackendResp struct {
	Url          string
	MethodName   string
	ErrorCode    string
	ErrorMessage string
}

func InitConfig(configPath string) {
	viper.SetConfigName("errorCodeConfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Sprintf("ErrorCode read errorCodeConfig.yml file occur error: %s", err.Error()))
	}
}

func (e CrmErrorCode) getErrorMessage(errCode string, values []interface{}) (string, error) {
	msg := viper.GetString("ErrorCode." + errCode)

	if msg == "" {
		return "", errors.New("Error Code=" + errCode + " not found")
	}

	return fmt.Sprintf(msg, values...), nil
}

func (e CrmErrorCode) generateErrorMsg(errorCode string, errorCodeDesc string, errorMsg string) string {
	var msgBuilder strings.Builder
	msgBuilder.WriteString("[")
	msgBuilder.WriteString(errorCode)
	msgBuilder.WriteString("] ")
	msgBuilder.WriteString("[")
	msgBuilder.WriteString(errorCodeDesc)
	msgBuilder.WriteString("] ")
	msgBuilder.WriteString("[")
	msgBuilder.WriteString(errorMsg)
	msgBuilder.WriteString("]")

	return msgBuilder.String()
}

func (e CrmErrorCode) Generate(errCode string, values []interface{}) CrmErrorCodeResp {
	msg, err := e.getErrorMessage(errCode, values)

	if err != nil {
		return e.GenerateDataNotFound("ErrorCode="+errCode, "CRM")
	}

	msgArray := strings.Split(msg, "|")

	var codeBuilder strings.Builder
	codeBuilder.WriteString(e.SystemCode)
	codeBuilder.WriteString(e.ModuleCode)
	codeBuilder.WriteString(msgArray[2])
	codeBuilder.WriteString(errCode)
	codeBuilder.WriteString(msgArray[3])

	var errorCodeResp CrmErrorCodeResp
	errorCodeResp.ErrorCode = codeBuilder.String()
	errorCodeResp.ErrorMessage = e.generateErrorMsg(errorCodeResp.ErrorCode, msgArray[1], msgArray[0])

	return errorCodeResp
	/*
		errorMsgPattern := msgArray[0];
		errorCodeDesc := msgArray[1];
		errorClassify := msgArray[2];
		severityType := msgArray[3];
	*/
}

func (e CrmErrorCode) GenerateOneVal(errCode string, value string) CrmErrorCodeResp {
	values := []interface{}{value}
	return e.Generate(errCode, values)
}

func (e CrmErrorCode) GenerateByAPI(errCode string, backendResp BackendResp) CrmErrorCodeResp {
	values := []interface{}{backendResp.Url, backendResp.MethodName}
	crmErrorCodeResp := e.Generate(errCode, values)

	var msgBuilder strings.Builder
	msgBuilder.WriteString(e.generateErrorMsg(backendResp.ErrorCode, "", backendResp.ErrorMessage))
	msgBuilder.WriteString(" ")
	msgBuilder.WriteString(crmErrorCodeResp.ErrorMessage)

	crmErrorCodeResp.ErrorCode = backendResp.ErrorCode
	crmErrorCodeResp.ErrorMessage = msgBuilder.String()

	return crmErrorCodeResp
}

func (e CrmErrorCode) GenerateByAPIRESTConnFail(url string, methodName string, errorMsg string) CrmErrorCodeResp {
	values := []interface{}{url, methodName, errorMsg}
	return e.Generate("920003", values)
}

func (e CrmErrorCode) GenerateByAPIWSConnFail(url string, methodName string, errorMsg string) CrmErrorCodeResp {
	values := []interface{}{url, methodName, errorMsg}
	return e.Generate("920001", values)
}

func (e CrmErrorCode) GenerateByAPIXMLConnFail(url string, methodName string, errorMsg string) CrmErrorCodeResp {
	values := []interface{}{url, methodName, errorMsg}
	return e.Generate("920000", values)
}

func (e CrmErrorCode) GenerateByAPIHttpError(url string, methodName string, errorCode string,
	httpResp httpclient.HttpResponse) CrmErrorCodeResp {

	var backendResp BackendResp
	backendResp.Url = url
	backendResp.MethodName = methodName
	backendResp.ErrorCode = "HttpStatusCode: " + common.IntToString(httpResp.HttpStatusCode)
	backendResp.ErrorMessage = "HttpStatusMsg: " + httpResp.HttpStatusMsg

	return e.GenerateByAPI(errorCode, backendResp)
}

func (e CrmErrorCode) GenerateAppError(appName string, errMsg string) CrmErrorCodeResp {
	values := []interface{}{appName, errMsg}
	return e.Generate("400000", values)
}

func (e CrmErrorCode) GenerateParameterRequire(fieldName string) CrmErrorCodeResp {
	return e.GenerateOneVal("101000", fieldName)
}

func (e CrmErrorCode) GenerateParameterInvalid(fieldName string) CrmErrorCodeResp {
	return e.GenerateOneVal("102000", fieldName)
}

func (e CrmErrorCode) GenerateDataNotFound(data string, systemName string) CrmErrorCodeResp {
	values := []interface{}{data, systemName}
	return e.Generate("201000", values)
}

func (e CrmErrorCode) GenerateBusinessLogic(reason string) CrmErrorCodeResp {
	return e.GenerateOneVal("300000", reason)
}

func (e CrmErrorCode) GenerateCRMDatabaseError(reason string) CrmErrorCodeResp {
	return e.GenerateOneVal("802014", reason)
}

func (e CrmErrorCode) IsDataNotFound(crmErrorCodeResp CrmErrorCodeResp) bool {
	isDataNotFound := false

	if strings.Contains(crmErrorCodeResp.ErrorCode, "201000") {
		isDataNotFound = true
	}

	return isDataNotFound
}
