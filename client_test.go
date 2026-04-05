package xtream_codes_go

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiClient_Flow(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		action := query.Get("action")

		if query.Get("username") != "testuser" {
			t.Errorf("Expected username testuser, got '%s'", query.Get("username"))
		}
		if query.Get("password") != "testpass" {
			t.Errorf("Expected password testpass, got '%s'", query.Get("password"))
		}

		switch action {
		case "":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"user_info":{"username":"testuser","password":"testpass","auth":1,"status":"Active","message":"Login successful","exp_date":"1700000000"},"server_info":{"timezone":"UTC","server_protocol":"http","url":"` + r.Host + `"}}`))
		case "get_series":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[{"num":1,"name":"Test Series","category_id":"1","series_id":123,"backdrop_path":"url1"}]`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, err := NewApiClient(server.URL, "testuser", "testpass")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if err := client.Connect(context.Background()); err != nil {
		t.Fatalf("Failed to connect client: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}

	series, err := client.GetSeries(context.Background(), 0)
	if err != nil {
		t.Fatalf("Failed to get series: %v", err)
	}

	if len(series) != 1 {
		t.Errorf("Expected 1 series, got %d", len(series))
	}

	if series[0].Name != "Test Series" {
		t.Errorf("Expected series name 'Test Series', got '%s'", series[0].Name)
	}
}

func TestNewApiClient_LoginFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user_info":{"auth":0,"status":"Expired","message":"Login failed"},"server_info":{"timezone":"UTC"}}`))
	}))
	defer server.Close()

	client, err := NewApiClient(server.URL, "testuser", "testpass")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if err := client.Connect(context.Background()); err == nil {
		t.Fatal("Expected error on failed login, got nil")
	}
}
