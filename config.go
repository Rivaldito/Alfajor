package alfajor

import "go.uber.org/zap/zapcore"

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"

	EncodingJSON    = "json"
	EncodingConsole = "console"
)

type SQLDialect string

const (
	DialectMySQL    SQLDialect = "mysql"    
	DialectPostgres SQLDialect = "postgres" 
)

type RotationConfig struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
}

type Config struct {
	Level         string
	Encoding      string
	EnableCaller  bool
	EnableFile    bool
	EnableConsole bool
	Rotation      RotationConfig
}

func NewDefaultConfig() *Config {
	return &Config{
		Level:         LevelInfo,
		Encoding:      EncodingConsole,
		EnableCaller:  true,
		EnableFile:    false,
		EnableConsole: true,
		Rotation: RotationConfig{
			Filename:   "app.log",
			MaxSize:    100,
			MaxAge:     30,
			MaxBackups: 5,
			LocalTime:  true,
			Compress:   true,
		},
	}
}

func (c *Config) getZapLevel() zapcore.Level {
	switch c.Level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	case LevelInfo:
		fallthrough
	default:
		return zapcore.InfoLevel
	}
}