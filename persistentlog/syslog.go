// +build !windows,!nacl,!plan9

package persistentlog

import (
	"log/syslog"

	log "github.com/sirupsen/logrus"
	logSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

// AddPersistentLogHook configures persisting logs in syslog
func AddPersistentLogHook() {
	if hook, err := logSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "browserpass"); err != nil {
		log.Warn("Unable to connect to syslog, logs will NOT be persisted: ", err)
	} else {
		log.AddHook(hook)
	}
}
