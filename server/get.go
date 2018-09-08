package server

import (
	"net/http"
)

// GetHandler handles all GET commands
func (s Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	h := http.StripPrefix(s.PathPrefix, http.FileServer(http.Dir(s.DocumentRoot)))
	h.ServeHTTP(w, r)
}
