package restapi

type MonitorResp struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	TransID string `json:"transID"`
}

type responseSelfHedgeProfile struct {
	Code             string              `json:"code"`
	Msg              string              `json:"msg"`
	BackendUrl       string              `json:"backendUrl,omitempty"`
	TransID          string              `json:"transID"`
	SelfHedgeProfile selfHedgeProfileObj `json:"selfHedgeProfile,omitempty"`
}

type selfHedgeProfileObj struct {
	GroupId string `json:"groupId,omitempty"`
	ResSize int    `json:"-"`
}
