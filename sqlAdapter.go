package alfajor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type sqlAdapter struct {
	db          *sql.DB
	insertQuery string
	config      *Config
}

func newSQLAdapter(config *Config, db *sql.DB, tableName string, dialect SQLDialect) *sqlAdapter {
	if tableName == "" {
		tableName = "alfajor_logs"
	}
	query := buildInsertQuery(tableName, dialect)

	return &sqlAdapter{
		db:          db,
		insertQuery: query,
		config:      config,
	}
}

func buildInsertQuery(tableName string, dialect SQLDialect) string {
	cols := "level, message, error, context, created_at"
	switch dialect {
	case DialectPostgres:
		return fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5)", tableName, cols)
	default:
		return fmt.Sprintf("INSERT INTO %s (%s) VALUES (?, ?, ?, ?, ?)", tableName, cols)
	}
}

func (s *sqlAdapter) insertLog(level string, msg string, err error, fields []map[string]interface{}) {

	var errStr sql.NullString
	if err != nil {
		errStr = sql.NullString{String: err.Error(), Valid: true}
	}

	var contextJSON []byte
	if len(fields) > 0 && fields[0] != nil {
		contextJSON, _ = json.Marshal(fields[0])
	}

	go func() {
		_, dbErr := s.db.Exec(s.insertQuery, level, msg, errStr, string(contextJSON), time.Now())
		if dbErr != nil {
			fmt.Printf("ALFAJOR ERROR: No se pudo escribir log en DB: %v\n", dbErr)
		}
	}()
}



func (s *sqlAdapter) Debug(msg string, fields ...map[string]interface{}) {
	s.insertLog(LevelDebug, msg, nil, fields)
}
func (s *sqlAdapter) Info(msg string, fields ...map[string]interface{}) {
	s.insertLog(LevelInfo, msg, nil, fields)
}
func (s *sqlAdapter) Warn(msg string, fields ...map[string]interface{}) {
	s.insertLog(LevelWarn, msg, nil, fields)
}
func (s *sqlAdapter) Error(msg string, err error, fields ...map[string]interface{}) {
	s.insertLog(LevelError, msg, err, fields)
}
func (s *sqlAdapter) Fatal(msg string, err error, fields ...map[string]interface{}) {
	s.insertLog(LevelFatal, msg, err, fields)
}
func (s *sqlAdapter) Sync() error {
	return nil 
}