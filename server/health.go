package server

import "net/http"

// HealthCheckHandler is a health endpoint, returns a 200
func (s Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	var res response
	res.Path = r.URL.RequestURI()
	res.Method = r.Method

	res.Status = http.StatusOK
	res.Error = nil
	res.Message = "Alive"

	s.writeResponse(w, res)
}
