package main

import (
	"fmt"
	"learn-go/recipes/domain"
)

func main() {

	recipeJSONString := `{ "_id" : { "$oid" : "5160756b96cc62079cc2db15" }, "name" : "Drop Biscuits and Sausage Gravy", "ingredients" : "Biscuits\n3 cups All-purpose Flour\n2 Tablespoons Baking Powder\n1/2 teaspoon Salt\n1-1/2 stick (3/4 Cup) Cold Butter, Cut Into Pieces\n1-1/4 cup Butermilk\n SAUSAGE GRAVY\n1 pound Breakfast Sausage, Hot Or Mild\n1/3 cup All-purpose Flour\n4 cups Whole Milk\n1/2 teaspoon Seasoned Salt\n2 teaspoons Black Pepper, More To Taste", "url" : "http://thepioneerwoman.com/cooking/2013/03/drop-biscuits-and-sausage-gravy/", "image" : "http://static.thepioneerwoman.com/cooking/files/2013/03/bisgrav.jpg", "ts" : { "$date" : 1365276011104 }, "cookTime" : "PT30M", "source" : "thepioneerwoman", "recipeYield" : "12", "datePublished" : "2013-03-11", "prepTime" : "PT10M", "description" : "Late Saturday afternoon, after Marlboro Man had returned home with the soccer-playing girls, and I had returned home with the..." }`

	recipe := domain.FromJSON(recipeJSONString)
	// fmt.Println(recipe)
	fmt.Println(recipe.Name)
	fmt.Println(recipe.CookTime)
	fmt.Println(recipe.URL)
	fmt.Println(recipe.Ingredients)
	fmt.Println(recipe.Image)
	fmt.Println(recipe.DatePublished)
	fmt.Println(recipe.RecipeID.Oid)
	fmt.Println(recipe.Description)

}
