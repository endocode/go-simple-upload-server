package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPostHandler(t *testing.T) {
	s, _ := CreateTestServer()
	content := []byte("content")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "uploadme")

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, bytes.NewReader(content))
	writer.Close()

	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.PostHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	r := response{}
	if err := json.Unmarshal(rr.Body.Bytes(), &r); err != nil {
		t.Errorf("Cannot unmarshal %v",
			err)
	}

	if r.Status != 200 {
		t.Errorf("handler returned not OK: got %v want true",
			r.Status)
	}

	if r.Error != nil {
		t.Errorf("handler returned an error: got %v want nil",
			r.Error)
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, bytes.NewReader(content)); err != nil {
		t.Fatal(err)
	}
	sum := fmt.Sprintf("sha256:%x", hash.Sum(nil))
	if r.Hash != sum {
		t.Errorf("Hash mismatch: got %v want %v",
			r.Hash, sum)
	}
	m := "File successfully uploaded"
	if r.Message != m {
		t.Errorf("Message differs: got %v want %v",
			r.Message, m)
	}
	defer os.RemoveAll(s.DocumentRoot)
}
