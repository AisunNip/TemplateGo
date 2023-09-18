package errorcode

import (
	"FPSelfHegdeAPIGo/common"
	"FPSelfHegdeAPIGo/httpclient"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type FpErrorCode struct {
	SystemCode string
	ModuleCode string
}

type FpErrorCodeResp struct {
	ErrorCode    string
	ErrorMessage string
}

type BackendResp struct {
	Url          string
	MethodName   string
	ErrorCode    string
	ErrorMessage string
}

func ListErrorBackendResp(url, MethodName, ErrorCode, ErrorMessage string) (BackendRespList BackendResp) {
	BackendRespList = BackendResp{
		Url:          url,
		MethodName:   MethodName,
		ErrorCode:    ErrorCode,
		ErrorMessage: ErrorMessage,
	}
	return
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

func (e FpErrorCode) getErrorMessage(errCode string, values []interface{}) (string, error) {
	msg := viper.GetString("ErrorCode." + errCode)

	if msg == "" {
		return "", errors.New("Error Code=" + errCode + " not found")
	}

	return fmt.Sprintf(msg, values...), nil
}

func (e FpErrorCode) generateErrorMsg(errorCode string, errorCodeDesc string, errorMsg string) string {
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

func (e FpErrorCode) Generate(errCode string, values []interface{}) FpErrorCodeResp {
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

	var errorCodeResp FpErrorCodeResp
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

func (e FpErrorCode) GenerateOneVal(errCode string, value string) FpErrorCodeResp {
	values := []interface{}{value}
	return e.Generate(errCode, values)
}

func (e FpErrorCode) GenerateByAPI(errCode string, backendResp BackendResp) FpErrorCodeResp {
	values := []interface{}{backendResp.Url, backendResp.MethodName}
	fpErrorCodeResp := e.Generate(errCode, values)

	var msgBuilder strings.Builder
	msgBuilder.WriteString(e.generateErrorMsg(backendResp.ErrorCode, "", backendResp.ErrorMessage))
	msgBuilder.WriteString(" ")
	msgBuilder.WriteString(fpErrorCodeResp.ErrorMessage)

	fpErrorCodeResp.ErrorCode = backendResp.ErrorCode
	fpErrorCodeResp.ErrorMessage = msgBuilder.String()

	return fpErrorCodeResp
}

func (e FpErrorCode) GenerateByAPIRESTConnFail(url string, methodName string, errorMsg string) FpErrorCodeResp {
	values := []interface{}{url, methodName, errorMsg}
	return e.Generate("920003", values)
}

func (e FpErrorCode) GenerateByAPIWSConnFail(url string, methodName string, errorMsg string) FpErrorCodeResp {
	values := []interface{}{url, methodName, errorMsg}
	return e.Generate("920001", values)
}

func (e FpErrorCode) GenerateByAPIXMLConnFail(url string, methodName string, errorMsg string) FpErrorCodeResp {
	values := []interface{}{url, methodName, errorMsg}
	return e.Generate("920000", values)
}

func (e FpErrorCode) GenerateByAPIHttpError(url string, methodName string, errorCode string,
	httpResp httpclient.HttpResponse) FpErrorCodeResp {

	var backendResp BackendResp
	backendResp.Url = url
	backendResp.MethodName = methodName
	backendResp.ErrorCode = "HttpStatusCode: " + common.IntToString(httpResp.HttpStatusCode)
	backendResp.ErrorMessage = "HttpStatusMsg: " + httpResp.HttpStatusMsg

	return e.GenerateByAPI(errorCode, backendResp)
}

func (e FpErrorCode) GenerateAppError(appName string, errMsg string) FpErrorCodeResp {
	values := []interface{}{appName, errMsg}
	return e.Generate("400000", values)
}

func (e FpErrorCode) GenerateParameterRequire(fieldName string) FpErrorCodeResp {
	return e.GenerateOneVal("101000", fieldName)
}

func (e FpErrorCode) GenerateParameterInvalid(fieldName string) FpErrorCodeResp {
	return e.GenerateOneVal("102000", fieldName)
}

func (e FpErrorCode) GenerateDataNotFound(data string, systemName string) FpErrorCodeResp {
	values := []interface{}{data, systemName}
	return e.Generate("201000", values)
}

func (e FpErrorCode) GenerateBusinessLogic(reason string) FpErrorCodeResp {
	return e.GenerateOneVal("300000", reason)
}

func (e FpErrorCode) GenerateFPDatabaseError(reason string) FpErrorCodeResp {
	return e.GenerateOneVal("802014", reason)
}

func (e FpErrorCode) IsDataNotFound(fpErrorCodeResp FpErrorCodeResp) bool {
	isDataNotFound := false

	if strings.Contains(fpErrorCodeResp.ErrorCode, "201000") {
		isDataNotFound = true
	}

	return isDataNotFound
}
