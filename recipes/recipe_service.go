package main

import (
	"learn-go/recipes/domain"
	"learn-go/recipes/repo"

	"github.com/jinzhu/gorm"
)

// SqliteRecipeService -
type SqliteRecipeService struct {
	repository repo.RecipeRepository
}

// RecipeService -
type RecipeService interface {
	FindRecipe(id string) (*domain.RecipeDTO, error)
	FindAllRecipes() ([]*domain.RecipeDTO, error)
	FindAllByRecipeName(name string) ([]domain.RecipeDTO, error)
	Count() (int64, error)
}

// NewSqliteRecipeService -
func NewSqliteRecipeService(db *gorm.DB) RecipeService {
	return &SqliteRecipeService{
		repository: repo.NewSqliteRecipeRepository(db),
	}
}

// FindRecipe -
func (impl *SqliteRecipeService) FindRecipe(id string) (*domain.RecipeDTO, error) {
	recipeDto, _ := impl.repository.Find(id)
	return recipeDto, nil
}

// FindAllRecipes -
func (impl *SqliteRecipeService) FindAllRecipes() ([]*domain.RecipeDTO, error) {
	recipeDtos := make([]*domain.RecipeDTO, 0)
	return recipeDtos, nil
}

// FindAllByRecipeName -
func (impl *SqliteRecipeService) FindAllByRecipeName(name string) ([]domain.RecipeDTO, error) {
	recipeDtos, err := impl.repository.FindByName(name)
	return recipeDtos, err
}

// FindAllByRecipeName -
func (impl *SqliteRecipeService) Count() (int64, error) {
	count, err := impl.repository.Count()
	return count, err
}
