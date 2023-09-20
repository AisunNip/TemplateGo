package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var sourceHostName string

func init() {
	initialHostName()
}

const (
	PermFileMode os.FileMode = 0666
	TimeFormat   string      = "2006-01-02T15:04:05.000Z07:00"

	MaxStackLength int = 40
)

func initialHostName() {
	if sourceHostName == "" {
		var err error
		sourceHostName, err = os.Hostname()
		if err != nil {
			panic(fmt.Errorf("PatternLogger get host name error: %s\n", err))
		}
	}
}

func (p *PatternLogger) AllowLogging(level LogLevel) bool {
	var isAllow bool

	if p.Level == LEVEL_ALL {
		isAllow = true
	} else if p.Level == LEVEL_OFF {
		isAllow = false
	} else {
		if p.Level.Integer() >= level.Integer() {
			isAllow = true
		}
	}

	return isAllow
}

func getCallersFrames(skipNoOfStack int) *runtime.Frames {
	stackBuf := make([]uintptr, MaxStackLength)
	length := runtime.Callers(skipNoOfStack, stackBuf[:])
	stack := stackBuf[:length]

	return runtime.CallersFrames(stack)
}

func getStackTraceError(err error) string {
	frames := getCallersFrames(5)
	trace := err.Error()

	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			trace = trace + fmt.Sprintf("\n File: %s, Line: %d, Func: %s", frame.File, frame.Line, frame.Function)
		}
		if !more {
			break
		}
	}

	return trace
}

func GetStackTrace() (stacktrace StackTrace) {
	frames := getCallersFrames(4)

	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			stacktrace.Function = frame.Function
			stacktrace.File = frame.File
			stacktrace.Line = frame.Line
			break
		}
		if !more {
			break
		}
	}

	return stacktrace
}

func getCallerInfoBean(stacktrace StackTrace) (callerInfoBean CallerInfoBean) {
	runes := []rune(stacktrace.Function)
	indexFirst := strings.Index(stacktrace.Function, ".")
	lastChar := string(stacktrace.Function[len(stacktrace.Function)-1:])
	indexLast := strings.LastIndex(stacktrace.Function, lastChar)
	callerInfoBean.ClassName = string(runes[0:indexFirst])
	callerInfoBean.FileName = fmt.Sprintf("\n\t %s:%d", stacktrace.File, stacktrace.Line)
	callerInfoBean.MethodName = string(runes[indexFirst+1 : indexLast+1])
	return callerInfoBean
}

func getElapsedTime(start time.Time) int64 {
	var elapsedTime int64 = 0

	if !start.IsZero() {
		elapsedTime = time.Since(start).Milliseconds()
	}

	return elapsedTime
}

func (p *PatternLogger) EnableFileLogger(path string, fileName string) {
	p.SetLogger.WriteFile = true
	p.SetLogger.Path = path
	p.SetLogger.FileName = fileName
}

func (p *PatternLogger) logMonitoring(messageType string, correlationID string, monitorType LogMonitorType,
	targetURL string, action string, elapsedTime int64, responseCode string) time.Time {

	currDateTime := time.Now()

	var messageBean = LogMonMessageBean{}
	messageBean.Timestamp = currDateTime.Format(TimeFormat)
	messageBean.ApplicationName = p.ApplicationName
	messageBean.LogType = Monitor
	messageBean.Level = LEVEL_INFO
	messageBean.CorrelationID = correlationID
	messageBean.MessageType = messageType
	messageBean.ProductName = p.ProductName
	messageBean.MonitorType = monitorType
	messageBean.SourceSystem = p.SourceSystem
	messageBean.SourceHostName = sourceHostName
	messageBean.TargetSystem = p.TargetSystem
	messageBean.TargetURL = targetURL
	messageBean.Action = action
	messageBean.ElapsedTime = elapsedTime
	messageBean.ResponseCode = responseCode

	if p.SetLogger.IsJSON {
		jsonBinary, _ := json.Marshal(&messageBean)

		if p.SetLogger.WriteFile {
			p.writeLogToFile(string(jsonBinary))
		} else {
			fmt.Println(string(jsonBinary))
		}
	} else {
		if p.SetLogger.WriteFile {
			p.writeLogToFile(logMonStringPattern(&messageBean))
		} else {
			fmt.Fprintln(os.Stdout, logMonStringPattern(&messageBean))
		}
	}

	return currDateTime
}

