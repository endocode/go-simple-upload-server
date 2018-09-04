package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

var logger *logrus.Logger

func run(args []string) int {
	bindAddress := flag.String("ip", "0.0.0.0", "IP address to bind")
	listenPort := flag.Int("port", 25478, "port number to listen on")
	// 5,242,880 bytes == 5 MiB
	maxUploadSize := flag.Int64("upload_limit", 5242880, "max size of uploaded file (byte)")
	logLevelFlag := flag.String("loglevel", "info", "logging level")
	pathPrefix := flag.String("pathPrefix", "/tmp", "inject path prefix")
	flag.Parse()
	serverRoot := flag.Arg(0)
	if len(serverRoot) == 0 {
		flag.Usage()
		return 2
	}
	serverRoot, err := filepath.Abs(serverRoot)
	if err != nil {
		return 2
	}

	if logLevel, err := logrus.ParseLevel(*logLevelFlag); err != nil {
		logrus.WithError(err).Error("failed to parse logging level, so set to default")
	} else {
		logger.Level = logLevel
	}

	logger.WithFields(logrus.Fields{
		"ip":           *bindAddress,
		"port":         *listenPort,
		"upload_limit": *maxUploadSize,
		"pathPrefix":   *pathPrefix,
		"root":         serverRoot,
	}).Info("start listening")
	server := NewServer(serverRoot, *maxUploadSize, *pathPrefix)
	http.Handle(*pathPrefix, server)
	http.ListenAndServe(fmt.Sprintf("%s:%d", *bindAddress, *listenPort), nil)
	return 0
}

func main() {
	logger = logrus.New()
	logger.Info("starting up simple-upload-server")

	result := run(os.Args)
	os.Exit(result)
}
