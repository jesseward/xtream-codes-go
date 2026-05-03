package xtream_codes_go

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestBackdropPath_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		jsonData  string
		want      BackdropPath
		wantError bool
	}{
		{
			name:      "null",
			jsonData:  `null`,
			want:      nil,
			wantError: false,
		},
		{
			name:      "string",
			jsonData:  `"image1.jpg"`,
			want:      []string{"image1.jpg"},
			wantError: false,
		},
		{
			name:      "slice of strings",
			jsonData:  `["image1.jpg", "image2.jpg"]`,
			want:      []string{"image1.jpg", "image2.jpg"},
			wantError: false,
		},
		{
			name:      "slice with non-string elements",
			jsonData:  `["image1.jpg", 123]`,
			want:      nil,
			wantError: true,
		},
		{
			name:      "slice with mixed strings and nulls",
			jsonData:  `["image1.jpg", null]`,
			want:      []string{"image1.jpg"}, // null is skipped
			wantError: false,
		},
		{
			name:      "invalid type",
			jsonData:  `123`,
			want:      nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var b BackdropPath
			err := json.Unmarshal([]byte(tt.jsonData), &b)
			if (err != nil) != tt.wantError {
				t.Fatalf("UnmarshalJSON() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && !reflect.DeepEqual(b, tt.want) {
				t.Errorf("UnmarshalJSON() got = %v, want %v", b, tt.want)
			}
		})
	}
}

func TestApiClient_GetSeriesCategories(t *testing.T) {
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
			mockResponse:   `[{"category_id": "1", "category_name": "Action"}, {"category_id": "2", "category_name": "Comedy"}]`,
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
				if r.URL.Query().Get("action") != "get_series_categories" {
					t.Errorf("Expected action get_series_categories, got %s", r.URL.Query().Get("action"))
				}
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client, err := NewApiClient(server.URL, "user", "pass")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			categories, err := client.GetSeriesCategories(context.Background())
			if (err != nil) != tt.wantError {
				t.Fatalf("GetSeriesCategories() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if len(categories) != tt.wantCount {
					t.Errorf("Expected %d categories, got %d", tt.wantCount, len(categories))
				}
			}
		})
	}
}

func TestApiClient_GetSeries(t *testing.T) {
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
			mockResponse:   `[{"series_id": 1, "name": "Series 1"}]`,
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
			mockResponse:   `[{"series_id": 2, "name": "Series 2"}]`,
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
				if r.URL.Query().Get("action") != "get_series" {
					t.Errorf("Expected action get_series, got %s", r.URL.Query().Get("action"))
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

			series, err := client.GetSeries(context.Background(), tt.category)
			if (err != nil) != tt.wantError {
				t.Fatalf("GetSeries() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if len(series) != tt.wantCount {
					t.Errorf("Expected %d series, got %d", tt.wantCount, len(series))
				}
			}
		})
	}
}

func TestApiClient_GetSeriesInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		seriesID       int
		mockResponse   string
		mockStatusCode int
		wantError      bool
	}{
		{
			name:           "success",
			seriesID:       123,
			mockResponse:   `{"info": {"name": "Test Series"}, "seasons": []}`,
			mockStatusCode: http.StatusOK,
			wantError:      false,
		},
		{
			name:           "http error",
			seriesID:       123,
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
				if r.URL.Query().Get("action") != "get_series_info" {
					t.Errorf("Expected action get_series_info, got %s", r.URL.Query().Get("action"))
				}
				if r.URL.Query().Get("series_id") != "123" {
					t.Errorf("Expected series_id 123, got %s", r.URL.Query().Get("series_id"))
				}
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client, err := NewApiClient(server.URL, "user", "pass")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			info, err := client.GetSeriesInfo(context.Background(), tt.seriesID)
			if (err != nil) != tt.wantError {
				t.Fatalf("GetSeriesInfo() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if info == nil {
					t.Fatalf("Expected info to not be nil")
				}
				if info.Info.Name != "Test Series" {
					t.Errorf("Expected series name Test Series, got %s", info.Info.Name)
				}
			}
		})
	}
}

func TestApiClient_GetSeriesUri(t *testing.T) {
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
		seriesID  int
		extension string
		want      string
	}{
		{
			name:      "with extension",
			seriesID:  123,
			extension: "mp4",
			want:      "http://example.com/series/user/pass/123.mp4",
		},
		{
			name:      "without extension",
			seriesID:  456,
			extension: "",
			want:      "http://example.com/series/user/pass/456",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := client.GetSeriesUri(tt.seriesID, tt.extension)
			if got != tt.want {
				t.Errorf("GetSeriesUri() = %v, want %v", got, tt.want)
			}
		})
	}
}
