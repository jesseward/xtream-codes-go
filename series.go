package xtream_codes_go

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

type BackdropPath []string

func (b *BackdropPath) UnmarshalJSON(data []byte) error {
	var x interface{}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	switch y := x.(type) {
	case string:
		*b = []string{y}
	case []string:
		*b = y
	case []interface{}:
		*b = make([]string, 0)
		for _, v := range y {
			if s, ok := v.(string); ok {
				*b = append(*b, s)
			} else if v != nil {
				return fmt.Errorf("expecting string got %#v", v)
			}
		}
	default:
		return fmt.Errorf("invalid BackdropPath type %#v", x)
	}
	return nil
}

type Series struct {
	ModelBase
	ModelVideo
	SeriesId       int          `json:"series_id"`
	Cover          string       `json:"cover"`
	Plot           string       `json:"plot"`
	Cast           string       `json:"cast"`
	Director       string       `json:"director"`
	Genre          string       `json:"genre"`
	ReleaseDate    string       `json:"releaseDate"`
	LastModified   string       `json:"last_modified"`
	BackdropPath   BackdropPath `json:"backdrop_path"`
	YoutubeTrailer string       `json:"youtube_trailer"`
	EpisodeRunTime int          `json:"episode_run_time,string"`
}

type Season struct {
	Name         string  `json:"name"`
	EpisodeCount numeric `json:"episode_count"`
	Overview     string  `json:"overview"`
	AirDate      string  `json:"air_date"`
	Cover        string  `json:"cover"`
	CoverTmdb    string  `json:"cover_tmdb"`
	SeasonNumber numeric `json:"season_number"`
	CoverBig     string  `json:"cover_big"`
	ReleaseDate  string  `json:"releaseDate"`
	Duration     string  `json:"duration"`
}

type EpisodeInfo struct {
	AirDate      string  `json:"air_date"`
	Crew         string  `json:"crew"`
	Rating       float   `json:"rating"`
	Id           numeric `json:"id"`
	MovieImage   string  `json:"movie_image"`
	DurationSecs int     `json:"duration_secs"`
	Duration     string  `json:"duration"`
	Bitrate      int     `json:"bitrate"`
	Video        *Video  `json:"video"`
	Audio        *Video  `json:"audio"`
}

type Episode struct {
	Id                 int          `json:"id,string"`
	EpisodeNum         int          `json:"episode_num"`
	Title              string       `json:"title"`
	ContainerExtension string       `json:"container_extension"`
	Info               *EpisodeInfo `json:"info"`
}

type SeriesInfo struct {
	Seasons  []*Season          `json:"seasons"`
	Info     Series             `json:"info"`
	Episodes map[int][]*Episode `json:"episodes"` // indexed by season
}

func (a *ApiClient) GetSeries(ctx context.Context, category int) ([]*Series, error) {
	var series []*Series
	var values map[string]string

	if category > -1 {
		values = map[string]string{"category_id": strconv.Itoa(category)}
	}

	if err := a.fetch(ctx, "get_series", values, a.apiPath, &series); err != nil {
		return nil, err
	}

	return series, nil
}

func (a *ApiClient) GetSeriesInfo(ctx context.Context, id int) (*SeriesInfo, error) {
	var seriesInfo *SeriesInfo

	if err := a.fetch(ctx, "get_series_info", map[string]string{"series_id": strconv.Itoa(id)}, a.apiPath, &seriesInfo); err != nil {
		return nil, err
	}

	return seriesInfo, nil
}

func (a *ApiClient) GetSeriesCategories(ctx context.Context) ([]CategoryInterface, error) {
	return a.getCategories(ctx, CategoryTypeSeries)
}

// GetSeriesUri build serie url
//
// var serie *Serie
// ...
// client.GetSeriesUri(serie.Id, serie.ContainerExtension)
func (a *ApiClient) GetSeriesUri(id int, extension string) string {
	return a.streamUrl("series", id, extension)
}
