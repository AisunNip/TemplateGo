package main

import (
	"FPSelfHegdeAPIGo/common"
	"FPSelfHegdeAPIGo/restapi"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	config "github.com/spf13/viper"
	"github.com/tylerb/graceful"
)

func main() {

	transID := common.NewUUID()

	// Controller
	ctrl := restapi.NewSelfHegdeController()

	// LoadConfig from yml file
	err := ctrl.LoadConfigFile(transID)
	if err != nil {
		ctrl.Logger.Error(transID, "LoadConfigFile from yml file error: "+err.Error())
		panic("LoadConfigFile from yml file error: " + err.Error())
	}
	fmt.Printf("\ntransID : " + transID)
	fmt.Printf("\nservice.endpoint : " + config.GetString("service.endpoint"))
	fmt.Printf("\nservice.port : " + config.GetString("service.port"))

	// Initial Echo Framework
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(200)))

	// Router Group
	r := e.Group(config.GetString("service.endpoint"))
	r.GET("/monitoring", ctrl.Monitor)
	r.POST("/getToken", ctrl.GetToken)
	//Start http://localhost/FPSelfHegdeAPIGo/monitoring

	// Start Server, Graceful Shutdown with in 5 sec.
	// e.Server.Addr = ":80"
	e.Server.Addr = ":" + config.GetString("service.port")
	err = graceful.ListenAndServe(e.Server, 5*time.Second)

	if err != nil {
		ctrl.Logger.Info(transID, "Start Echo server error because "+err.Error())
	}
}
