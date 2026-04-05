package xtream_codes_go

import (
	"context"
	"sort"
)

type CategoryType string

const (
	CategoryTypeLive   CategoryType = "live"
	CategoryTypeVod    CategoryType = "vod"
	CategoryTypeSeries CategoryType = "series"
)

type Category struct {
	Id       int    `json:"category_id,string"`
	Name     string `json:"category_name"`
	ParentId int    `json:"parent_id"`
}

func (a *ApiClient) getCategories(ctx context.Context, categoryType CategoryType) ([]*Category, error) {
	var categories []*Category

	if err := a.fetch(ctx, "get_"+string(categoryType)+"_categories", nil, a.apiPath, &categories); err != nil {
		return nil, err
	}

	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Name < categories[j].Name
	})

	return categories, nil
}
