package xtream_codes_go

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"
)

type date time.Time

func (s *date) UnmarshalJSON(data []byte) error {
	var x string

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	val, err := time.Parse("2006-01-02 15:04:05", x)

	if err != nil {
		return err
	}

	*s = date(val)

	return nil
}

type base64string string

func (s *base64string) UnmarshalJSON(data []byte) error {
	var x string

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	out, err := base64.StdEncoding.DecodeString(x)

	if err == nil {
		*s = base64string(out)
	} else {
		*s = base64string(data)
	}

	return nil
}

type EpgListing struct {
	EpgId          int          `json:"epg_id,string"`
	Title          base64string `json:"title"`
	Lang           string       `json:"lang"`
	Start          date         `json:"start"`
	End            date         `json:"end"`
	Description    base64string `json:"description"`
	ChannelId      string       `json:"channel_id"`
	StartTimestamp int          `json:"start_timestamp,string"`
	StopTimestamp  int          `json:"stop_timestamp,string"`
	//short
	StreamId int `json:"stream_id,string"`
	// long
	Id         int     `json:"id,string"`
	NowPlaying boolean `json:"now_playing"`
	HasArchive boolean `json:"has_archive"`
}
type EpgListening struct {
	EpgListings []*EpgListing `json:"epg_listings"`
}

func (a *ApiClient) GetSimpleDataTable(ctx context.Context, streamId int) (*EpgListening, error) {
	var egpInfo *EpgListening
	var values = map[string]string{"stream_id": strconv.Itoa(streamId)}

	if err := a.fetch(a.context(ctx, "get_simple_data_table", values), playerApi, &egpInfo); err != nil {
		return nil, err
	}

	return egpInfo, nil
}

func (a *ApiClient) GetShortEpg(ctx context.Context, streamId int) (*EpgListening, error) {

	var egpInfo *EpgListening
	var values = map[string]string{"stream_id": strconv.Itoa(streamId)}

	if err := a.fetch(a.context(ctx, "get_short_epg", values), playerApi, &egpInfo); err != nil {
		return nil, err
	}

	return egpInfo, nil

}
