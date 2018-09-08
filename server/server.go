package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Server is a struct that represents the ServerConfig
type Server struct {
	DocumentRoot  string
	MaxUploadSize int64
	PathPrefix    string
	Logger        *log.Logger
}

type response struct {
	Status  int    `json:",omitempty"`
	Method  string `json:",omitempty"`
	Error   error  `json:",omitempty"`
	Message string `json:",omitempty"`
	Path    string `json:",omitempty"`
	Hash    string `json:",omitempty"`
}

var err error

// New creates a new simple-upload server.
func New(documentRoot string, maxUploadSize int64, pathPrefix string, log *log.Logger) Server {
	return Server{
		DocumentRoot:  documentRoot,
		MaxUploadSize: maxUploadSize,
		PathPrefix:    strings.TrimSuffix(pathPrefix, "/"),
		Logger:        log,
	}
}

// LoggingMiddleware injects a debug logger
func (s Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Printf(fmt.Sprintf("[DEBUG] %v %v %v", r.Method, r.RequestURI, r.Header.Get("Content-Type")))
		next.ServeHTTP(w, r)
	})
}

func (s Server) writeResponse(w http.ResponseWriter, res response) {
	if res.Error != nil {
		s.Logger.Printf(fmt.Sprintf("[ERROR] %v %v", res.Message, res.Error))
	}
	w.WriteHeader(res.Status)
	w.Header().Set("Content-Type", "application/json")

	b, e := json.Marshal(res)
	if e != nil {
		w.Write([]byte{})
	}
	w.Write(b)
}

func getSize(content io.Seeker) (int64, error) {
	size, err := content.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}
	_, err = content.Seek(0, os.SEEK_SET)
	if err != nil {
		return 0, err
	}
	return size, nil
}
