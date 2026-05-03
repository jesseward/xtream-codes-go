package xtream_codes_go

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiClient_GetVodCategories(t *testing.T) {
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
			mockResponse:   `[{"category_id": "1", "category_name": "Movies"}, {"category_id": "2", "category_name": "Documentaries"}]`,
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
				if r.URL.Query().Get("action") != "get_vod_categories" {
					t.Errorf("Expected action get_vod_categories, got %s", r.URL.Query().Get("action"))
				}
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client, err := NewApiClient(server.URL, "user", "pass")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			categories, err := client.GetVodCategories(context.Background())
			if (err != nil) != tt.wantError {
				t.Fatalf("GetVodCategories() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if len(categories) != tt.wantCount {
					t.Errorf("Expected %d categories, got %d", tt.wantCount, len(categories))
				}
			}
		})
	}
}

func TestApiClient_GetVodStreams(t *testing.T) {
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
			mockResponse:   `[{"stream_id": 1, "name": "Movie 1"}]`,
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
			mockResponse:   `[{"stream_id": 2, "name": "Movie 2"}]`,
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
				if r.URL.Query().Get("action") != "get_vod_streams" {
					t.Errorf("Expected action get_vod_streams, got %s", r.URL.Query().Get("action"))
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

			streams, err := client.GetVodStreams(context.Background(), tt.category)
			if (err != nil) != tt.wantError {
				t.Fatalf("GetVodStreams() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if len(streams) != tt.wantCount {
					t.Errorf("Expected %d streams, got %d", tt.wantCount, len(streams))
				}
			}
		})
	}
}

func TestApiClient_GetVodInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		vodID          int
		mockResponse   string
		mockStatusCode int
		wantError      bool
	}{
		{
			name:           "success",
			vodID:          123,
			mockResponse:   `{"info": {"name": "Test Movie"}, "movie_data": {"stream_id": 123}}`,
			mockStatusCode: http.StatusOK,
			wantError:      false,
		},
		{
			name:           "http error",
			vodID:          123,
			mockResponse:   `Not Found`,
			mockStatusCode: http.StatusNotFound,
			wantError:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("action") != "get_vod_info" {
					t.Errorf("Expected action get_vod_info, got %s", r.URL.Query().Get("action"))
				}
				if r.URL.Query().Get("vod_id") != "123" {
					t.Errorf("Expected vod_id 123, got %s", r.URL.Query().Get("vod_id"))
				}
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client, err := NewApiClient(server.URL, "user", "pass")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			info, err := client.GetVodInfo(context.Background(), tt.vodID)
			if (err != nil) != tt.wantError {
				t.Fatalf("GetVodInfo() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if info == nil {
					t.Fatalf("Expected info to not be nil")
				}
				if info.Info.Name != "Test Movie" {
					t.Errorf("Expected movie name Test Movie, got %s", info.Info.Name)
				}
			}
		})
	}
}

func TestApiClient_GetVodUri(t *testing.T) {
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
			extension: "mp4",
			want:      "http://example.com/movie/user/pass/123.mp4",
		},
		{
			name:      "without extension",
			streamID:  456,
			extension: "",
			want:      "http://example.com/movie/user/pass/456",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := client.GetVodUri(tt.streamID, tt.extension)
			if got != tt.want {
				t.Errorf("GetVodUri() = %v, want %v", got, tt.want)
			}
		})
	}
}
