package domain

import (
	"time"

	"github.com/araddon/dateparse"
)

// RecipeDTO - a data access object used to insert and remove
// data into db
type RecipeDTO struct {
	ID   string `gorm:"primary_key"`
	Name string `gorm:"type:varchar(256);not null"`
	// Name          string `gorm:"type:varchar(256);not null;unique_index"`
	Ingredients   string
	URL           string
	Image         string
	CookTime      string
	Source        string
	RecipeYield   string
	DatePublished *time.Time
	PrepTime      string
	Description   string
	TotalTime     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
	// Created       Timestamp
}

// FromRecipe - Converts a recipe to recipeDTO
func FromRecipe(recipe Recipe) RecipeDTO {
	centralTimeZone, _ := time.LoadLocation("America/Chicago")

	if &recipe == nil {
		return RecipeDTO{}
	}
	recipeDto := new(RecipeDTO)
	if &recipe.RecipeID != nil {
		recipeDto.ID = recipe.RecipeID.Oid
	}
	recipeDto.Name = recipe.Name
	recipeDto.Ingredients = recipe.Ingredients
	recipeDto.URL = recipe.URL
	recipeDto.Image = recipe.Image
	recipeDto.CookTime = recipe.CookTime
	recipeDto.Source = recipe.Source
	recipeDto.RecipeYield = recipe.RecipeYield
	recipeDto.TotalTime = recipe.TotalTime
	recipeDto.PrepTime = recipe.PrepTime
	recipeDto.Description = recipe.Description

	datePublished, err := dateparse.ParseIn(recipe.DatePublished, centralTimeZone)
	if err == nil {
		recipeDto.DatePublished = &datePublished
	}

	return *recipeDto
}

// TableName - set RecipeDTO's table name to be `recipe`
func (RecipeDTO) TableName() string {
	return "recipe"
}

// Type refers to the document type in bleve
func (recipeDTO *RecipeDTO) Type() string {
	return "recipe"
}

// NewRecipeDTO -
func NewRecipeDTO() *RecipeDTO {
	return &RecipeDTO{}
}
