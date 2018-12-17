package server

import (
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

// PutHandler handles all put requests
func (s Server) PutHandler(w http.ResponseWriter, r *http.Request) {

	var res response

	r.Body = http.MaxBytesReader(w, r.Body, s.MaxUploadSize)
	res.Path = r.URL.RequestURI()
	res.Method = r.Method

	targetFile := path.Base(r.URL.Path)
	targetPath := path.Join(s.DocumentRoot, targetFile)
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {

		res.Status = http.StatusInternalServerError
		res.Error = err
		res.Message = "Failed to acquire the uploaded content"

		s.writeResponse(w, res)
		return
	}
	err = ioutil.WriteFile(targetPath, body, 0644)
	if err != nil {

		res.Status = http.StatusInternalServerError
		res.Error = err
		res.Message = "Failed to create file handle"

		s.writeResponse(w, res)
		return
	}

	res.Status = http.StatusOK
	res.Path = s.PathPrefix + strings.TrimPrefix(targetPath, s.DocumentRoot)
	res.Error = err
	res.Message = "File successfully uploaded"
	res.Hash = "n/a"
	s.writeResponse(w, res)
}
