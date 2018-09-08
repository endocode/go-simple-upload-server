package server

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

// PostHandler handles all post requests
func (s Server) PostHandler(w http.ResponseWriter, r *http.Request) {
	var res response
	res.Path = r.URL.RequestURI()
	res.Method = r.Method

	srcFile, info, err := r.FormFile("file")
	if err != nil {

		res.Status = http.StatusInternalServerError
		res.Error = err
		res.Message = "Failed to acquire the uploaded content"

		s.writeResponse(w, res)
		return
	}
	defer srcFile.Close()

	if info.Filename == "" {

		res.Status = http.StatusInternalServerError
		res.Error = errors.New("Filename field empty")
		res.Message = "Filename field empty"

		s.writeResponse(w, res)
		return
	}

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

	body, err := ioutil.ReadAll(srcFile)
	if err != nil {

		res.Status = http.StatusInternalServerError
		res.Error = err
		res.Message = "Failed to read the uploaded file"

		s.writeResponse(w, res)
		return
	}

	targetPath := path.Join(s.DocumentRoot, info.Filename)
	targetFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {

		res.Status = http.StatusInternalServerError
		res.Error = err
		res.Message = "Failed to create file"

		s.writeResponse(w, res)
		return
	}
	defer targetFile.Close()

	err = nil
	if written, err := targetFile.Write(body); err != nil {

		res.Status = http.StatusInternalServerError
		res.Error = err
		res.Message = "Failed to write the content"

		s.writeResponse(w, res)
		return

	} else if int64(written) != size {

		res.Status = http.StatusInternalServerError
		res.Error = err
		res.Message = "Uploaded file size and written size differ"

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
	/*
		uploadedURL := strings.TrimPrefix(targetPath, s.DocumentRoot)
		if !strings.HasPrefix(uploadedURL, "/") {
			uploadedURL = "/" + uploadedURL
		}
		uploadedURL = s.PathPrefix + uploadedURL
	*/
	res.Status = http.StatusOK
	res.Path = s.PathPrefix + strings.TrimPrefix(targetPath, s.DocumentRoot)
	res.Error = err
	res.Message = "File successfully uploaded"
	res.Hash = fmt.Sprintf("sha256:%x", writtenHash.Sum(nil))

	s.writeResponse(w, res)

}
