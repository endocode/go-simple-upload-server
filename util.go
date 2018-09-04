package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type response struct {
	OK bool `json:"ok"`
}

type uploadedResponse struct {
	response
	Path string `json:"path"`
	Hash string `json:"hash"`
}

func newUploadedResponse(path, hash string) uploadedResponse {
	return uploadedResponse{response: response{OK: true}, Path: path, Hash: hash}
}

type errorResponse struct {
	response
	Message string `json:"error"`
}

func newErrorResponse(err error) errorResponse {
	return errorResponse{response: response{OK: false}, Message: err.Error()}
}

func writeError(w http.ResponseWriter, err error) (int, error) {
	body := newErrorResponse(err)
	b, e := json.Marshal(body)
	// if an error is occured on marshaling, write empty value as response.
	if e != nil {
		return w.Write([]byte{})
	}
	return w.Write(b)
}

func writeSuccess(w http.ResponseWriter, path string, hash string) (int, error) {
	body := newUploadedResponse(path, hash)
	b, e := json.Marshal(body)
	// if an error is occured on marshaling, write empty value as response.
	if e != nil {
		return w.Write([]byte{})
	}
	return w.Write(b)
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
