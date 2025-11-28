package alfajor_test

import (
	"errors"
	"path/filepath"
	"testing"

	alfajor "github.com/Rivaldito/alfajor"
)

func TestMain(t *testing.T) {
	cfg := alfajor.NewDefaultConfig()

	logFolder := "app-logs"
	logFileName := "alfajor_app.log"

	logPath := filepath.Join(logFolder, logFileName)

	cfg.EnableFile = true
	cfg.Rotation.Filename = logPath

	cfg.Level = alfajor.LevelDebug
	cfg.Rotation.MaxSize = 5

	log, err := alfajor.New(cfg, "zap")
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer log.Sync()

	log.Info("Logger configurado para escribir en la carpeta 'app-logs'")
	log.Error("Este error también irá a un archivo dentro de 'app-logs'", errors.New("error de prueba"))
}