func (p *PatternLogger) logSecurityAudit(correlationID string, clientIPAddr string,
	employeeID string, action string, objectName string, request interface{}, oldValue interface{},
	response interface{}, isSuccess bool, remark string) {

	var securityAuditBean = LogSecurityAuditBean{}
	securityAuditBean.CorrelationID = correlationID
	securityAuditBean.Timestamp = time.Now().Format(TimeFormat)
	securityAuditBean.ApplicationName = p.ApplicationName
	securityAuditBean.Team = "CRM"
	securityAuditBean.ClientIPAddr = clientIPAddr
	securityAuditBean.LogType = SecurityAudit
	securityAuditBean.Level = LEVEL_INFO
	securityAuditBean.EmployeeID = employeeID
	securityAuditBean.Action = action
	securityAuditBean.ObjectName = objectName

	reqBinary, reqErr := json.Marshal(request)

	if reqErr == nil {
		if securityAuditBean.Action == "Modify" {
			securityAuditBean.Request = "New Value: " + string(reqBinary) + ", Old Value: "

			if oldValue != nil {
				oldValueBinary, oldValueErr := json.Marshal(oldValue)

				if oldValueErr == nil {
					securityAuditBean.Request = securityAuditBean.Request + string(oldValueBinary)
				}
			}
		} else if securityAuditBean.Action == "Delete" {
			securityAuditBean.Request = string(reqBinary) + ", Old Value: "

			if oldValue != nil {
				oldValueBinary, oldValueErr := json.Marshal(oldValue)

				if oldValueErr == nil {
					securityAuditBean.Request = securityAuditBean.Request + string(oldValueBinary)
				}
			}
		} else {
			securityAuditBean.Request = string(reqBinary)
		}
	}

	respBinary, respErr := json.Marshal(response)
	if respErr == nil {
		securityAuditBean.Response = string(respBinary)
	}

	if isSuccess {
		securityAuditBean.ResponseIndicator = "SUCCESS"
	} else {
		securityAuditBean.ResponseIndicator = "FAIL"
	}

	securityAuditBean.Remark = remark

	if p.SetLogger.IsJSON {
		jsonBinary, _ := json.Marshal(&securityAuditBean)

		if p.SetLogger.WriteFile {
			p.writeLogToFile(string(jsonBinary))
		} else {
			fmt.Println(string(jsonBinary))
		}
	} else {
		if p.SetLogger.WriteFile {
			p.writeLogToFile(logSecurityAuditStringPattern(&securityAuditBean))
		} else {
			fmt.Fprintln(os.Stdout, logSecurityAuditStringPattern(&securityAuditBean))
		}
	}
}

func (p *PatternLogger) SecurityAuditView(correlationID string, clientIPAddr string,
	employeeID string, objectName string, request interface{}, response interface{},
	isSuccess bool, remark string) {
	p.logSecurityAudit(correlationID, clientIPAddr, employeeID, "View",
		objectName, request, nil, response, isSuccess, remark)
}

func (p *PatternLogger) SecurityAuditCreate(correlationID string, clientIPAddr string,
	employeeID string, objectName string, request interface{}, response interface{},
	isSuccess bool, remark string) {
	p.logSecurityAudit(correlationID, clientIPAddr, employeeID, "Create",
		objectName, request, nil, response, isSuccess, remark)
}

func (p *PatternLogger) SecurityAuditDelete(correlationID string, clientIPAddr string,
	employeeID string, objectName string, request interface{}, oldValue interface{},
	response interface{}, isSuccess bool, remark string) {
	p.logSecurityAudit(correlationID, clientIPAddr, employeeID, "Delete",
		objectName, request, oldValue, response, isSuccess, remark)
}

