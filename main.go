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
	webspace := flag.String("serverRoot", "/tmp/htdocs", "Webspace")
	flag.Parse()

	logger := logging.Start(*debugFlag)

	serverRoot, err := filepath.Abs(*webspace)
	if err != nil {
		return 2
	}
	if _, err := os.Stat(serverRoot); os.IsNotExist(err) {
		logger.Printf(fmt.Sprintf("[FATAL] Server root is not available: %v", serverRoot))
		return 2

	}
	logger.Printf(fmt.Sprintf("[INFO] Server root: %v", serverRoot))

	server := server.New(serverRoot, *maxUploadSize, *pathPrefix, logger)

	r := mux.NewRouter()

	r.HandleFunc("/health", server.HealthCheckHandler)
	r.PathPrefix(*pathPrefix).HandlerFunc(server.GetHandler).Methods("Get")
	r.PathPrefix(*pathPrefix).HandlerFunc(server.PostHandler).Methods("Post")
	r.PathPrefix(*pathPrefix).HandlerFunc(server.PostHandler).Methods("Put")

	r.Use(server.LoggingMiddleware)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *bindAddress, *listenPort))
	if err != nil {
		panic(err)
	}

	log.Fatal(http.Serve(listener, r))
	return 0
}

func main() {
	result := run(os.Args)
	os.Exit(result)
}
