// +build !windows,!nacl,!plan9

package persistentlog

import (
	"log/syslog"

	log "github.com/sirupsen/logrus"
	logSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

// AddPersistentLogHook configures persisting logs in syslog
func AddPersistentLogHook() {
	hook, err := logSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "browserpass")

	if err == nil {
		log.AddHook(hook)
	} else {
		log.Warn("Unable to connect to syslog, logs will NOT be persisted")
	}
}
