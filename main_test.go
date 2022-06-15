package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestGenerateApiKey(t *testing.T) {
	key := generateApiKey(32)
	if len(key) != 32 {
		t.Error("GenerateApiKey does not return the correct length")
	}
}

func TestReadConfig(t *testing.T) {
	var err error
	// Test invalid config
	temp_config, _ := ioutil.TempFile("", "")
	defer os.Remove(temp_config.Name())
	temp_config.Write([]byte("awrgaergerherherherh"))
	c, err = ReadConfig(temp_config.Name())
	if err == nil {
		t.Error("Config should be invalid")
	}
	// Test valid config
	cfg := Config{AllowedHosts: []string{}, ApiKey: "12345", Port: 7000}
	b, _ := json.Marshal(cfg)
	temp_config.WriteAt(b, 0)
	c, err = ReadConfig(temp_config.Name())
	if err != nil {
		t.Error("Config should be valid")
	}
	// Test Config does not exist
	tdir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tdir)
	c, err = ReadConfig(path.Join(tdir, "config.json"))
	if err != nil {
		t.Error("Config should be created")
	}
}

func TestShutdownHandler(t *testing.T) {
	var err error
	temp_config, _ := ioutil.TempFile("", "")
	defer os.Remove(temp_config.Name())
	cfg := Config{AllowedHosts: []string{"192.168.1.22"}, ApiKey: "12345", Port: 7000}
	b, _ := json.Marshal(cfg)
	temp_config.Write(b)

	c, err = ReadConfig(temp_config.Name())
	if err != nil {
		t.Error(err)
	}
	// Test invalid Method
	mux := http.NewServeMux()
	mux.HandleFunc("/", shutdownHandler)
	writer := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	mux.ServeHTTP(writer, req)
	if writer.Code != http.StatusMethodNotAllowed {
		t.Errorf("Handler should return MethodNotAllowed")
	}
	// Test unauthorized
	req.Method = "POST"
	writer = httptest.NewRecorder()
	mux.ServeHTTP(writer, req)
	if writer.Code != http.StatusUnauthorized {
		t.Errorf("Handler should return Unauthorized")
	}
	// Test AllowedHosts
	req.Header.Set("Authorization", "Bearer 12345")
	writer = httptest.NewRecorder()
	mux.ServeHTTP(writer, req)
	if writer.Code != http.StatusUnauthorized {
		t.Error("AllowedHosts not respected")
	}
	req.RemoteAddr = "192.168.1.22:51234"

	writer = httptest.NewRecorder()
	mux.ServeHTTP(writer, req)
	if writer.Code != http.StatusOK && writer.Code != http.StatusInternalServerError {
		t.Errorf("Handler returned %d", writer.Code)
	}

	// Test ok
	writer = httptest.NewRecorder()
	mux.ServeHTTP(writer, req)
	if writer.Code != http.StatusOK && writer.Code != http.StatusInternalServerError {
		t.Errorf("Handler returned %d", writer.Code)
	}
}
