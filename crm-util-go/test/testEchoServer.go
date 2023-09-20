package test

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/http2"
	"net/http"
	"time"
)

func getRequestInfo(c echo.Context) error {
	req := c.Request()
	format := `<code>
		Protocol: %s<br>
		Host: %s<br>
		Remote Address: %s<br>
		Method: %s<br>
		Path: %s<br>
	</code>`

	return c.HTML(http.StatusOK, fmt.Sprintf(format, req.Proto, req.Host, req.RemoteAddr, req.Method, req.URL.Path))
}

func StartEchoHTTP2Server() {
	// Initial Echo Server
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(200)))

	// Router Group
	// http://localhost:8080/GoAPI/getRequestInfo
	r := e.Group("GoAPI")
	r.GET("/getRequestInfo", getRequestInfo)

	h2Server := &http2.Server{
		MaxConcurrentStreams: 200,
		MaxReadFrameSize:     1048576,
		IdleTimeout:          10 * time.Second,
	}

	err := e.StartH2CServer(":8080", h2Server)

	if err != http.ErrServerClosed {
		fmt.Println(err.Error())
	}
}
