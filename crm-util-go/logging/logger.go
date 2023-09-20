package logging

func InitLogger(appName string, sourceSys LogSystem, targetSys LogSystem) *PatternLogger {
	var logger = new(PatternLogger)
	logger.SetLogger.IsJSON = true
	logger.Level = LEVEL_INFO
	logger.ApplicationName = appName
	logger.ProductName = All
	logger.SourceSystem = sourceSys
	logger.TargetSystem = targetSys
	return logger
}

func InitUtilLogger(appName string, targetSys LogSystem) *PatternLogger {
	return InitLogger(appName, CrmUtil, targetSys)
}

func InitInboundLogger(appName string, targetSys LogSystem) *PatternLogger {
	return InitLogger(appName, CrmInbound, targetSys)
}

func InitOutboundLogger(appName string, targetSys LogSystem) *PatternLogger {
	return InitLogger(appName, CrmOutbound, targetSys)
}

func InitValidationLogger(appName string, targetSys LogSystem) *PatternLogger {
	return InitLogger(appName, CrmValidation, targetSys)
}

func InitScheduleLogger(appName string, targetSys LogSystem) *PatternLogger {
	return InitLogger(appName, CrmSchedule, targetSys)
}