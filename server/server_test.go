package server

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/endocode/go-simple-upload-server/logging"
)

func CreateTestServer() (Server, string) {
	tmpDir, err := ioutil.TempDir(".", "")

	if err != nil {
		log.Fatal(err)
	}

	content := []byte("<b>iaredigital</b>")
	tmpFile, err := ioutil.TempFile(tmpDir, "")
	if err != nil {
		log.Fatal(err)
	}

	if _, err = tmpFile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err = tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	err = nil
	tmpDirPath, err := filepath.Abs(tmpDir)
	if err != nil {
		log.Fatal(err)
	}
	logger := logging.Start(true)
	tmpFileBase := filepath.Base(tmpFile.Name())
	return Server{
		DocumentRoot:  tmpDirPath,
		PathPrefix:    "/wer/das/liest/ist/bloed",
		Logger:        logger,
		MaxUploadSize: 100,
	}, tmpFileBase
}

func TestNewServer(t *testing.T) {
	var root = "/moo"
	var max int64 = 1000

	s := New(root, max, "", nil)

	if s.MaxUploadSize != max {
		t.Errorf("MaxUploadSize: got %v want %v",
			s.MaxUploadSize, max)
	}
	if s.DocumentRoot != root {
		t.Errorf("DocumentRoot: got %v want %v",
			s.DocumentRoot, root)
	}

}
func TestNewServerPrefix(t *testing.T) {
	var goodPrefix = "/wuff"
	var badPrefix = "/wuff/"

	s := New("", 123, goodPrefix, nil)
	if s.PathPrefix != goodPrefix {
		t.Errorf("PathPrefix: got %v want %v",
			s.PathPrefix, goodPrefix)
	}

	q := New("", 123, badPrefix, nil)
	if q.PathPrefix != goodPrefix {
		t.Errorf("PathPrefix: got %v want %v",
			q.PathPrefix, goodPrefix)
	}
}
