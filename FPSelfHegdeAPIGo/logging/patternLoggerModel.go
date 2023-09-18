package logging

type SetLogger struct {
	IsJSON    bool   `json:"isJson"`
	WriteFile bool   `json:"writeFile"`
	Path      string `json:"path"`
	FileName  string `json:"fileName"`
}

type PatternLogger struct {
	Level           LogLevel
	ApplicationName string
	ProductName     LogProductName
	SourceSystem    LogSystem
	TargetSystem    LogSystem
	SetLogger       SetLogger
}

type LogAppMessageBean struct {
	Timestamp       string         `json:"@timestamp"`
	ApplicationName string         `json:"@suffix"`
	LogType         LogType        `json:"logType"`
	CorrelationID   string         `json:"correlationID"`
	Level           LogLevel       `json:"level"`
	Message         string         `json:"message,omitempty"`
	MessageType     string         `json:"messageType,omitempty"`
	StackTrace      interface{}    `json:"stackTrace,omitempty"`
	ProductName     LogProductName `json:"productName,omitempty"`
	MonitorType     LogMonitorType `json:"monitorType,omitempty"`
	SourceSystem    LogSystem      `json:"sourceSystem"`
	SourceHostName  string         `json:"sourceHostName"`
	TargetSystem    LogSystem      `json:"targetSystem"`
}

type LogMonMessageBean struct {
	LogAppMessageBean
	TargetURL    string `json:"targetURL,omitempty"`
	Action       string `json:"action"`
	ElapsedTime  int64  `json:"elapsedTime"`
	ResponseCode string `json:"responseCode"`
}

type LogSecurityAuditBean struct {
	CorrelationID     string   `json:"txid"`
	Timestamp         string   `json:"@timestamp"`
	ApplicationName   string   `json:"@suffix"`
	Team              string   `json:"@team"`
	ClientIPAddr      string   `json:"endpoint"`
	LogType           LogType  `json:"log_cat"`
	Level             LogLevel `json:"level"`
	EmployeeID        string   `json:"employee_id"`
	Action            string   `json:"service_type"`
	ObjectName        string   `json:"method"`
	Request           string   `json:"request"`
	Response          string   `json:"response"`
	ResponseIndicator string   `json:"result_indicator"`
	Remark            string   `json:"remark"`
}

type logTypes string

func (logType logTypes) isType() logTypes {
	return logType
}

type LogType interface {
	isType() logTypes
}

const (
	Monitor       = logTypes("Monitor")
	AppLog        = logTypes("AppLog")
	SecurityAudit = logTypes("SecurityAudit")
	ServerLog     = logTypes("ServerLog")
)

type CallerInfoBean struct {
	ClassName  string `json:"className"`
	MethodName string `json:"methodName"`
	FileName   string `json:"fileName"`
}

type LogLevel string

const (
	LEVEL_ALL   = LogLevel("ALL")
	LEVEL_TRACE = LogLevel("TRACE")
	LEVEL_DEBUG = LogLevel("DEBUG")
	LEVEL_INFO  = LogLevel("INFO")
	LEVEL_WARN  = LogLevel("WARN")
	LEVEL_ERROR = LogLevel("ERROR")
	LEVEL_FATAL = LogLevel("FATAL")
	LEVEL_OFF   = LogLevel("OFF")
)

func (l LogLevel) Integer() int {
	var level int

	switch l {
	case LEVEL_ALL:
		level = 7
	case LEVEL_TRACE:
		level = 6
	case LEVEL_DEBUG:
		level = 5
	case LEVEL_INFO:
		level = 4
	case LEVEL_WARN:
		level = 3
	case LEVEL_ERROR:
		level = 2
	case LEVEL_FATAL:
		level = 1
	case LEVEL_OFF:
		level = 0
	default:
		level = 0
	}

	return level
}

type LogSystem string

const (
	FpUtil       LogSystem = "FP_UTIL"
	FpInbound    LogSystem = "FP_INBOUND"
	FpOutbound   LogSystem = "FP_OUTBOUND"
	FpValidation LogSystem = "FP_VALIDATION"
	FpSchedule   LogSystem = "FP_SCHEDULE"
	ICE          LogSystem = "ICE"
)

type productNames string

func (productName productNames) isProduct() productNames {
	return productName
}

type LogProductName interface {
	isProduct() productNames
}

const (
	TrueOnline  = productNames("TrueOnline")
	TrueMoveH   = productNames("TrueMoveH")
	TrueVisions = productNames("TrueVisions")
	IoT         = productNames("IoT")
	Convergence = productNames("Convergence")
	All         = productNames("All")
)

type monitorTypes string

func (monitorType monitorTypes) isMonitorType() monitorTypes {
	return monitorType
}

type LogMonitorType interface {
	isMonitorType() monitorTypes
}

const (
	Application         = monitorTypes("Application")
	FormProvider        = monitorTypes("FormProvider")
	FormClient          = monitorTypes("FormClient")
	XMLProvider         = monitorTypes("XMLProvider")
	XMLClient           = monitorTypes("XMLClient")
	WebServiceProvider  = monitorTypes("WebServiceProvider")
	WebServiceClient    = monitorTypes("WebServiceClient")
	RESTServiceProvider = monitorTypes("RESTServiceProvider")
	RESTServiceClient   = monitorTypes("RESTServiceClient")
	DatabaseClient      = monitorTypes("DatabaseClient")
)

type StackTrace struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}
