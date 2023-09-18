package restapi

import (
	"FPSelfHegdeAPIGo/common"
	"FPSelfHegdeAPIGo/constant"
	"FPSelfHegdeAPIGo/validate"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (ctrl SelfHegdeController) GetToken(c echo.Context) error {
	transID := common.NewUUID()
	responseToken := new(ResponseToken)

	var reqToken ReqToken
	err := c.Bind(reqToken)
	if err != nil {
		ctrl.Logger.Error(transID, "HTTP Json Request to struct")
		fpErrorCodeResp := ctrl.ErrCode.GenerateAppError(constant.AppName, "HTTP Json Request to struct "+err.Error())

		responseToken.Code = fpErrorCodeResp.ErrorCode
		responseToken.Msg = fpErrorCodeResp.ErrorMessage
		return c.JSON(http.StatusOK, responseToken)
	}

	err = validate.ValidateStruct(reqToken)
	if err != nil {
		ctrl.Logger.Error(transID, "HTTP Json Request to struct")
		fpErrorCodeResp := ctrl.ErrCode.GenerateParameterRequire(err.Error())

		responseToken.Code = fpErrorCodeResp.ErrorCode
		responseToken.Msg = fpErrorCodeResp.ErrorMessage

		return c.JSON(http.StatusOK, responseToken)
	}
	getTokenICEResp, MOIerr, url := ctrl.GetTokenICE(transID)
	/*
		Expire -> get again
	*/

	if validate.HasStringValue(MOIerr.ErrorCode) {
		responseToken.Code = MOIerr.ErrorCode
		responseToken.Msg = MOIerr.ErrorMessage
		responseToken.BackendUrl = url
		return c.JSON(http.StatusOK, getTokenICEResp)
	}

	responseToken.Code = constant.ERR_CODE_SUCCESS
	responseToken.Msg = constant.SuccessMsg
	responseToken.Token = getTokenICEResp.Token
	responseToken.ExpireDate = "dd/mm/yyyy"

	return c.JSON(200, responseToken)
}
