package config

var debugLog bool = true // default value

func SetDebugLog(debugger bool) {
	debugLog = debugger
}

func DebugLog() bool {
	return debugLog
}
