package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/endocode/go-simple-upload-server/logging"
	"github.com/endocode/go-simple-upload-server/server"

	"github.com/gorilla/mux"
)

func run(args []string) int {
	bindAddress := flag.String("ip", "0.0.0.0", "IP address to bind")
	listenPort := flag.Int("port", 25478, "port number to listen on")
	maxUploadSize := flag.Int64("upload_limit", 5242880, "max size of uploaded file (byte)")
	debugFlag := flag.Bool("debug", false, "Turn debug on/off")
	pathPrefix := flag.String("pathPrefix", "/moo", "inject path prefix")
	webSpace := flag.String("serverRoot", "/tmp/htdocs", "Webspace")
	flag.Parse()

	logger := logging.Start(*debugFlag)
	logger.Printf("[DEBUG] Debug mode is on")

	webSpaceAbsPath, err := filepath.Abs(*webSpace)
	if err != nil {
		return 2
	}

	if _, err := os.Stat(webSpaceAbsPath); os.IsNotExist(err) {
		logger.Printf(fmt.Sprintf("[FATAL] Server root is not available: %v", webSpaceAbsPath))
		return 2
	}

	server := server.New(webSpaceAbsPath, *maxUploadSize, *pathPrefix, logger)

	r := mux.NewRouter()

	r.HandleFunc("/health", server.HealthCheckHandler)
	r.PathPrefix(*pathPrefix).HandlerFunc(server.GetHandler).Methods("Get")
	r.PathPrefix(*pathPrefix).HandlerFunc(server.PostHandler).Methods("Post")
	r.PathPrefix(*pathPrefix).HandlerFunc(server.PutHandler).Methods("Put")

	r.Use(server.LoggingMiddleware)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *bindAddress, *listenPort))
	if err != nil {
		logger.Printf("[FATAL] Failed to initialize listener: %v", err)
	}

	logger.Printf(fmt.Sprintf("[INFO] Maximum upload size: %v bytes", *maxUploadSize))
	logger.Printf(fmt.Sprintf("[INFO] Server root: %v", webSpaceAbsPath))
	logger.Printf(fmt.Sprintf("[INFO] Path prefix: %v", *pathPrefix))
	logger.Printf(fmt.Sprintf("[INFO] Health endpoint on: http://%v:%v/health", *bindAddress, *listenPort))
	log.Fatal(http.Serve(listener, r))
	return 0
}

func main() {
	result := run(os.Args)
	os.Exit(result)
}
