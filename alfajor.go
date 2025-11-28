package alfajor

import (
	"database/sql"
	"fmt"
)

type Logger interface {
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Error(msg string, err error, fields ...map[string]interface{})
	Fatal(msg string, err error, fields ...map[string]interface{})
	Sync() error
}

type LoggerOption func(*loggerOptions)

type loggerOptions struct {
	dbInst    *sql.DB
	tableName string
	dialect   SQLDialect
}

func WithSQLDB(db *sql.DB, tableName string, dialect SQLDialect) LoggerOption {
	return func(o *loggerOptions) {
		o.dbInst = db
		o.tableName = tableName
		if dialect == "" {
			o.dialect = DialectMySQL
		} else {
			o.dialect = dialect
		}
	}
}


func New(config *Config, loggerType string, opts ...LoggerOption) (Logger, error) {

	options := &loggerOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var mainLogger Logger
	var err error


	switch loggerType {
	case "zap":
		mainLogger, err = newZapAdapter(config)
	case "sql":
		if options.dbInst == nil {
			return nil, fmt.Errorf("para logger tipo 'sql' es obligatorio usar WithSQLDB")
		}
		return newSQLAdapter(config, options.dbInst, options.tableName, options.dialect), nil
	default:
		return nil, fmt.Errorf("logger type '%s' no soportado", loggerType)
	}

	if err != nil {
		return nil, err
	}


	if options.dbInst != nil {
		sqlLog := newSQLAdapter(config, options.dbInst, options.tableName, options.dialect)
		return newMultiAdapter(mainLogger, sqlLog), nil
	}

	return mainLogger, nil
}