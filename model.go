package xtream_codes_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type ModelBase struct {
	Num         int    `json:"num"`
	Name        string `json:"name"`
	CategoryId  int    `json:"category_id,string"`
	CategoryIds []int  `json:"category_ids"`
}

type ModelStream struct {
	StreamType   string  `json:"stream_type"`
	StreamId     int     `json:"stream_id"`
	StreamIcon   string  `json:"stream_icon"`
	Added        int     `json:"added,string"`
	IsAdult      boolean `json:"is_adult"`
	CustomSid    string  `json:"custom_sid"`
	DirectSource string  `json:"direct_source"`
}

type ModelVideo struct {
	Rating       float   `json:"rating"`
	Rating5Based float   `json:"rating_5based"`
	Tmdb         varchar `json:"tmdb"`
}

type numeric int

func (n *numeric) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	s := string(data)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		slog.Warn("unexpected value for numeric, defaulting to 0", "value", string(data), "error", err)
		*n = 0
		return nil
	}

	*n = numeric(int(f))
	return nil
}

type boolean bool

func (b *boolean) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	s := string(data)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	if s == "true" {
		*b = true
		return nil
	}
	if s == "false" {
		*b = false
		return nil
	}

	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		*b = boolean(f > 0)
		return nil
	}

	return fmt.Errorf("unexpected value %q for boolean", string(data))
}

type varchar string

func (v *varchar) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*v = varchar(s)
		return nil
	}

	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		if f == float64(int64(f)) {
			*v = varchar(strconv.FormatInt(int64(f), 10))
		} else {
			*v = varchar(strconv.FormatFloat(f, 'f', -1, 64))
		}
		return nil
	}

	return fmt.Errorf("unexpected value %q for varchar", string(data))
}

type float float32

func (f *float) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	data = bytes.Replace(data, []byte{','}, []byte{'.'}, -1)

	s := string(data)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil
	}

	val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		slog.Warn("unexpected value for float, defaulting to 0", "value", string(data), "error", err)
		*f = 0
		return nil
	}

	*f = float(val)
	return nil
}
