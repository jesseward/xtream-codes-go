package xtream_codes_go

import (
	"context"
	"strconv"
)

type VodStream struct {
	ModelBase
	ModelStream
	ModelVideo

	Trailer            string `json:"trailer"`
	ContainerExtension string `json:"container_extension"`
}

type MovieData struct {
	StreamId           int    `json:"stream_id"`
	Name               string `json:"name"`
	Added              string `json:"added"`
	CategoryId         string `json:"category_id"`
	CategoryIds        []int  `json:"category_ids"`
	ContainerExtension string `json:"container_extension"`
	CustomSid          string `json:"custom_sid"`
	DirectSource       string `json:"direct_source"`
}

type Movie struct {
	TmdbUrl              string   `json:"tmdb_url"`
	TmdbId               string   `json:"tmdb_id"`
	Name                 string   `json:"name"`
	OName                string   `json:"o_name"`
	CoverBig             string   `json:"cover_big"`
	MovieImage           string   `json:"movie_image"`
	Releasedate          string   `json:"releasedate"`
	EpisodeRunTime       string   `json:"episode_run_time"`
	YoutubeTrailer       string   `json:"youtube_trailer"`
	Director             string   `json:"director"`
	Actors               string   `json:"actors"`
	Cast                 string   `json:"cast"`
	Description          string   `json:"description"`
	Plot                 string   `json:"plot"`
	Age                  string   `json:"age"`
	MpaaRating           string   `json:"mpaa_rating"`
	RatingCountKinopoisk int      `json:"rating_count_kinopoisk"`
	Country              string   `json:"country"`
	Genre                string   `json:"genre"`
	BackdropPath         []string `json:"backdrop_path"`
	DurationSecs         int      `json:"duration_secs"`
	Duration             string   `json:"duration"`
	Video                *Video   `json:"video"`
	Audio                *Audio   `json:"audio"`
	Bitrate              int      `json:"bitrate"`
	Rating               string   `json:"rating"`
}

type VodInfo struct {
	Info      *Movie     `json:"info"`
	MovieData *MovieData `json:"movie_data"`
}

func (a *ApiClient) GetVodStreams(ctx context.Context, category int) ([]*VodStream, error) {
	var streams []*VodStream
	var values map[string]string

	if category > -1 {
		values = map[string]string{"category_id": strconv.Itoa(category)}
	}

	if err := a.fetch(ctx, "get_vod_streams", values, a.apiPath, &streams); err != nil {
		return nil, err
	}

	return streams, nil
}

func (a *ApiClient) GetVodInfo(ctx context.Context, id int) (*VodInfo, error) {
	var seriesInfo *VodInfo

	var values = map[string]string{
		"vod_id": strconv.Itoa(id),
	}

	if err := a.fetch(ctx, "get_vod_info", values, a.apiPath, &seriesInfo); err != nil {
		return nil, err
	}

	return seriesInfo, nil
}

func (a *ApiClient) GetVodCategories(ctx context.Context) ([]*Category, error) {
	return a.getCategories(ctx, CategoryTypeVod)
}

func (a *ApiClient) GetVodUri(id int, extension string) string {
	return a.streamUrl("movie", id, extension)
}
