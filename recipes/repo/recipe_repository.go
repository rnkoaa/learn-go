package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"learn-go/recipes/domain"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

// DbTemplate -
type DbTemplate struct {
	db *gorm.DB
}

const (
	INSERT_STMT = `INSERT INTO recipe (id, name, ingredients, url, image, 
		cook_time, source, recipe_yield, date_published, prep_time, 
		description,total_time, created_at, updated_at, deleted_at) 
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

// RecipeRepository -
type RecipeRepository interface {
	CreateTableIFNotExists() error
	Find(id string) (*domain.RecipeDTO, error)
	FindAll() ([]*domain.RecipeDTO, error)
	Save(receiptDto *domain.RecipeDTO) (*domain.RecipeDTO, error)
	FindByName(name string) ([]*domain.RecipeDTO, error)
	SaveAll(receiptDtos []domain.RecipeDTO) error
}

// NewSqliteRecipeRepository -
func NewSqliteRecipeRepository(dbConn *gorm.DB) RecipeRepository {
	return &DbTemplate{
		db: dbConn,
	}
}

// Find -
func (template *DbTemplate) Find(id string) (*domain.RecipeDTO, error) {
	receiptDto := domain.NewRecipeDTO()
	template.db.First(receiptDto, id)
	return receiptDto, nil
}

// Find -
func (template *DbTemplate) SaveAll(receiptDtos []domain.RecipeDTO) error {
	db := template.db.DB()

	if db == nil {
		// failed to get underlying db
		return errors.New("db connection is Null")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(INSERT_STMT)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // danger!

	err = persistItems(tx, stmt, receiptDtos)
	if err != nil {
		return err
	}
	return nil
}

func persistItems(tx *sql.Tx, stmt *sql.Stmt, items []domain.RecipeDTO) error {
	defer tx.Rollback()
	for i := 0; i < len(items); i++ {
		recipeDto := items[i]
		_, err := stmt.Exec(recipeDto.ID,
			recipeDto.Name,
			recipeDto.Ingredients,
			recipeDto.URL,
			recipeDto.Image,
			recipeDto.CookTime,
			recipeDto.Source,
			recipeDto.RecipeYield,
			recipeDto.DatePublished,
			recipeDto.PrepTime,
			recipeDto.Description,
			recipeDto.TotalTime,
			time.Now(),
			time.Now(),
			nil)
		if err != nil {
			// log.Fatal(err)
			return err
		}
	}
	// return tx.Commit().Error
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// CreateTableIFNotExists -
func (template *DbTemplate) CreateTableIFNotExists() error {
	hasTable := template.db.HasTable(&domain.RecipeDTO{})
	if !hasTable {
		template.db.CreateTable(&domain.RecipeDTO{})
		hasTable = template.db.HasTable(&domain.RecipeDTO{})
		if !hasTable {
			return errors.New("Failed to create Table for RecipeDTO")
		}
	}
	return nil
}

// FindAll -
func (template *DbTemplate) FindAll() ([]*domain.RecipeDTO, error) {
	results := make([]*domain.RecipeDTO, 0)
	err := template.db.Find(&results)
	if err != nil {
		fmt.Printf("Error when looking up Table, the error is '%v'", err)
	}
	return results, nil
}

// Save -
func (template *DbTemplate) Save(recipeDTO *domain.RecipeDTO) (*domain.RecipeDTO, error) {
	result := domain.NewRecipeDTO()
	if template.db.NewRecord(&recipeDTO) {
		// is not a new record update the recipeDTO
	} else {
		template.db.Create(recipeDTO)
	}
	return result, nil
}

// FindByName -
func (template *DbTemplate) FindByName(name string) ([]*domain.RecipeDTO, error) {
	results := make([]*domain.RecipeDTO, 0)
	err := template.db.Where("name = ?", name).Find(&results)
	if err != nil {
		// log.Debugf("Error when looking up Table, the error is '%v'", err)
		fmt.Printf("Error when looking up Table, the error is '%v'", err)
	}
	return results, nil
}
