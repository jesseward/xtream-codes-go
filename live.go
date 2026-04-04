package xtream_codes_go

import (
	"context"
	"strconv"
)

type LiveStream struct {
	ModelBase
	ModelStream

	EpgChannelId      string  `json:"epg_channel_id"`
	TvArchive         boolean `json:"tv_archive"`
	TvArchiveDuration numeric `json:"tv_archive_duration"`
}

func (a *ApiClient) GetLiveCategories(ctx context.Context) ([]CategoryInterface, error) {
	return a.getCategories(ctx, CategoryTypeLive)
}

func (a *ApiClient) GetLiveStreams(ctx context.Context, category int) ([]*LiveStream, error) {
	var streams []*LiveStream
	var values map[string]string

	if category > -1 {
		values = map[string]string{"category_id": strconv.Itoa(category)}
	}

	if err := a.fetch(ctx, "get_live_streams", values, playerApi, &streams); err != nil {
		return nil, err
	}

	return streams, nil
}

func (a *ApiClient) GetLiveStreamUri(streamId int, extension string) string {
	return a.streamUrl("live", streamId, extension)
}
