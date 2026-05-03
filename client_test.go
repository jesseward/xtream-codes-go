package xtream_codes_go

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewApiClient_Options(t *testing.T) {
	t.Parallel()

	logger := slog.Default()
	httpClient := &http.Client{}
	dumper := &bytes.Buffer{}
	apiPath := "custom_api.php"

	client, err := NewApiClient(
		"http://example.com",
		"user",
		"pass",
		WithLogger(logger),
		WithHTTPClient(httpClient),
		WithDumper(dumper),
		WithAPIPath(apiPath),
	)

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client.logger != logger {
		t.Errorf("Expected logger %v, got %v", logger, client.logger)
	}

	if client.apiPath != apiPath {
		t.Errorf("Expected apiPath %s, got %s", apiPath, client.apiPath)
	}

	if client.client != httpClient {
		t.Errorf("Expected httpClient %v, got %v", httpClient, client.client)
	}

	transport, ok := client.client.Transport.(*ApiTransport)
	if !ok {
		t.Fatalf("Expected transport to be *ApiTransport")
	}

	if transport.dumper != dumper {
		t.Errorf("Expected dumper %v, got %v", dumper, transport.dumper)
	}
}

func TestNewApiClient_OptionValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opt     Option
		wantErr string
	}{
		{
			name:    "nil http client",
			opt:     WithHTTPClient(nil),
			wantErr: "http client cannot be nil",
		},
		{
			name:    "empty api path",
			opt:     WithAPIPath(""),
			wantErr: "apiPath cannot be empty",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewApiClient("http://example.com", "user", "pass", tt.opt)
			if err == nil {
				t.Fatalf("Expected error, got nil")
			}
			if !bytes.Contains([]byte(err.Error()), []byte(tt.wantErr)) {
				t.Errorf("Expected error to contain %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

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
