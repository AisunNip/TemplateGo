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
	return InitLogger(appName, FpOutbound, targetSys)
}

func InitInboundLogger(appName string, targetSys LogSystem) *PatternLogger {
	return InitLogger(appName, FpInbound, targetSys)
}

func InitOutboundLogger(appName string, targetSys LogSystem) *PatternLogger {
	return InitLogger(appName, FpOutbound, targetSys)
}

func InitValidationLogger(appName string, targetSys LogSystem) *PatternLogger {
	return InitLogger(appName, FpValidation, targetSys)
}

func InitScheduleLogger(appName string, targetSys LogSystem) *PatternLogger {
	return InitLogger(appName, FpSchedule, targetSys)
}
