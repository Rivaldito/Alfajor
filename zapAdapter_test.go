package alfajor

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewLogger_FileCreation(t *testing.T) {

	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	cfg := NewDefaultConfig()
	cfg.EnableFile = true
	cfg.EnableConsole = false
	cfg.Rotation.Filename = logFile
	cfg.Encoding = EncodingJSON

	log, err := New(cfg, "zap")
	require.NoError(t, err, "La creación del logger no debería fallar")
	require.NotNil(t, log, "El logger no debería ser nulo")

	// Act
	testError := errors.New("error de prueba")
	log.Error("Este es un mensaje de error", testError, map[string]interface{}{"user_id": 123})
	err = log.Sync()
	require.NoError(t, err)

	// Assert
	content, err := os.ReadFile(logFile)
	require.NoError(t, err, "La lectura del fichero de log no debería fallar")

	logContent := string(content)
	assert.Contains(t, logContent, `"level":"ERROR"`, "El nivel de log debe ser ERROR")
	assert.Contains(t, logContent, `"msg":"Este es un mensaje de error"`, "El mensaje de log no es el esperado")
	assert.Contains(t, logContent, `"error":"error de prueba"`, "El error no está en el log")
	assert.Contains(t, logContent, `"user_id":123`, "El campo estructurado no está en el log")
	assert.Contains(t, logContent, `"caller":"alfajor/zapAdapter_test.go"`, "La información del caller no está presente")
}

func TestLogger_LevelsAndObserver(t *testing.T) {

	observedZapCore, observedLogs := observer.New(zapcore.DebugLevel)
	observedLogger := &zapAdapter{logger: zap.New(observedZapCore)}

	observedLogger.Debug("Mensaje de debug")
	observedLogger.Info("Mensaje de info", map[string]interface{}{"data": "value"})

	allLogs := observedLogs.All()
	require.Equal(t, 2, len(allLogs), "Deberían haberse registrado 2 logs")

	assert.Equal(t, zapcore.DebugLevel, allLogs[0].Level)
	assert.Equal(t, "Mensaje de debug", allLogs[0].Message)

	assert.Equal(t, zapcore.InfoLevel, allLogs[1].Level)
	assert.Equal(t, "Mensaje de info", allLogs[1].Message)

	expectedContext := []zapcore.Field{
		zap.Any("data", "value"),
	}
	assert.ElementsMatch(t, expectedContext, allLogs[1].Context, "Los campos de contexto no coinciden")
}

func TestLogger_Concurrency(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zapcore.DebugLevel)
	logger := &zapAdapter{logger: zap.New(observedZapCore)}

	var wg sync.WaitGroup
	logCount := 100
	wg.Add(logCount)

	for i := 0; i < logCount; i++ {
		go func(n int) {
			defer wg.Done()
			logger.Info("Log concurrente", map[string]interface{}{"goroutine_id": n})
		}(i)
	}

	wg.Wait()

	require.Equal(t, logCount, observedLogs.Len(), "El número de logs registrados no coincide con el esperado")
}
