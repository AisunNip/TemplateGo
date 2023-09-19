package restapi

import (
	"FPSelfHegdeAPIGo/common"
	"FPSelfHegdeAPIGo/constant"
	"FPSelfHegdeAPIGo/errorcode"
	"FPSelfHegdeAPIGo/httpclient"
	"FPSelfHegdeAPIGo/logging"
	"FPSelfHegdeAPIGo/util"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	config "github.com/spf13/viper"
)

// NewAssetController new asset controller
func NewSelfHegdeController() SelfHegdeController {
	logger := logging.NewLogger()
	logger.Level = logging.LEVEL_ALL

	selfHegdeController := SelfHegdeController{}
	selfHegdeController.Logger = logger
	selfHegdeController.HttpClient = NewSelfHegdeHttpClient(logger)
	selfHegdeController.ErrCode = errorcode.FpErrorCode{
		SystemCode: "COB",
		ModuleCode: "SH",
	}
	selfHegdeController.util = NewUtil(logger)
	return selfHegdeController
}

func NewUtil(logger *logging.PatternLogger) *util.Util {
	util := new(util.Util)
	util.Logger = logger
	return util
}

func NewSelfHegdeHttpClient(logger *logging.PatternLogger) httpclient.HttpClient {
	httpClient := httpclient.HttpClient{}
	httpClient.Charset = "utf-8"
	httpClient.CertSkipVerify = true
	httpClient.Timeout = 10 * time.Second
	httpClient.Logger = logger
	return httpClient
}

func (ctrl SelfHegdeController) LoadConfigFile(transID string) error {
	env := os.Getenv("ENV")
	if env == "" {
		env = os.Args[1]
	}

	ctrl.Logger.Info(transID, fmt.Sprintf("Server start running on %s environment configuration", env))
	config.SetConfigName(env)
	config.SetConfigType("yaml")
	config.AddConfigPath("./config")
	err := config.ReadInConfig()
	if err != nil {
		errMsg := fmt.Sprintf("Read config file %s.yml occur error: %s", env, err.Error())
		ctrl.Logger.Error(transID, errMsg)
		return err
	}

	config.SetConfigName("errorCodeConfig")
	config.SetConfigType("yaml")
	config.AddConfigPath("./config")
	err = config.MergeInConfig()
	if err != nil {
		errMsg := "Read config file errorCodeConfig.yml occur error: " + err.Error()
		ctrl.Logger.Error(transID, errMsg, err)
	}
	return err
}

func (ctrl SelfHegdeController) Monitor(c echo.Context) error {
	transID := common.NewUUID()
	startDAO := ctrl.Logger.LogRequestDBClient(transID)

	var monitorResp MonitorResp
	monitorResp.TransID = transID

	monitorResp.Code = constant.SuccessCode
	monitorResp.Msg = constant.SuccessMsg

	ctrl.Logger.WriteResponseMsg(transID, monitorResp)
	ctrl.Logger.LogResponseDBClient(transID, monitorResp.Code, startDAO)
	return c.JSON(http.StatusOK, monitorResp)
}
