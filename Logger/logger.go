package Logger

import (
	"fmt"
	"go.uber.org/zap"
)

var null = Logger{}

type Logger struct {
	Config zap.Config

	Log *zap.Logger

	SyncFunc func() error
}

func GetLogger(config zap.Config) (Logger, error) {
	var res Logger
	var err error

	res.Config = config
	res.Log, err = res.Config.Build()
	if err != nil {
		return null, err
	}
	res.SyncFunc = res.Log.Sync
	return res, nil
}
func (l *Logger) Debug(format string, values ...any) {
	l.Log.Debug(fmt.Sprintf(format, values...))
}
func (l *Logger) Info(format string, values ...any) {
	l.Log.Info(fmt.Sprintf(format, values...))
}
func (l *Logger) Warning(format string, values ...any) {
	l.Log.Warn(fmt.Sprintf(format, values...))
}
func (l *Logger) Error(format string, values ...any) {
	l.Log.Error(fmt.Sprintf(format, values...))
}
func (l *Logger) DPanic(format string, values ...any) {
	l.Log.DPanic(fmt.Sprintf(format, values...))
}
func (l *Logger) Panic(format string, values ...any) {
	l.Log.Panic(fmt.Sprintf(format, values...))
}
func (l *Logger) Fatal(format string, values ...any) {
	l.Log.Fatal(fmt.Sprintf(format, values...))
}
