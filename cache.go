package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	cacheDir      = "kentetsu"
	cacheFile     = "recipes.json"
	remoteCacheURL = "https://raw.githubusercontent.com/mrtgm/kentetsu/main/recipes.json"
)

func cacheFilePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, cacheDir, cacheFile), nil
}

func LoadCache() (*RecipeData, error) {
	path, err := cacheFilePath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data RecipeData
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func SaveCache(data *RecipeData) error {
	path, err := cacheFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0o644)
}

func IsCacheValid(data *RecipeData, ttl time.Duration) bool {
	t, err := time.Parse(time.RFC3339, data.UpdatedAt)
	if err != nil {
		return false
	}
	return time.Since(t) < ttl
}

func UpdateFromRemote() (*RecipeData, error) {
	resp, err := http.Get(remoteCacheURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch remote recipes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var data RecipeData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to decode recipes: %w", err)
	}

	if err := SaveCache(&data); err != nil {
		return nil, fmt.Errorf("failed to save cache: %w", err)
	}

	return &data, nil
}