func (p *PatternLogger) SecurityAuditModify(correlationID string, clientIPAddr string,
	employeeID string, objectName string, request interface{}, oldValue interface{},
	response interface{}, isSuccess bool, remark string) {
	p.logSecurityAudit(correlationID, clientIPAddr, employeeID, "Modify",
		objectName, request, oldValue, response, isSuccess, remark)
}

func (p *PatternLogger) SecurityAuditExport(correlationID string, clientIPAddr string,
	employeeID string, objectName string, request interface{}, response interface{},
	isSuccess bool, remark string) {
	p.logSecurityAudit(correlationID, clientIPAddr, employeeID, "Export",
		objectName, request, nil, response, isSuccess, remark)
}

func (p *PatternLogger) LogRequest(correlationID string, monitorType LogMonitorType,
	targetURL string, action string) time.Time {
	return p.logMonitoring("Request", correlationID, monitorType, targetURL, action,
		getElapsedTime(time.Time{}), "")
}

func (p *PatternLogger) LogRequestApplication(correlationID string) time.Time {
	callerInfoBean := getCallerInfoBean(GetStackTrace())

	return p.logMonitoring("Request", correlationID, Application, "",
		"package: "+callerInfoBean.ClassName+", func: "+callerInfoBean.MethodName,
		getElapsedTime(time.Time{}), "")
}

func (p *PatternLogger) LogRequestFormProvider(correlationID string) time.Time {
	callerInfoBean := getCallerInfoBean(GetStackTrace())

	return p.logMonitoring("Request", correlationID, FormProvider, "",
		"package: "+callerInfoBean.ClassName+", func: "+callerInfoBean.MethodName,
		getElapsedTime(time.Time{}), "")
}

func (p *PatternLogger) LogRequestFormClient(correlationID string, targetURL string, action string) time.Time {
	return p.logMonitoring("Request", correlationID, FormClient, targetURL,
		action, getElapsedTime(time.Time{}), "")
}

func (p *PatternLogger) LogRequestXMLProvider(correlationID string) time.Time {
	callerInfoBean := getCallerInfoBean(GetStackTrace())

	return p.logMonitoring("Request", correlationID, XMLProvider, "",
		"package: "+callerInfoBean.ClassName+", func: "+callerInfoBean.MethodName,
		getElapsedTime(time.Time{}), "")
}

func (p *PatternLogger) LogRequestXMLClient(correlationID string, targetURL string, action string) time.Time {
	return p.logMonitoring("Request", correlationID, XMLClient, targetURL,
		action, getElapsedTime(time.Time{}), "")
}

func (p *PatternLogger) LogRequestWSProvider(correlationID string) time.Time {
	callerInfoBean := getCallerInfoBean(GetStackTrace())

	return p.logMonitoring("Request", correlationID, WebServiceProvider, "",
		"package: "+callerInfoBean.ClassName+", func: "+callerInfoBean.MethodName,
		getElapsedTime(time.Time{}), "")
}

func (p *PatternLogger) LogRequestWSClient(correlationID string, targetURL string, action string) time.Time {
	return p.logMonitoring("Request", correlationID, WebServiceClient, targetURL,
		action, getElapsedTime(time.Time{}), "")
}

func (p *PatternLogger) LogRequestRESTProvider(correlationID string) time.Time {
	callerInfoBean := getCallerInfoBean(GetStackTrace())

	return p.logMonitoring("Request", correlationID, RESTServiceProvider, "",
		"package: "+callerInfoBean.ClassName+", func: "+callerInfoBean.MethodName,
		getElapsedTime(time.Time{}), "")
}

func (p *PatternLogger) LogRequestRESTClient(correlationID string, targetURL string, action string) time.Time {
	return p.logMonitoring("Request", correlationID, RESTServiceClient, targetURL,
		action, getElapsedTime(time.Time{}), "")
}

