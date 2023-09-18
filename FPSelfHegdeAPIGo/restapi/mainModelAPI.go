package restapi

import (
	"FPSelfHegdeAPIGo/errorcode"
	"FPSelfHegdeAPIGo/httpclient"
	"FPSelfHegdeAPIGo/logging"
	"FPSelfHegdeAPIGo/util"
	//"FPSelfHegdeAPIGo/util"
)

// SelfHegdeController selfHegde component controller
type SelfHegdeController struct {
	Logger     *logging.PatternLogger
	ErrCode    errorcode.FpErrorCode
	HttpClient httpclient.HttpClient
	util       *util.Util
}
