package logging

import (
	"FPSelfHegdeAPIGo/constant"
)

func NewLogger() *PatternLogger {
	logger := InitOutboundLogger(constant.AppName, FpOutbound)
	return logger
}
