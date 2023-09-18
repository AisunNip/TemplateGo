package restapi

type ReqTokenICE struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type ResponseTokenICE struct {
	ReqeustId        string `json:"reqeustId"`
	Status           string `json:"status"`
	Token            string `json:"token"`
	Errorcode        string `json:"errorcode,omitempty"`
	ErrorDescription string `json:"errorDescription,omitempty"`
	ErrorReference   string `json:"errorReference,omitempty"`
}

type ReqToken struct {
	Username string `json:"Username" validate:"required"`
	Password string `json:"Password" validate:"required"`
}

type ResponseToken struct {
	Code       string `json:"code"`
	Msg        string `json:"msg"`
	ReqeustId  string `json:"reqeustId,omitempty"`
	Token      string `json:"token,omitempty"`
	ExpireDate string `json:"expireDate,omitempty"`
	BackendUrl string `json:"backendUrl,omitempty"`
}
