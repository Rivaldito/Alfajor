package alfajor

import (
	"os"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapAdapter struct {
	logger *zap.Logger
}

func newZapAdapter(config *Config) (*zapAdapter, error) {
	var cores []zapcore.Core

	if config.EnableFile {
		fileEncoderConfig := zap.NewProductionEncoderConfig()
		fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.Rotation.Filename,
			MaxSize:    config.Rotation.MaxSize,
			MaxBackups: config.Rotation.MaxBackups,
			MaxAge:     config.Rotation.MaxAge,
			Compress:   config.Rotation.Compress,
			LocalTime:  config.Rotation.LocalTime,
		})

		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(fileEncoderConfig),
			fileWriter,
			config.getZapLevel(),
		)
		cores = append(cores, fileCore)
	}

	if config.EnableConsole {
		consoleEncoderConfig := zap.NewProductionEncoderConfig()
		consoleEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		var consoleEncoder zapcore.Encoder

		if config.Encoding == EncodingConsole {
			consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			consoleEncoder = zapcore.NewConsoleEncoder(consoleEncoderConfig)
		} else {
			consoleEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
			consoleEncoder = zapcore.NewJSONEncoder(consoleEncoderConfig)
		}

		consoleWriter := zapcore.Lock(os.Stdout)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			consoleWriter,
			config.getZapLevel(),
		)
		cores = append(cores, consoleCore)
	}

	combinedCore := zapcore.NewTee(cores...)
	var options []zap.Option
	if config.EnableCaller {
		options = append(options, zap.AddCaller(), zap.AddCallerSkip(1)) // Skip 1 para que no marque el wrapper
	}

	logger := zap.New(combinedCore, options...)
	return &zapAdapter{logger: logger}, nil
}

func (z *zapAdapter) mapToZapFields(fields ...map[string]interface{}) []zap.Field {
	if len(fields) == 0 || fields[0] == nil {
		return nil
	}
	zapFields := make([]zap.Field, 0, len(fields[0]))
	for k, v := range fields[0] {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

func (z *zapAdapter) Debug(msg string, fields ...map[string]interface{}) {
	z.logger.Debug(msg, z.mapToZapFields(fields...)...)
}
func (z *zapAdapter) Info(msg string, fields ...map[string]interface{}) {
	z.logger.Info(msg, z.mapToZapFields(fields...)...)
}
func (z *zapAdapter) Warn(msg string, fields ...map[string]interface{}) {
	z.logger.Warn(msg, z.mapToZapFields(fields...)...)
}
func (z *zapAdapter) Error(msg string, err error, fields ...map[string]interface{}) {
	allFields := z.mapToZapFields(fields...)
	allFields = append(allFields, zap.Error(err))
	z.logger.Error(msg, allFields...)
}
func (z *zapAdapter) Fatal(msg string, err error, fields ...map[string]interface{}) {
	allFields := z.mapToZapFields(fields...)
	allFields = append(allFields, zap.Error(err))
	z.logger.Fatal(msg, allFields...)
}
func (z *zapAdapter) Sync() error {
	return z.logger.Sync()
}