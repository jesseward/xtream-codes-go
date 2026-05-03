package xtream_codes_go

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestApiTransport_GetRequestUri(t *testing.T) {
	t.Parallel()

	transport := &ApiTransport{}

	tests := []struct {
		name    string
		uriStr  string
		wantUri string
	}{
		{
			name:    "no username or password",
			uriStr:  "http://example.com/api?action=test&id=123",
			wantUri: "/api?action=test&id=123",
		},
		{
			name:    "with username and password > 5 chars",
			uriStr:  "http://example.com/api?username=mytestuser&password=mytestpassword",
			wantUri: "/api?password=myt******&username=myt******", // url.Values.Encode() sorts keys alphabetically
		},
		{
			name:    "with username and password <= 5 chars",
			uriStr:  "http://example.com/api?username=usr&password=pwd",
			wantUri: "/api?password=******&username=******",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			parsedURL, err := url.Parse(tt.uriStr)
			if err != nil {
				t.Fatalf("Failed to parse URL: %v", err)
			}
			got := transport.getRequestUri(*parsedURL)
			if got != tt.wantUri {
				t.Errorf("getRequestUri() = %v, want %v", got, tt.wantUri)
			}
		})
	}
}

func TestApiTransport_Dump(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	transport := &ApiTransport{
		dumper: &buf,
	}

	input := []byte("Line 1\r\nLine 2\r\nLine 3")
	expected := "> Line 1\n> Line 2\n> Line 3\n"

	transport.dump(input, "> ")

	if buf.String() != expected {
		t.Errorf("dump() got %q, want %q", buf.String(), expected)
	}
}

func TestApiTransport_RoundTrip(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	t.Run("with dumper and logger", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

		dumperBuf := &bytes.Buffer{}

		transport := &ApiTransport{
			inner:  http.DefaultTransport,
			logger: logger,
			dumper: dumperBuf,
		}

		req, err := http.NewRequest("GET", server.URL, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := transport.RoundTrip(req)
		if err != nil {
			t.Fatalf("RoundTrip() failed: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if string(body) != "OK" {
			t.Errorf("Expected body 'OK', got %q", body)
		}

		// Verify dumper output (should have request and response dumps)
		dumperOutput := dumperBuf.String()
		if !strings.Contains(dumperOutput, "> GET / HTTP/1.1") {
			t.Errorf("Expected dumper to contain request dump, got %q", dumperOutput)
		}
		if !strings.Contains(dumperOutput, "< HTTP/1.1 200 OK") {
			t.Errorf("Expected dumper to contain response dump, got %q", dumperOutput)
		}

		// Verify logger output
		logOutput := buf.String()
		if !strings.Contains(logOutput, "GET / HTTP/1.1 200") {
			t.Errorf("Expected logger to contain debug log, got %q", logOutput)
		}
	})

	t.Run("without logger or dumper", func(t *testing.T) {
		transport := &ApiTransport{
			inner: http.DefaultTransport,
		}

		req, err := http.NewRequest("GET", server.URL, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := transport.RoundTrip(req)
		if err != nil {
			t.Fatalf("RoundTrip() failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}
