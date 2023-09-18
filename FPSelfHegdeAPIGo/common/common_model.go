package common

type ResponseBean struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	TransID string `json:"transID"`
}
