package xtream_codes_go

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiClient_GetLiveCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		wantCount      int
		wantError      bool
	}{
		{
			name:           "success empty",
			mockResponse:   `[]`,
			mockStatusCode: http.StatusOK,
			wantCount:      0,
			wantError:      false,
		},
		{
			name:           "success with categories",
			mockResponse:   `[{"category_id": "1", "category_name": "News"}, {"category_id": "2", "category_name": "Sports"}]`,
			mockStatusCode: http.StatusOK,
			wantCount:      2,
			wantError:      false,
		},
		{
			name:           "http error",
			mockResponse:   `Internal Server Error`,
			mockStatusCode: http.StatusInternalServerError,
			wantCount:      0,
			wantError:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("action") != "get_live_categories" {
					t.Errorf("Expected action get_live_categories, got %s", r.URL.Query().Get("action"))
				}
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client, err := NewApiClient(server.URL, "user", "pass")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			categories, err := client.GetLiveCategories(context.Background())
			if (err != nil) != tt.wantError {
				t.Fatalf("GetLiveCategories() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if len(categories) != tt.wantCount {
					t.Errorf("Expected %d categories, got %d", tt.wantCount, len(categories))
				}
			}
		})
	}
}

func TestApiClient_GetLiveStreams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		category       int
		mockResponse   string
		mockStatusCode int
		wantCount      int
		wantError      bool
		checkQuery     func(t *testing.T, r *http.Request)
	}{
		{
			name:           "success all categories",
			category:       -1,
			mockResponse:   `[{"stream_id": 1, "name": "Channel 1"}]`,
			mockStatusCode: http.StatusOK,
			wantCount:      1,
			wantError:      false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Has("category_id") {
					t.Errorf("Expected no category_id for category -1")
				}
			},
		},
		{
			name:           "success specific category",
			category:       5,
			mockResponse:   `[{"stream_id": 2, "name": "Channel 2"}]`,
			mockStatusCode: http.StatusOK,
			wantCount:      1,
			wantError:      false,
			checkQuery: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("category_id") != "5" {
					t.Errorf("Expected category_id 5, got %s", r.URL.Query().Get("category_id"))
				}
			},
		},
		{
			name:           "http error",
			category:       0,
			mockResponse:   `Not Found`,
			mockStatusCode: http.StatusNotFound,
			wantCount:      0,
			wantError:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("action") != "get_live_streams" {
					t.Errorf("Expected action get_live_streams, got %s", r.URL.Query().Get("action"))
				}
				if tt.checkQuery != nil {
					tt.checkQuery(t, r)
				}
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client, err := NewApiClient(server.URL, "user", "pass")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			streams, err := client.GetLiveStreams(context.Background(), tt.category)
			if (err != nil) != tt.wantError {
				t.Fatalf("GetLiveStreams() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if len(streams) != tt.wantCount {
					t.Errorf("Expected %d streams, got %d", tt.wantCount, len(streams))
				}
			}
		})
	}
}

func TestApiClient_GetLiveStreamUri(t *testing.T) {
	t.Parallel()

	client, err := NewApiClient("http://example.com", "user", "pass")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	client.setLoginInfo(&LoginInfo{
		UserInfo: &UserInfo{
			Username: "user",
			Password: "pass",
		},
		ServerInfo: &ServerInfo{
			ServerProtocol: "http",
			Url:            "example.com",
		},
	})

	tests := []struct {
		name      string
		streamID  int
		extension string
		want      string
	}{
		{
			name:      "with extension",
			streamID:  123,
			extension: "m3u8",
			want:      "http://example.com/live/user/pass/123.m3u8",
		},
		{
			name:      "without extension",
			streamID:  456,
			extension: "",
			want:      "http://example.com/live/user/pass/456",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := client.GetLiveStreamUri(tt.streamID, tt.extension)
			if got != tt.want {
				t.Errorf("GetLiveStreamUri() = %v, want %v", got, tt.want)
			}
		})
	}
}
