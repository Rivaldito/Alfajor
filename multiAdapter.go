package alfajor

import (
	"fmt"
	"strings"
)

type multiAdapter struct {
	loggers []Logger
}

func newMultiAdapter(loggers ...Logger) Logger {
	return &multiAdapter{loggers: loggers}
}

func (m *multiAdapter) Debug(msg string, fields ...map[string]interface{}) {
	for _, l := range m.loggers {
		l.Debug(msg, fields...)
	}
}
func (m *multiAdapter) Info(msg string, fields ...map[string]interface{}) {
	for _, l := range m.loggers {
		l.Info(msg, fields...)
	}
}
func (m *multiAdapter) Warn(msg string, fields ...map[string]interface{}) {
	for _, l := range m.loggers {
		l.Warn(msg, fields...)
	}
}
func (m *multiAdapter) Error(msg string, err error, fields ...map[string]interface{}) {
	for _, l := range m.loggers {
		l.Error(msg, err, fields...)
	}
}
func (m *multiAdapter) Fatal(msg string, err error, fields ...map[string]interface{}) {
	for _, l := range m.loggers {
		l.Fatal(msg, err, fields...)
	}
}
func (m *multiAdapter) Sync() error {
	var errs []string
	for _, l := range m.loggers {
		if err := l.Sync(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("sync errors: %s", strings.Join(errs, "; "))
	}
	return nil
}
