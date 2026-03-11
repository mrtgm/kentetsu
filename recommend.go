package main

import (
	"math/rand/v2"
	"strings"
)

func PickRandom(recipes []Recipe) Recipe {
	return recipes[rand.IntN(len(recipes))]
}

func FilterByKeyword(recipes []Recipe, keyword string) []Recipe {
	keyword = strings.ToLower(keyword)
	var matched []Recipe
	for _, r := range recipes {
		if strings.Contains(strings.ToLower(r.Title), keyword) {
			matched = append(matched, r)
		}
	}
	return matched
}
