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
	FindAll() ([]domain.RecipeDTO, error)
	Save(receiptDto *domain.RecipeDTO) (*domain.RecipeDTO, error)
	FindByName(name string) ([]domain.RecipeDTO, error)
	SaveAll(receiptDtos []domain.RecipeDTO) error
	Count() (int64, error)
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
	log.Println("Looking for item with id:", id)
	template.db.Where("id = ?", id).First(&receiptDto)
	return receiptDto, nil
}

// SaveAll - Saves RecipeDTO's in bulk in Transactions. This is very fast
// due to the fact that it uses prepared statements
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
func (template *DbTemplate) FindAll() ([]domain.RecipeDTO, error) {
	results := []domain.RecipeDTO{}
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
func (template *DbTemplate) FindByName(name string) ([]domain.RecipeDTO, error) {
	results := []domain.RecipeDTO{}
	template.db.Where("name = ?", name).Find(&results)
	// if err != nil {
	// 	// log.Debugf("Error when looking up Table, the error is '%v'", err)
	// 	fmt.Printf("Finding recipes by name, the error is '%v'", err.)
	// }
	return results, nil
}

// Count -
func (template *DbTemplate) Count() (int64, error) {
	var count = 0
	tableName := domain.NewRecipeDTO().TableName()
	template.db.Table(tableName).Count(&count)
	return int64(count), nil
}
