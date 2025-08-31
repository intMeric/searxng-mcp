package searxng

import (
	"context"
	"fmt"
)

// SearchOptions provides convenient ways to configure search requests
type SearchOptions struct {
	Language  string
	TimeRange TimeRange
	PageNo    int
}

// SimpleSearch performs a basic search with minimal configuration
func SimpleSearch(ctx context.Context, client Client, query string) (*SearchResponse, error) {
	req := SearchRequest{
		Query:    query,
		Category: []Category{CategoryGeneral},
	}
	
	return client.Search(ctx, req)
}

// SearchWithCategory performs a search in specific categories
func SearchWithCategory(ctx context.Context, client Client, query string, categories ...Category) (*SearchResponse, error) {
	req := SearchRequest{
		Query:    query,
		Category: categories,
	}
	
	return client.Search(ctx, req)
}

// SearchWithOptions performs a search with advanced options
func SearchWithOptions(ctx context.Context, client Client, query string, opts SearchOptions) (*SearchResponse, error) {
	req := SearchRequest{
		Query:     query,
		Language:  opts.Language,
		TimeRange: opts.TimeRange,
		PageNo:    opts.PageNo,
		Category:  []Category{CategoryGeneral},
	}
	
	return client.Search(ctx, req)
}

// ValidateCategory checks if a category string is valid
func ValidateCategory(category string) (Category, error) {
	switch Category(category) {
	case CategoryGeneral, CategoryImages, CategoryVideos, CategoryNews, 
		 CategoryMap, CategoryMusic, CategoryIT, CategoryScience, 
		 CategoryFiles, CategorySocialMedia:
		return Category(category), nil
	default:
		return "", fmt.Errorf("invalid category: %s", category)
	}
}

// ValidateTimeRange checks if a time range string is valid
func ValidateTimeRange(timeRange string) (TimeRange, error) {
	switch TimeRange(timeRange) {
	case TimeRangeDay, TimeRangeMonth, TimeRangeYear:
		return TimeRange(timeRange), nil
	default:
		return "", fmt.Errorf("invalid time range: %s", timeRange)
	}
}

// GetAllCategories returns all available categories
func GetAllCategories() []Category {
	return []Category{
		CategoryGeneral,
		CategoryImages,
		CategoryVideos,
		CategoryNews,
		CategoryMap,
		CategoryMusic,
		CategoryIT,
		CategoryScience,
		CategoryFiles,
		CategorySocialMedia,
	}
}

// GetAllTimeRanges returns all available time ranges
func GetAllTimeRanges() []TimeRange {
	return []TimeRange{
		TimeRangeDay,
		TimeRangeMonth,
		TimeRangeYear,
	}
}