func (p *PatternLogger) LogRequestDBClient(correlationID string) time.Time {
	callerInfoBean := getCallerInfoBean(GetStackTrace())

	return p.logMonitoring("Request", correlationID, DatabaseClient, "",
		"package: "+callerInfoBean.ClassName+", func: "+callerInfoBean.MethodName,
		getElapsedTime(time.Time{}), "")
}

func (p *PatternLogger) LogResponse(correlationID string, monitorType LogMonitorType, targetURL string,
	action string, responseCode string, requestDateTime time.Time) time.Time {
	return p.logMonitoring("Response", correlationID, monitorType, targetURL,
		action, getElapsedTime(requestDateTime), responseCode)
}

func (p *PatternLogger) LogResponseApplication(correlationID string, responseCode string, requestDateTime time.Time) time.Time {
	callerInfoBean := getCallerInfoBean(GetStackTrace())

	return p.logMonitoring("Response", correlationID, Application, "",
		"package: "+callerInfoBean.ClassName+", func: "+callerInfoBean.MethodName,
		getElapsedTime(requestDateTime), responseCode)
}

func (p *PatternLogger) LogResponseFormProvider(correlationID string, responseCode string, requestDateTime time.Time) time.Time {
	callerInfoBean := getCallerInfoBean(GetStackTrace())

	return p.logMonitoring("Response", correlationID, FormProvider, "",
		"package: "+callerInfoBean.ClassName+", func: "+callerInfoBean.MethodName,
		getElapsedTime(requestDateTime), responseCode)
}

func (p *PatternLogger) LogResponseFormClient(correlationID string, targetURL string,
	action string, responseCode string, requestDateTime time.Time) time.Time {
	return p.logMonitoring("Response", correlationID, FormClient, targetURL,
		action, getElapsedTime(requestDateTime), responseCode)
}

func (p *PatternLogger) LogResponseXMLProvider(correlationID string, responseCode string, requestDateTime time.Time) time.Time {
	callerInfoBean := getCallerInfoBean(GetStackTrace())

	return p.logMonitoring("Response", correlationID, XMLProvider, "",
		"package: "+callerInfoBean.ClassName+", func: "+callerInfoBean.MethodName,
		getElapsedTime(requestDateTime), responseCode)
}

func (p *PatternLogger) LogResponseXMLClient(correlationID string, targetURL string,
	action string, responseCode string, requestDateTime time.Time) time.Time {
	return p.logMonitoring("Response", correlationID, XMLClient, targetURL,
		action, getElapsedTime(requestDateTime), responseCode)
}

func (p *PatternLogger) LogResponseWSProvider(correlationID string, responseCode string, requestDateTime time.Time) time.Time {
	callerInfoBean := getCallerInfoBean(GetStackTrace())

	return p.logMonitoring("Response", correlationID, WebServiceProvider, "",
		"package: "+callerInfoBean.ClassName+", func: "+callerInfoBean.MethodName,
		getElapsedTime(requestDateTime), responseCode)
}

func (p *PatternLogger) LogResponseWSClient(correlationID string, targetURL string, action string, responseCode string,
	requestDateTime time.Time) time.Time {
	return p.logMonitoring("Response", correlationID, WebServiceClient, targetURL,
		action, getElapsedTime(requestDateTime), responseCode)
}

func (p *PatternLogger) LogResponseRESTProvider(correlationID string, responseCode string, requestDateTime time.Time) time.Time {
	callerInfoBean := getCallerInfoBean(GetStackTrace())

	return p.logMonitoring("Response", correlationID, RESTServiceProvider, "",
		"package: "+callerInfoBean.ClassName+", func: "+callerInfoBean.MethodName,
		getElapsedTime(requestDateTime), responseCode)
}

func (p *PatternLogger) LogResponseRESTClient(correlationID string, targetURL string, action string,
	responseCode string, requestDateTime time.Time) time.Time {
	return p.logMonitoring("Response", correlationID, RESTServiceClient, targetURL,
		action, getElapsedTime(requestDateTime), responseCode)
}

