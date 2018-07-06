package search

import (
	"fmt"
	"learn-go/recipes/domain"
	// "log"

	log "github.com/Sirupsen/logrus"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/analysis/token/edgengram"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
	"github.com/blevesearch/bleve/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/mapping"
	// "github.com/blevesearch/bleve/mapping"
)

// BleveSearch -
type BleveSearch struct {
	index bleve.Index
}

// NewBleveSearch - create a NewBleveSearch struct
func NewBleveSearch(dbPath string) *BleveSearch {
	bleveSearch := &BleveSearch{}

	index, err := bleveSearch.OpenIndex(dbPath)
	if err != nil {
		fmt.Println("Failed to open or create index file")
	}

	bleveSearch.index = index
	return bleveSearch

}

// OpenIndex returns the opened index
func (bleveSearch *BleveSearch) OpenIndex(dbPath string) (bleve.Index, error) {
	index, err := bleve.Open(dbPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Printf("Creating new index...")
		indexMapping := bleveSearch.CreateIndexMapping(dbPath)
		if indexMapping == nil {
			log.Printf("Failed to create mapping")
		}
		index, err := bleve.New(dbPath, indexMapping)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		return index, nil
	}
	bleveSearch.index = index
	return index, nil
}

// CreateIndexMapping creates the initial index
func (bleveSearch *BleveSearch) CreateIndexMapping(databasePath string) *mapping.IndexMappingImpl {
	mapping := bleve.NewIndexMapping()
	mapping = addCustomAnalyzers(mapping)
	mapping = CreateRecipeIndexMapping(mapping)
	return mapping
}

// IndexRecipe creates the initial index
func (bleveSearch *BleveSearch) IndexRecipe(recipeDto *domain.RecipeDTO) (*domain.RecipeDTO, error) {
	err := bleveSearch.index.Index(recipeDto.ID, *recipeDto)
	return recipeDto, err
}

// BatchIndexRecipe creates the initial index
func (bleveSearch *BleveSearch) BatchIndexRecipe(recipeDtos []domain.RecipeDTO) ([]domain.RecipeDTO, error) {
	batch := bleveSearch.index.NewBatch()
	for _, recipeDto := range recipeDtos {
		batch.Index(recipeDto.ID, recipeDto)
	}
	err := bleveSearch.index.Batch(batch)
	if err != nil {
		log.Printf("Failed to add recipeDto to batch")
	}
	return recipeDtos, nil
}

// Search creates the initial index
func (bleveSearch *BleveSearch) Search(searchRequest *bleve.SearchRequest) (*bleve.SearchResult, error) {
	//
	searchResult, err := bleveSearch.index.Search(searchRequest)
	return searchResult, err
}

// CreateRecipeIndexMapping -
func CreateRecipeIndexMapping(indexMapping *mapping.IndexMappingImpl) *mapping.IndexMappingImpl {

	// a generic disabled fieldMapping
	disabledFieldMapping := bleve.NewDocumentDisabledMapping()
	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Store = false
	englishTextFieldMapping.Analyzer = en.AnalyzerName
	// a generic reusable mapping for english text
	enWithEdgeNgram325TextFieldMapping := bleve.NewTextFieldMapping()
	enWithEdgeNgram325TextFieldMapping.Store = false
	enWithEdgeNgram325TextFieldMapping.Analyzer = en.AnalyzerName

	descriptionMapping := bleve.NewTextFieldMapping()
	descriptionMapping.Store = false
	descriptionMapping.Analyzer = "enWithEdgeNgram325"

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	// keywordFieldMapping.Analyzer = keyword.Name

	dateFieldMapping := bleve.NewDateTimeFieldMapping()
	// dateFieldMapping.DateFormat = "2006-01-02T15:04:05.999999-07:00"

	recipeMapping := bleve.NewDocumentMapping()
	// name
	recipeMapping.AddFieldMappingsAt("name", enWithEdgeNgram325TextFieldMapping)
	// description
	recipeMapping.AddFieldMappingsAt("description", descriptionMapping)
	recipeMapping.AddFieldMappingsAt("ingredients", descriptionMapping)
	recipeMapping.AddFieldMappingsAt("recipeyield", englishTextFieldMapping)
	recipeMapping.AddFieldMappingsAt("cooktime", keywordFieldMapping)
	recipeMapping.AddFieldMappingsAt("source", keywordFieldMapping)
	recipeMapping.AddFieldMappingsAt("preptime", keywordFieldMapping)
	recipeMapping.AddFieldMappingsAt("datepublished", dateFieldMapping)

	// disable createdAt and updatedAt from being indexed
	recipeMapping.AddSubDocumentMapping("createdat", disabledFieldMapping)
	recipeMapping.AddSubDocumentMapping("updatedat", disabledFieldMapping)
	recipeMapping.AddSubDocumentMapping("deletedat", disabledFieldMapping)

	//  := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("recipe", recipeMapping)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping
}

func addCustomTokenFilter(indexMapping *mapping.IndexMappingImpl) *mapping.IndexMappingImpl {
	err := indexMapping.AddCustomTokenFilter("edgeNgram325",
		map[string]interface{}{
			"side": edgengram.FRONT,
			"type": edgengram.Name,
			"min":  3.0,
			"max":  25.0,
		})
	if err != nil {
		log.Fatal(err)
	}
	return indexMapping
}

func addCustomAnalyzers(indexMapping *mapping.IndexMappingImpl) *mapping.IndexMappingImpl {
	indexMapping = addCustomTokenFilter(indexMapping)
	err := indexMapping.AddCustomAnalyzer("enWithEdgeNgram325",
		map[string]interface{}{
			"type":      custom.Name,
			"tokenizer": unicode.Name,
			"token_filters": []string{
				en.PossessiveName,
				lowercase.Name,
				en.StopName,
				"edgeNgram325",
			},
		})
	if err != nil {
		log.Fatal(err)
	}

	return indexMapping
}
