package server

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// PutHandler handles all put requests
func (s Server) PutHandler(w http.ResponseWriter, r *http.Request) {
	var res response
	res.Path = r.URL.RequestURI()
	res.Method = r.Method

	targetFile := path.Base(r.URL.Path)
	targetPath := path.Join(s.DocumentRoot, targetFile)

	file, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {

		res.Status = http.StatusInternalServerError
		res.Error = err
		res.Message = "Failed to open the file"

		s.writeResponse(w, res)
	}
	defer file.Close()
	defer r.Body.Close()

	srcFile, _, err := r.FormFile("file")
	if err != nil {

		res.Status = http.StatusInternalServerError
		res.Error = err
		res.Message = "Failed to acquire the uploaded content"

		s.writeResponse(w, res)
		return
	}
	defer srcFile.Close()

	size, err := getSize(srcFile)
	if err != nil {

		res.Status = http.StatusInternalServerError
		res.Error = err
		res.Message = "Failed to get the size of the uploaded content"

		s.writeResponse(w, res)
		return
	}

	if size > s.MaxUploadSize {

		res.Status = http.StatusRequestEntityTooLarge
		res.Error = errors.New("Uploaded file size exceeds the limit")
		res.Message = "Uploaded file size exceeds the limit"

		s.writeResponse(w, res)
		return
	}

	_, err = io.Copy(file, srcFile)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Error = err
		res.Message = "Failed to read the uploaded file"

		s.writeResponse(w, res)
		return
	}

	writtenFile, err := os.Open(targetPath)
	if err != nil {

		res.Status = http.StatusInternalServerError
		res.Error = err
		res.Message = "Failed to read the newly written file"

		s.writeResponse(w, res)
		return
	}
	defer writtenFile.Close()

	err = nil
	writtenHash := sha256.New()
	if _, err := io.Copy(writtenHash, writtenFile); err != nil {

		res.Status = http.StatusInternalServerError
		res.Error = err
		res.Message = "Failed to hash the newly written file"

		s.writeResponse(w, res)
		return
	}
	res.Status = http.StatusOK
	res.Path = s.PathPrefix + strings.TrimPrefix(targetPath, s.DocumentRoot)
	res.Error = err
	res.Message = "File successfully uploaded"
	res.Hash = fmt.Sprintf("sha256:%x", writtenHash.Sum(nil))

	s.writeResponse(w, res)
}
