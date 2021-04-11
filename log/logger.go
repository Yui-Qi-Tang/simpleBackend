// Package builds a logger for this projcet. Just a wrapper on uber zap log
package log

import "go.uber.org/zap"

// Logger is the logger
var Logger *zap.Logger

// for now, just set it as default
func init() {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Logger = l

	// Logger.Info("success to set logger")
}
