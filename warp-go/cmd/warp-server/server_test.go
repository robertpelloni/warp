package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCorsMiddleware(t *testing.T) {
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("OPTIONS", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if origin := rr.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("handler returned wrong origin: got %v want %v", origin, "*")
	}
}

func TestAgentEndpoint(t *testing.T) {
	reqBody := strings.NewReader(`{"message": "Hello Warp"}`)
	req, err := http.NewRequest("POST", "/agent", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAgent)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"reply":"I am Warp Go Agent. I received your message: Hello Warp"}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestRpcEndpoint(t *testing.T) {
	reqBody := strings.NewReader(`{"jsonrpc": "2.0", "method": "testMethod", "params": {}, "id": 1}`)
	req, err := http.NewRequest("POST", "/rpc", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRPC)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"jsonrpc":"2.0","result":"Method testMethod executed successfully","id":1}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
