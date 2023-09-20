package httpclient

import (
	"bytes"
	"crm-util-go/validate"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	maxHttpConnections    = 100
	defaultTimeout        = 15 * time.Second
	transportTimeout      = 3 * time.Second
	tlsHandshakeTimeout   = 3 * time.Second
	expectContinueTimeout = 1 * time.Second
)

func NewHttpClient() HttpClient {
	return HttpClient{
		MaxConnections: maxHttpConnections,
		CertSkipVerify: true,
		Charset:        "utf-8",
		Timeout:        defaultTimeout,
	}
}

func (hc HttpClient) GenerateBasicAuthorization(userName string, password string) map[string]string {
	auth := userName + ":" + password
	httpHeaderMap := make(map[string]string)
	httpHeaderMap["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	return httpHeaderMap
}

func (hc HttpClient) GenerateBearerAuthorization(token string) map[string]string {
	httpHeaderMap := make(map[string]string)
	httpHeaderMap["Authorization"] = "Bearer " + token
	return httpHeaderMap
}

func (hc HttpClient) logResponseTime(transID string, startDT time.Time, reqURL string) {
	hc.Logger.Info(transID, "Response Time:", time.Since(startDT).Milliseconds(), "ms.,",
		"Request URL:", reqURL)
}

func (hc HttpClient) getServerName(reqURL string) (string, error) {
	urlRequest, err := url.Parse(reqURL)

	if err != nil {
		return "", err
	}

	return urlRequest.Hostname(), nil
}

func (hc HttpClient) send(transID string, method string, reqURL string,
	body string, httpHeaderMap map[string]string) (httpResp HttpResponse, err error) {

	startDT := time.Now()

	hc.Logger.WriteRequestMsg(transID, reqURL, method, body)
	defer hc.logResponseTime(transID, startDT, reqURL)

	reqURL = strings.TrimSpace(reqURL)

	if hc.Timeout.Seconds() == 0 {
		hc.Timeout = defaultTimeout
	}

	// Proxy
	var proxyFunc func(*http.Request) (*url.URL, error)

	if validate.HasStringValue(hc.ProxyURL) {
		proxyURL, err := url.Parse(hc.ProxyURL)

		if err != nil {
			hc.Logger.Error(transID, "Parse ProxyURL Error", err)
			return httpResp, err
		}

		proxyFunc = http.ProxyURL(proxyURL)
	}

	var certPool *x509.CertPool

	if validate.HasStringValue(hc.CertPEMFileName) {
		certPool = x509.NewCertPool()
		pemData, err := os.ReadFile(hc.CertPEMFileName)

		if err != nil {
			hc.Logger.Error(transID, "Error read a certificate file", err)
			return httpResp, err
		}

		certPool.AppendCertsFromPEM(pemData)
	}

	customTransport := &http.Transport{
		Proxy: proxyFunc,
		DialContext: (&net.Dialer{
			Timeout: transportTimeout,
		}).DialContext,
		DisableCompression: true,
		ForceAttemptHTTP2:  false,
		TLSClientConfig: &tls.Config{
			RootCAs:            certPool,
			ServerName:         hc.CertServerName,
			InsecureSkipVerify: hc.CertSkipVerify,
		},
		TLSHandshakeTimeout:   tlsHandshakeTimeout,
		ExpectContinueTimeout: expectContinueTimeout,
		MaxIdleConns:          hc.MaxConnections,
		MaxIdleConnsPerHost:   hc.MaxConnections,
		MaxConnsPerHost:       hc.MaxConnections,
	}

	client := &http.Client{
		Transport: customTransport,
		Timeout:   hc.Timeout,
	}

	var httpReq *http.Request
	httpReq, err = http.NewRequest(method, reqURL, bytes.NewBufferString(body))

	if err != nil {
		hc.Logger.Error(transID, "Can not create new request", err)
		return httpResp, err
	}

	if httpHeaderMap != nil {
		for k, v := range httpHeaderMap {
			httpReq.Header.Set(k, v)
		}
	}

	httpReq.Header.Set("Cache-Control", "no-cache")

	if validate.HasStringValue(hc.BasicAuthen.UserName) && validate.HasStringValue(hc.BasicAuthen.Password) {
		httpReq.SetBasicAuth(hc.BasicAuthen.UserName, hc.BasicAuthen.Password)
	}

	hc.Logger.Info(transID, "Send a request to http server. Request URL: "+reqURL)
	var resp *http.Response
	resp, err = client.Do(httpReq)

	if err != nil {
		hc.Logger.Error(transID, "Error send a http request. Request URL: "+reqURL+", Error: "+err.Error())
		return httpResp, err
	}

	defer resp.Body.Close()

	var rawBody []byte
	rawBody, err = io.ReadAll(resp.Body)

	if err != nil {
		hc.Logger.Error(transID, "Can not read response message", err)
		return httpResp, err
	}

	httpResp.HttpStatusCode = resp.StatusCode
	httpResp.HttpStatusMsg = resp.Status
	httpResp.ResponseMsg = string(rawBody)
	httpResp.HttpHeader = resp.Header

	if httpResp.HttpStatusCode == 301 || httpResp.HttpStatusCode == 302 || httpResp.HttpStatusCode == 303 ||
		httpResp.HttpStatusCode == 307 || httpResp.HttpStatusCode == 308 {
		httpResp.IsRedirect = true
		httpResp.RedirectUrl = resp.Header.Get("Location")
	}

	hc.Logger.WriteResponseMsg(transID, httpResp)

	return httpResp, err
}

func (hc HttpClient) Put(transID string, url string, body string, httpHeaderMap map[string]string) (httpResp HttpResponse, err error) {
	return hc.send(transID, "PUT", url, body, httpHeaderMap)
}

func (hc HttpClient) Patch(transID string, url string, body string, httpHeaderMap map[string]string) (httpResp HttpResponse, err error) {
	return hc.send(transID, "PATCH", url, body, httpHeaderMap)
}

func (hc HttpClient) Delete(transID string, url string, body string, httpHeaderMap map[string]string) (httpResp HttpResponse, err error) {
	return hc.send(transID, "DELETE", url, body, httpHeaderMap)
}

func (hc HttpClient) Get(transID string, url string, httpHeaderMap map[string]string) (httpResp HttpResponse, err error) {
	return hc.send(transID, "GET", url, "", httpHeaderMap)
}

func (hc HttpClient) Post(transID string, url string, body string, httpHeaderMap map[string]string) (httpResp HttpResponse, err error) {
	return hc.send(transID, "POST", url, body, httpHeaderMap)
}

func (hc HttpClient) PostJson(transID string, url string, body string, httpHeaderMap map[string]string) (httpResp HttpResponse, err error) {
	if httpHeaderMap == nil {
		httpHeaderMap = make(map[string]string)
	}

	if validate.HasStringValue(hc.Charset) {
		httpHeaderMap["Content-Type"] = "application/json; charset=" + strings.ToLower(hc.Charset)
	} else {
		httpHeaderMap["Content-Type"] = "application/json"
	}

	return hc.send(transID, "POST", url, body, httpHeaderMap)
}

func (hc HttpClient) PostXML(transID string, url string, body string, httpHeaderMap map[string]string) (httpResp HttpResponse, err error) {
	if httpHeaderMap == nil {
		httpHeaderMap = make(map[string]string)
	}

	if validate.HasStringValue(hc.Charset) {
		httpHeaderMap["Content-Type"] = "text/xml; charset=" + strings.ToLower(hc.Charset)
	} else {
		httpHeaderMap["Content-Type"] = "text/xml"
	}

	return hc.send(transID, "POST", url, body, httpHeaderMap)
}

func (hc HttpClient) PostLineNotify(transID string, lineToken string, message string) (httpResp HttpResponse, err error) {

	httpHeaderMap := make(map[string]string)

	if validate.HasStringValue(hc.Charset) {
		httpHeaderMap["Content-Type"] = "application/x-www-form-urlencoded; charset=" + strings.ToLower(hc.Charset)
	} else {
		httpHeaderMap["Content-Type"] = "application/x-www-form-urlencoded"
	}

	httpHeaderMap["Authorization"] = "Bearer " + lineToken

	lineURL := "https://notify-api.line.me/api/notify"

	body := make(map[string]string)
	body["message"] = message
	encodedBody := hc.EncodeFormBody(body)

	return hc.send(transID, "POST", lineURL, encodedBody, httpHeaderMap)
}

func (hc HttpClient) PostForm(transID string, url string, body map[string]string, httpHeaderMap map[string]string) (httpResp HttpResponse, err error) {
	if httpHeaderMap == nil {
		httpHeaderMap = make(map[string]string)
	}

	if validate.HasStringValue(hc.Charset) {
		httpHeaderMap["Content-Type"] = "application/x-www-form-urlencoded; charset=" + strings.ToLower(hc.Charset)
	} else {
		httpHeaderMap["Content-Type"] = "application/x-www-form-urlencoded"
	}

	encodedBody := hc.EncodeFormBody(body)

	return hc.send(transID, "POST", url, encodedBody, httpHeaderMap)
}

func (hc HttpClient) EncodeFormBody(body map[string]string) string {
	var encodedData string

	if body != nil {
		formData := url.Values{}

		for k, v := range body {
			formData.Add(k, v)
		}

		encodedData = formData.Encode()
	}

	return encodedData
}
