package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

//go:embed recipes.json
var embeddedRecipes []byte

const cacheTTL = 7 * 24 * time.Hour

func loadRecipes(forceUpdate bool) ([]Recipe, error) {
	// Try local cache first (from --update)
	if !forceUpdate {
		if data, err := LoadCache(); err == nil && IsCacheValid(data, cacheTTL) {
			return data.Recipes, nil
		}
	}

	// If --update, fetch from GitHub
	if forceUpdate {
		fmt.Println("⏳ レシピデータを更新中...")
		data, err := UpdateFromRemote()
		if err != nil {
			return nil, err
		}
		fmt.Printf("✅ %d 件のレシピを取得しました\n", len(data.Recipes))
		return data.Recipes, nil
	}

	// Fall back to embedded data
	var data RecipeData
	if err := json.Unmarshal(embeddedRecipes, &data); err != nil {
		return nil, fmt.Errorf("failed to load embedded recipes: %w", err)
	}
	return data.Recipes, nil
}

func openURL(url string) error {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "start"
	default:
		return fmt.Errorf("unsupported platform")
	}
	return exec.Command(cmd, url).Start()
}

func main() {
	search := flag.String("search", "", "キーワードで絞り込み")
	update := flag.Bool("update", false, "レシピデータを最新に更新")
	open := flag.Bool("open", false, "ブラウザで動画を開く")
	flag.Parse()

	recipes, err := loadRecipes(*update)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	if len(recipes) == 0 {
		fmt.Fprintln(os.Stderr, "レシピが見つかりません")
		os.Exit(1)
	}

	if *search != "" {
		recipes = FilterByKeyword(recipes, *search)
		if len(recipes) == 0 {
			fmt.Fprintf(os.Stderr, "「%s」に一致するレシピが見つかりません\n", *search)
			os.Exit(1)
		}
	}

	if *update && *search == "" {
		return
	}

	recipe := PickRandom(recipes)
	DisplayRecipe(recipe)

	if *open {
		if err := openURL(recipe.URL); err != nil {
			fmt.Fprintf(os.Stderr, "ブラウザを開けませんでした: %v\n", err)
		}
	}
}
