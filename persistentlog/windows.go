// +build windows

package persistentlog

import (
	"os"
	"path/filepath"

	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

// AddPersistentLogHook configures persisting logs in a file
func AddPersistentLogHook() {
	appDataPath := os.Getenv("LOCALAPPDATA")
	if appDataPath == "" {
		log.Warn("Unable to determine %%APPDATA%% folder location, logs will NOT be persisted")
		return
	}
	logFolderPath := filepath.Join(appDataPath, "browserpass")
	if err := os.MkdirAll(logFolderPath, os.ModePerm); err != nil {
		log.Warn("Unable to create browserpass folder in %%APPDATA%%, logs will NOT be persisted")
		return
	}
	logFilePath := filepath.Join(logFolderPath, "browserpass.log")
	log.Debug("Logs will being written to: ", logFilePath)
	log.AddHook(lfshook.NewHook(logFilePath, &log.TextFormatter{FullTimestamp: true}))
}
