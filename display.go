package main

import "fmt"

const (
	bold  = "\033[1m"
	cyan  = "\033[36m"
	dim   = "\033[2m"
	reset = "\033[0m"
)

func DisplayRecipe(r Recipe) {
	fmt.Println()
	fmt.Printf("%s🍳 今日のコウケンテツレシピ%s\n", bold, reset)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Printf("  %s%s%s\n", bold, r.Title, reset)
	fmt.Println()
	fmt.Printf("  %s📅 %s%s\n", dim, r.PublishedAt.Format("2006-01-02"), reset)
	fmt.Printf("  %s🔗 %s%s\n", cyan, r.URL, reset)
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
}
