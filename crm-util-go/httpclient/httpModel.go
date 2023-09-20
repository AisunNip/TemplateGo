package httpclient

import (
	"crm-util-go/logging"
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

type ConfigHttpProxy struct {
	ConfigList []ConfigHttpProxyList `json:"configList"`
}

type ConfigHttpProxyList struct {
	Path       []string `json:"path"`
	ForwardURL string   `json:"forwardURL"`
	TimeoutSec int      `json:"timeoutSec"`
}
