package cryptography

import (
	"FPSelfHegdeAPIGo/common"
	"FPSelfHegdeAPIGo/constant"
	"FPSelfHegdeAPIGo/logging"
	"fmt"
	"strings"

	config "github.com/spf13/viper"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var logger = logging.InitScheduleLogger(constant.AppName, logging.FpUtil)

func CheckHeaderAuthorization(tranId string, c echo.Context) map[string]string {
	logger.Info(tranId, " Start process function CheckHeaderAuthorization.")

	rsMap := make(map[string]string)
	var code = ""
	var message = ""
	var auth = c.Request().Header.Get("Authorization")
	if len(auth) == 0 {
		logger.Info(tranId, " Headers Authorization is require.")
		code = constant.RequiredField
		message = "Headers authorization is require."

	} else {
		if !strings.HasPrefix(auth, "Bearer") {
			logger.Info(tranId, "Authorization : "+auth)
			code = constant.AuthorizeInCorrect
			message = "Headers Authorization incorrect."
		} else {
			logger.Info(tranId, "Authorization : "+auth)
			code = constant.SuccessCode
			rsMap["authorization"] = auth
		}
	}
	rsMap["code"] = code
	rsMap["message"] = message
	return rsMap
}

func MapErrorCode(code string, message string, tranID string) (response common.ResponseBean) {
	response.Code = code
	response.Msg = message
	response.TransID = tranID
	return
}

func VerifyJWToken(tokenString string) (res JwtClaims, err error) {
	secret := config.GetString("crm.secret")

	if strings.HasPrefix(tokenString, "Bearer") {
		tokenString = strings.Split(tokenString, "Bearer ")[1]
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err == nil {
		res = MapToJwtClaims(claims)
	}

	return
}

func MapToJwtClaims(mapClaims jwt.MapClaims) (res JwtClaims) {
	//var jwtClaims JwtClaims
	for key, val := range mapClaims {
		switch key {
		case "empId":
			res.EmpId = fmt.Sprintf("%v", val)
		case "firstName":
			res.FirstName = fmt.Sprintf("%v", val)
		case "lastName":
			res.LastName = fmt.Sprintf("%v", val)
		case "division":
			res.Division = fmt.Sprintf("%v", val)
		case "divisionRowId":
			res.DivisionRowId = fmt.Sprintf("%v", val)
		case "email":
			res.Email = fmt.Sprintf("%v", val)
		case "position":
			res.Position = fmt.Sprintf("%v", val)
		case "dealerCode":
			res.DealerCode = fmt.Sprintf("%v", val)
		case "saleCode":
			res.SaleCode = fmt.Sprintf("%v", val)
		case "ChanelAlias":
			res.ChanelAlias = fmt.Sprintf("%v", val)
		case "roles":
			res.Roles = fmt.Sprintf("%v", val)
		case "divisionType":
			res.DivisionType = fmt.Sprintf("%v", val)
		case "user":
			res.User = fmt.Sprintf("%v", val)
		case "client":
			res.Client = fmt.Sprintf("%v", val)
		case "ShopCode":
			res.ShopCode = fmt.Sprintf("%v", val)
		case "host":
			res.Host = fmt.Sprintf("%v", val)
		case "type":
			res.Type = fmt.Sprintf("%v", val)
		default:
			fmt.Sprintf("%v", val)
		}
	}
	return
}
