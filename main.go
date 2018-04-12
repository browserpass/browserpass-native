package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/maximbaz/browserpass-native/openbsd"
	log "github.com/sirupsen/logrus"
)

// VERSION host app version
const VERSION = "3.0.0"

func main() {
	var verbose bool
	var version bool
	flag.BoolVar(&verbose, "v", false, "print verbose output")
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.Parse()

	if version {
		fmt.Println("Browserpass host app version:", VERSION)
		os.Exit(0)
	}

	openbsd.Pledge("stdio rpath proc exec")

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugf("Starting browserpass host app v%v", VERSION)
	process()
}
