package xtream_codes_go

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	var x interface{}
	var v int

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	switch y := x.(type) {
	case string:
		var err error
		v, err = strconv.Atoi(y)
		if err != nil {
			return err
		}
	case int8:
		v = int(y)
	case int16:
		v = int(y)
	case int32:
		v = int(y)
	case int64:
		v = int(y)
	case int:
		v = y
	case float32:
		v = int(y)
	case float64:
		v = int(y)
	default:
		return fmt.Errorf("unexpected type %T for numeric", y)
	}

	*n = numeric(v)

	return nil
}

type boolean bool

func (b *boolean) UnmarshalJSON(data []byte) error {
	var x interface{}
	var v bool

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	switch y := x.(type) {
	case string:
		var err error
		v, err = strconv.ParseBool(y)
		if err != nil {
			return err
		}
	case int8:
		v = y != 0
	case int16:
		v = y != 0
	case int32:
		v = y != 0
	case int64:
		v = y != 0
	case int:
		v = y != 0
	case float32:
		v = y > 0
	case float64:
		v = y > 0
	case bool:
		v = y
	default:
		return fmt.Errorf("unexpected type %T for boolean", y)
	}

	*b = boolean(v)

	return nil
}

type varchar string

func (v *varchar) UnmarshalJSON(data []byte) error {
	var x interface{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	switch y := x.(type) {
	case string:
		*v = varchar(y)
	case int8:
		*v = varchar(strconv.Itoa(int(y)))
	case int16:
		*v = varchar(strconv.Itoa(int(y)))
	case int32:
		*v = varchar(strconv.Itoa(int(y)))
	case int64:
		*v = varchar(strconv.Itoa(int(y)))
	case int:
		*v = varchar(strconv.Itoa(y))
	case float32:
		*v = varchar(strconv.FormatFloat(float64(y), 'f', -1, 32))
	case float64:
		*v = varchar(strconv.FormatFloat(y, 'f', -1, 64))
	default:
		return fmt.Errorf("unexpected type %T for varchar", y)
	}

	return nil
}

type float float32

func (f *float) UnmarshalJSON(data []byte) error {
	var x interface{}

	data = bytes.Replace(data, []byte{','}, []byte{'.'}, -1)

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	switch y := x.(type) {
	case string:
		if len(y) > 0 {
			x, err := strconv.ParseFloat(strings.TrimSpace(y), 32)
			if err != nil {
				return err
			}
			*f = float(x)
		}
	case int8:
		*f = float(y)
	case int16:
		*f = float(y)
	case int32:
		*f = float(y)
	case int64:
		*f = float(y)
	case int:
		*f = float(y)
	case float32:
		*f = float(y)
	case float64:
		*f = float(y)
	default:
		return fmt.Errorf("unexpected type %T for float", y)
	}

	return nil
}
