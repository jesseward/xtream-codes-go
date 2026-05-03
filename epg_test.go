package xtream_codes_go

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDate_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		jsonData  string
		want      time.Time
		wantError bool
	}{
		{
			name:      "valid date",
			jsonData:  `"2023-10-25 14:30:00"`,
			want:      time.Date(2023, 10, 25, 14, 30, 0, 0, time.UTC),
			wantError: false,
		},
		{
			name:      "invalid date format",
			jsonData:  `"2023/10/25 14:30:00"`,
			want:      time.Time{},
			wantError: true,
		},
		{
			name:      "invalid type",
			jsonData:  `123456`,
			want:      time.Time{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var d date
			err := json.Unmarshal([]byte(tt.jsonData), &d)
			if (err != nil) != tt.wantError {
				t.Fatalf("UnmarshalJSON() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && time.Time(d) != tt.want {
				t.Errorf("UnmarshalJSON() got = %v, want %v", time.Time(d), tt.want)
			}
		})
	}
}

func TestBase64string_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		jsonData  string
		want      base64string
		wantError bool
	}{
		{
			name:      "valid base64",
			jsonData:  `"SGVsbG8gV29ybGQ="`, // "Hello World"
			want:      "Hello World",
			wantError: false,
		},
		{
			name:      "invalid base64 uses fallback string",
			jsonData:  `"NotBase64!"`,
			want:      "NotBase64!",
			wantError: false,
		},
		{
			name:      "invalid type",
			jsonData:  `12345`,
			want:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var b base64string
			err := json.Unmarshal([]byte(tt.jsonData), &b)
			if (err != nil) != tt.wantError {
				t.Fatalf("UnmarshalJSON() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && b != tt.want {
				t.Errorf("UnmarshalJSON() got = %v, want %v", b, tt.want)
			}
		})
	}
}

func TestApiClient_GetSimpleDataTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		streamID       int
		mockResponse   string
		mockStatusCode int
		wantCount      int
		wantError      bool
	}{
		{
			name:           "success empty",
			streamID:       1,
			mockResponse:   `{"epg_listings": []}`,
			mockStatusCode: http.StatusOK,
			wantCount:      0,
			wantError:      false,
		},
		{
			name:           "success with items",
			streamID:       1,
			mockResponse:   `{"epg_listings": [{"epg_id": "123", "title": "SGVsbG8="}]}`,
			mockStatusCode: http.StatusOK,
			wantCount:      1,
			wantError:      false,
		},
		{
			name:           "http error",
			streamID:       1,
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
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client, err := NewApiClient(server.URL, "user", "pass")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			listings, err := client.GetSimpleDataTable(context.Background(), tt.streamID)
			if (err != nil) != tt.wantError {
				t.Fatalf("GetSimpleDataTable() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if listings == nil {
					t.Fatalf("Expected listings to not be nil")
				}
				if len(listings.EpgListings) != tt.wantCount {
					t.Errorf("Expected %d listings, got %d", tt.wantCount, len(listings.EpgListings))
				}
			}
		})
	}
}

func TestApiClient_GetShortEpg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		streamID       int
		mockResponse   string
		mockStatusCode int
		wantCount      int
		wantError      bool
	}{
		{
			name:           "success empty",
			streamID:       2,
			mockResponse:   `{"epg_listings": []}`,
			mockStatusCode: http.StatusOK,
			wantCount:      0,
			wantError:      false,
		},
		{
			name:           "success with items",
			streamID:       2,
			mockResponse:   `{"epg_listings": [{"epg_id": "456", "title": "VGVzdA=="}]}`,
			mockStatusCode: http.StatusOK,
			wantCount:      1,
			wantError:      false,
		},
		{
			name:           "http error",
			streamID:       2,
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
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client, err := NewApiClient(server.URL, "user", "pass")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			listings, err := client.GetShortEpg(context.Background(), tt.streamID)
			if (err != nil) != tt.wantError {
				t.Fatalf("GetShortEpg() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				if listings == nil {
					t.Fatalf("Expected listings to not be nil")
				}
				if len(listings.EpgListings) != tt.wantCount {
					t.Errorf("Expected %d listings, got %d", tt.wantCount, len(listings.EpgListings))
				}
			}
		})
	}
}
