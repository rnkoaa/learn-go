package domain

import (
	"encoding/json"
	"fmt"
)

// Id -
type Id struct {
	Oid string `json:"$oid"`
}

// NewId -
func NewId() Id {
	return Id{}
}

// Timestamp -
type Timestamp struct {
	Date int64 `json:"$date"`
}

// NewTimestamp
func NewTimestamp() Timestamp {
	return Timestamp{}
}

// Recipe - json object
type Recipe struct {
	RecipeID      Id `json:"_id,omitempty"`
	Name          string
	Ingredients   string
	URL           string `json:"url"`
	Image         string
	CookTime      string
	Source        string
	RecipeYield   string `json:"recipeYield"`
	DatePublished string
	PrepTime      string
	Category      string `json:"recipeCategory,omitempty"`
	TotalTime     string `json:"totalTime"`
	Description   string
	Created       Timestamp `json:"ts"`
}

// ToString -
func (rp Recipe) ToString() string {
	result := ""
	result += "[Name: " + rp.Name + "]"
	return result
}

// NewRecipe -
func NewRecipe() Recipe {
	return Recipe{}
}

// FromJSON - converts recipe string to json
func FromJSON(recipeJSON string) Recipe {
	recipe := NewRecipe()
	err := json.Unmarshal([]byte(recipeJSON), &recipe)
	if err != nil {
		fmt.Printf("Failed to Process Recipe: %v\n", err)
	}
	return recipe
}
