package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/browserpass/browserpass-native/openbsd"
	"github.com/browserpass/browserpass-native/persistentlog"
	"github.com/browserpass/browserpass-native/request"
	"github.com/browserpass/browserpass-native/version"
	log "github.com/sirupsen/logrus"
)

func main() {
	var isVerbose bool
	var isVersion bool
	flag.BoolVar(&isVerbose, "v", false, "print verbose output")
	flag.BoolVar(&isVersion, "version", false, "print version and exit")
	flag.Parse()

	if isVersion {
		fmt.Println("Browserpass host app version:", version.String())
		os.Exit(0)
	}

	openbsd.Pledge("stdio rpath proc exec")

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	if isVerbose {
		log.SetLevel(log.DebugLevel)
	}

	persistentlog.AddPersistentLogHook()

	log.Debugf("Starting browserpass host app v%v", version.String())
	request.Process()
}
