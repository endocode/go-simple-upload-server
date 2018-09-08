package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetHandler(t *testing.T) {
	s, tmpFile := CreateTestServer()
	defer os.RemoveAll(s.DocumentRoot)

	req, err := http.NewRequest("GET", s.PathPrefix+"/"+tmpFile, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `<b>iaredigital</b>`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
