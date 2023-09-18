package httpclient

import (
	"FPSelfHegdeAPIGo/logging"
	"net/http"
	"time"
)

type BasicAuthen struct {
	UserName string
	Password string
}

type HttpClient struct {
	MaxConnections  int
	BasicAuthen     BasicAuthen
	Timeout         time.Duration
	CertSkipVerify  bool
	CertServerName  string
	CertPEMFileName string
	Charset         string
	ProxyURL        string
	Logger          *logging.PatternLogger
}

type HttpResponse struct {
	HttpStatusCode int
	HttpStatusMsg  string
	ResponseMsg    string
	HttpHeader     http.Header
	IsRedirect     bool
	RedirectUrl    string
}
