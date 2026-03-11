// Command fetch retrieves all recipes from Koh Kentetsu's YouTube channel
// and writes them to recipes.json. Used by CI and for initial data setup.
//
// Usage: YOUTUBE_API_KEY=xxx go run ./cmd/fetch
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	uploadsPlaylistID = "UU3p5OTQsMEnmZktWUkw_Y0A"
	apiBaseURL        = "https://www.googleapis.com/youtube/v3/playlistItems"
)

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

type playlistResponse struct {
	NextPageToken string `json:"nextPageToken"`
	Items         []struct {
		Snippet struct {
			Title       string    `json:"title"`
			PublishedAt time.Time `json:"publishedAt"`
			ResourceID  struct {
				VideoID string `json:"videoId"`
			} `json:"resourceId"`
			Thumbnails struct {
				High struct {
					URL string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
}

func fetchAll(apiKey string) ([]Recipe, error) {
	var recipes []Recipe
	pageToken := ""

	for {
		url := fmt.Sprintf("%s?part=snippet&playlistId=%s&maxResults=50&key=%s",
			apiBaseURL, uploadsPlaylistID, apiKey)
		if pageToken != "" {
			url += "&pageToken=" + pageToken
		}

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("API request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
		}

		var result playlistResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		for _, item := range result.Items {
			s := item.Snippet
			recipes = append(recipes, Recipe{
				Title:       s.Title,
				VideoID:     s.ResourceID.VideoID,
				URL:         "https://youtu.be/" + s.ResourceID.VideoID,
				Thumbnail:   s.Thumbnails.High.URL,
				PublishedAt: s.PublishedAt,
			})
		}

		if result.NextPageToken == "" {
			break
		}
		pageToken = result.NextPageToken
	}

	return recipes, nil
}

func main() {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "YOUTUBE_API_KEY environment variable is required")
		os.Exit(1)
	}

	fmt.Println("⏳ YouTube API からレシピを取得中...")
	recipes, err := fetchAll(apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	data := RecipeData{
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Recipes:   recipes,
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("recipes.json", b, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ %d 件のレシピを recipes.json に保存しました\n", len(recipes))
}
