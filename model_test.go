package xtream_codes_go

import (
	"encoding/json"
	"testing"
)

func TestNumericUnmarshal(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		want      numeric
		wantError bool
	}{
		{"valid int", `42`, 42, false},
		{"valid string int", `"42"`, 42, false},
		{"valid float", `42.5`, 42, false},
		{"invalid string", `"abc"`, 0, false},
		{"invalid type", `[]`, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n numeric
			err := json.Unmarshal([]byte(tt.jsonData), &n)
			if (err != nil) != tt.wantError {
				t.Fatalf("UnmarshalJSON() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && n != tt.want {
				t.Errorf("UnmarshalJSON() got = %v, want %v", n, tt.want)
			}
		})
	}
}

func TestBooleanUnmarshal(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		want      boolean
		wantError bool
	}{
		{"valid bool true", `true`, true, false},
		{"valid bool false", `false`, false, false},
		{"valid string true", `"true"`, true, false},
		{"valid string false", `"false"`, false, false},
		{"valid int 1", `1`, true, false},
		{"valid int 0", `0`, false, false},
		{"valid float >0", `1.5`, true, false},
		{"invalid string", `"abc"`, false, true},
		{"invalid type", `[]`, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b boolean
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

func TestVarcharUnmarshal(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		want      varchar
		wantError bool
	}{
		{"valid string", `"hello"`, "hello", false},
		{"valid int", `42`, "42", false},
		{"valid float", `42.5`, "42.5", false},
		{"invalid type", `[]`, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v varchar
			err := json.Unmarshal([]byte(tt.jsonData), &v)
			if (err != nil) != tt.wantError {
				t.Fatalf("UnmarshalJSON() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && v != tt.want {
				t.Errorf("UnmarshalJSON() got = %v, want %v", v, tt.want)
			}
		})
	}
}

func TestFloatUnmarshal(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		want      float
		wantError bool
	}{
		{"valid float", `42.5`, 42.5, false},
		{"valid string float", `"42.5"`, 42.5, false},
		{"valid int", `42`, 42.0, false},
		{"valid comma float", `"42,5"`, 42.5, false}, // Test comma replacement
		{"invalid string", `"abc"`, 0, false},
		{"invalid type", `[]`, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f float
			err := json.Unmarshal([]byte(tt.jsonData), &f)
			if (err != nil) != tt.wantError {
				t.Fatalf("UnmarshalJSON() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && f != tt.want {
				t.Errorf("UnmarshalJSON() got = %v, want %v", f, tt.want)
			}
		})
	}
}