func (p *PatternLogger) LogResponseDBClient(correlationID string, responseCode string, requestDateTime time.Time) time.Time {
	callerInfoBean := getCallerInfoBean(GetStackTrace())

	return p.logMonitoring("Response", correlationID, DatabaseClient, "",
		"package: "+callerInfoBean.ClassName+", func: "+callerInfoBean.MethodName,
		getElapsedTime(requestDateTime), responseCode)
}

func (p *PatternLogger) logApp(correlationID string, level LogLevel, message string, stackTrace string) {
	if !p.AllowLogging(level) {
		return
	}

	var messageBean = LogAppMessageBean{}
	messageBean.Timestamp = time.Now().Format(TimeFormat)
	messageBean.ApplicationName = p.ApplicationName
	messageBean.LogType = AppLog
	messageBean.CorrelationID = correlationID
	messageBean.Level = level
	messageBean.Message = message
	messageBean.StackTrace = stackTrace
	messageBean.SourceHostName = sourceHostName
	messageBean.ProductName = p.ProductName
	messageBean.SourceSystem = p.SourceSystem
	messageBean.TargetSystem = p.TargetSystem

	if p.SetLogger.IsJSON {
		jsonBinary, _ := json.Marshal(&messageBean)

		if p.SetLogger.WriteFile {
			p.writeLogToFile(string(jsonBinary))
		} else {
			fmt.Println(string(jsonBinary))
		}
	} else {
		if p.SetLogger.WriteFile {
			p.writeLogToFile(logAppStringPattern(&messageBean))
		} else {
			fmt.Fprintln(os.Stdout, logAppStringPattern(&messageBean))
		}
	}
}

func (p *PatternLogger) WriteRequestMsg(correlationID string, url string, httpMethod string, message interface{}) {
	if message == nil {
		p.logApp(correlationID, LEVEL_INFO, "Request URL: "+url+
			", HttpMethod: "+httpMethod+", Request Message: null", "")
	} else {
		msg, err := json.Marshal(message)
		if err != nil {
			p.Error(correlationID, err.Error(), err)
		}
		p.logApp(correlationID, LEVEL_INFO, "Request URL: "+url+
			", HttpMethod: "+httpMethod+", Request Message: "+string(msg), "")
	}
}

func (p *PatternLogger) WriteResponseMsg(correlationID string, message interface{}) {
	if message == nil {
		p.logApp(correlationID, LEVEL_INFO, "Response Message: null", "")
	} else {
		msg, err := json.Marshal(message)
		if err != nil {
			p.Error(correlationID, err.Error(), err)
		}
		p.logApp(correlationID, LEVEL_INFO, "Response Message: "+string(msg), "")
	}
}

func (p *PatternLogger) Info(correlationID string, message string, args ...interface{}) {
	if len(args) > 0 {
		msg, stackTrace := checkArguments(args)
		p.logApp(correlationID, LEVEL_INFO, message+msg, stackTrace)
	} else {
		p.logApp(correlationID, LEVEL_INFO, message, "")
	}
}

func (p *PatternLogger) Fatal(correlationID string, message string, args ...interface{}) {
	if len(args) > 0 {
		msg, stackTrace := checkArguments(args)
		p.logApp(correlationID, LEVEL_FATAL, message+msg, stackTrace)
	} else {
		p.logApp(correlationID, LEVEL_FATAL, message, "")
	}
}

func (p *PatternLogger) Error(correlationID string, message string, args ...interface{}) {
	if len(args) > 0 {
		msg, stackTrace := checkArguments(args)
		p.logApp(correlationID, LEVEL_ERROR, message+msg, stackTrace)
	} else {
		p.logApp(correlationID, LEVEL_ERROR, message, "")
	}
}

func (p *PatternLogger) Warn(correlationID string, message string, args ...interface{}) {
	if len(args) > 0 {
		msg, stackTrace := checkArguments(args)
		p.logApp(correlationID, LEVEL_WARN, message+msg, stackTrace)
	} else {
		p.logApp(correlationID, LEVEL_WARN, message, "")
	}
}

