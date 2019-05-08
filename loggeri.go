// Copyright 2015,2016,2017,2018 SeukWon Kang (kasworld@gmail.com)

package runnablemanager

type loggerI interface {
	// Reload() error
	Fatal(format string, v ...interface{})
	// Error(format string, v ...interface{})
	// Warn(format string, v ...interface{})
	// Debug(format string, v ...interface{})
	// TraceService(format string, v ...interface{})
}
