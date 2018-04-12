// +build nacl,plan9

package persistentlog

import log "github.com/sirupsen/logrus"

// AddPersistentLogHook configures persisting logs, not supported on these systems
func AddPersistentLogHook() {
	log.Warn("Persistent logging is not implemented on this OS")
}