func (p *PatternLogger) Debug(correlationID string, message string, args ...interface{}) {
	if len(args) > 0 {
		msg, stackTrace := checkArguments(args)
		p.logApp(correlationID, LEVEL_DEBUG, message+msg, stackTrace)
	} else {
		p.logApp(correlationID, LEVEL_DEBUG, message, "")
	}
}

func (p *PatternLogger) Trace(correlationID string, message string, args ...interface{}) {
	if len(args) > 0 {
		msg, stackTrace := checkArguments(args)
		p.logApp(correlationID, LEVEL_TRACE, message+msg, stackTrace)
	} else {
		p.logApp(correlationID, LEVEL_TRACE, message, "")
	}
}

func checkArguments(args []interface{}) (message string, staceTrace string) {
	for i := 0; i < len(args); i++ {
		switch v := args[i].(type) {
		case float32, float64, complex64, complex128:
			message = message + " " + fmt.Sprintf("%g", v)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			message = message + " " + fmt.Sprintf("%d", v)
		case bool:
			message = message + " " + strconv.FormatBool(v)
		case string:
			message = message + " " + v
		case error:
			staceTrace = staceTrace + " " + getStackTraceError(v)
		case *strconv.NumError:
			staceTrace = staceTrace + " " + getStackTraceError(v)
		default:
			str, _ := json.Marshal(v)
			message = message + " " + string(str)
		}
	}
	return
}

func logAppStringPattern(bean *LogAppMessageBean) string {
	return fmt.Sprintf(`timestamp=%s|logType=%s|correlationID=%s|level=%v|message=%s|messageType=%s|stackTrace=%v|productName=%s|monitorType=%v|sourceSystem=%s|sourceHostName=%s|targetSystem=%s|targetURL=|action=|elapsedTime=|responseCode=`,
		bean.Timestamp, bean.LogType, bean.CorrelationID, bean.Level, bean.Message, bean.MessageType, bean.StackTrace, bean.ProductName,
		bean.MonitorType, bean.SourceSystem, bean.SourceHostName, bean.TargetSystem)
}

func logMonStringPattern(bean *LogMonMessageBean) string {
	return fmt.Sprintf(`timestamp=%s|logType=%s|correlationID=%s|level=%v|message=%s|messageType=%s|stackTrace=%v|productName=%s|monitorType=%v|sourceSystem=%s|sourceHostName=%s|targetSystem=%s|targetURL=%s|action=%s|elapsedTime=%d|responseCode=%s`,
		bean.Timestamp, bean.LogType, bean.CorrelationID, bean.Level, bean.Message, bean.MessageType, bean.StackTrace, bean.ProductName,
		bean.MonitorType, bean.SourceSystem, bean.SourceHostName, bean.TargetSystem, bean.TargetURL, bean.Action, bean.ElapsedTime, bean.ResponseCode)
}

func logSecurityAuditStringPattern(bean *LogSecurityAuditBean) string {
	return fmt.Sprintf(`timestamp=%s|logType=%s|correlationID=%s|level=%v|team=%s|clientIPAddr=%s|employeeID=%s|action=%s|objectName=%s|request=%s|response=%s|responseIndicator=%s|remark=%s`,
		bean.Timestamp, bean.LogType, bean.CorrelationID, bean.Level, bean.Team, bean.ClientIPAddr, bean.EmployeeID,
		bean.Action, bean.ObjectName, bean.Request, bean.Response, bean.ResponseIndicator, bean.Remark)
}

func (p *PatternLogger) writeLogToFile(message string) {
	fileName := p.SetLogger.Path + "/" + p.SetLogger.FileName + "_" + time.Now().Format("2006-01-02") + ".log"
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, PermFileMode)

	if err != nil {
		fmt.Printf("Opening log file error: %v", err)
	}

	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.LstdFlags)
	log.Println(message)
}
