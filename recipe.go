package main

import "time"

type Recipe struct {
	Title       string    `json:"title"`
	VideoID     string    `json:"videoId"`
	URL         string    `json:"url"`
	Thumbnail   string    `json:"thumbnail"`
	PublishedAt time.Time `json:"publishedAt"`
}

type RecipeData struct {
	UpdatedAt string   `json:"updatedAt"`
	Recipes   []Recipe `json:"recipes"`
}
