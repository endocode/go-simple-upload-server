package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	s, _ := CreateTestServer()
	defer os.RemoveAll(s.DocumentRoot)

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.HealthCheckHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want application/json",
			contentType)
	}

	r := response{}
	if err := json.Unmarshal(rr.Body.Bytes(), &r); err != nil {
		t.Errorf("Cannot unmarshal json: %v",
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
}
