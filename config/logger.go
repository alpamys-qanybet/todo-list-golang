package config

var debugLog bool

func SetDebugLog(debugger bool) {
	debugLog = debugger
}

func DebugLog() bool {
	return debugLog
}
