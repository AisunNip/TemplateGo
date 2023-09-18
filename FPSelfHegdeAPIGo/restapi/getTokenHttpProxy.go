package restapi

import (
	"FPSelfHegdeAPIGo/constant"
	"FPSelfHegdeAPIGo/errorcode"
	"FPSelfHegdeAPIGo/fpencoding"
	"encoding/json"

	config "github.com/spf13/viper"
)

func (ctrl SelfHegdeController) GetTokenICE(transID string) (*ResponseTokenICE, errorcode.FpErrorCodeResp, string) {
	var fpErrorCodeResp errorcode.FpErrorCodeResp
	action := "GetTokenICE"
	reqTokenICE := ReqTokenICE{
		Username: config.GetString("ice.authenticate.username"),
		Password: config.GetString("ice.authenticate.password"),
	}

	url := config.GetString("ice.authenticate.url")
	httpHeaderMap := map[string]string{"Content-Type": "application/json"}

	body, _ := json.Marshal(reqTokenICE)
	httpResp, err := ctrl.HttpClient.Post(transID, url, string(body), httpHeaderMap)

	if err != nil {
		fpErrorCodeResp = ctrl.ErrCode.GenerateByAPIRESTConnFail(url, action, err.Error())
		return nil, fpErrorCodeResp, url
	}

	if httpResp.HttpStatusCode != 200 {
		fpErrorCodeResp = ctrl.ErrCode.GenerateByAPIHttpError(url, action, constant.ERR_CODE_ICE_API, httpResp)
		return nil, fpErrorCodeResp, url
	}

	responseTokenICE := new(ResponseTokenICE)
	err = fpencoding.JsonToStruct(httpResp.ResponseMsg, responseTokenICE)
	if err != nil {
		fpErrorCodeResp = ctrl.ErrCode.GenerateAppError(constant.AppName,
			"URL: "+url+" json response message invalid "+err.Error())
		return nil, fpErrorCodeResp, url
	}

	if responseTokenICE.Status != constant.SuccessMsg {
		fpErrorCodeResp.ErrorCode = responseTokenICE.Errorcode
		fpErrorCodeResp.ErrorMessage = responseTokenICE.ErrorDescription

		return nil, fpErrorCodeResp, url
	}
	return responseTokenICE, fpErrorCodeResp, url
}
