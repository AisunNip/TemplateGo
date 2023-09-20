package httpclient

import (
	"bytes"
	"crm-util-go/common"
	"crm-util-go/logging"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultHttpProxyTimeout = 30 * time.Second
)

func LoadConfigHttpProxy(fileName string) (ConfigHttpProxy, error) {
	config := ConfigHttpProxy{}

	file, err := os.ReadFile(fileName)

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	return config, err
}

func getConfigHttpProxy(config ConfigHttpProxy, reqURLPath string) (baseForwardURL string, timeoutSec time.Duration) {
	timeoutSec = defaultHttpProxyTimeout

	for i := 0; i < len(config.ConfigList); i++ {
		for _, path := range config.ConfigList[i].Path {
			if strings.Contains(reqURLPath, path) {
				baseForwardURL = config.ConfigList[i].ForwardURL
				if config.ConfigList[i].TimeoutSec > 0 {
					timeoutSec = time.Duration(config.ConfigList[i].TimeoutSec) * time.Second
				}
				return
			}
		}
	}

	return
}

func StartHttpProxy(addr string) {
	logger := logging.InitOutboundLogger("HttpProxy", logging.AllSystem)
	logger.Level = logging.LEVEL_ALL

	transID := common.NewUUID()
	config, err := LoadConfigHttpProxy("./config/configHttpProxy.json")

	if err != nil {
		logger.Error(transID, "LoadConfigHttpProxy error: "+err.Error())
		panic("LoadConfigHttpProxy error: " + err.Error())
	}

	logger.Info(transID, "LoadConfigHttpProxy Success")

	customTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: transportTimeout,
		}).DialContext,
		DisableCompression: true,
		ForceAttemptHTTP2:  false,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		TLSHandshakeTimeout:   tlsHandshakeTimeout,
		ExpectContinueTimeout: expectContinueTimeout,
		MaxIdleConns:          maxHttpConnections,
		MaxIdleConnsPerHost:   maxHttpConnections,
		MaxConnsPerHost:       maxHttpConnections,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		transID = common.NewUUID()

		baseForwardURL, timeoutSec := getConfigHttpProxy(config, r.URL.Path)

		if baseForwardURL == "" {
			http.Error(w, "Please check configure in HttpProxy", http.StatusInternalServerError)
			return
		}

		// URL
		logger.Info(transID, r.Method, r.URL.Path)

		// Http Header
		var headersBuilder strings.Builder
		headersBuilder.WriteString("Request Http Header Key=Value")
		for key, values := range r.Header {
			for _, value := range values {
				headersBuilder.WriteString("\n" + key + "=" + value)
			}
		}
		logger.Info(transID, headersBuilder.String())

		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Read a request body error "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Body Payload
		logger.Info(transID, "Request Body:", string(requestBody))

		forwardURL := baseForwardURL + r.URL.Path
		logger.Info(transID, "ForwardURL:", forwardURL, "TimeoutSec:", timeoutSec)

		forwardRequest, err := http.NewRequest(r.Method, forwardURL, bytes.NewReader(requestBody))
		if err != nil {
			http.Error(w, "ForwardRequest error "+err.Error(), http.StatusInternalServerError)
			return
		}
		forwardRequest.Header = r.Header

		client := &http.Client{
			Transport: customTransport,
			Timeout:   timeoutSec,
		}

		startDT := time.Now()
		resp, err := client.Do(forwardRequest)
		if err != nil {
			logger.Info(transID, "Send Request Error. ResponseTime:", time.Since(startDT).Milliseconds(), "ms")
			http.Error(w, "Send a request error "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		logger.Info(transID, "Send Request Success. ResponseTime:", time.Since(startDT).Milliseconds(), "ms")

		if resp.StatusCode != http.StatusOK {
			w.WriteHeader(resp.StatusCode)
		}

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		io.Copy(w, resp.Body)
	})

	http.ListenAndServe(addr, nil)
}
