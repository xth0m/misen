package pantry

import (
	"time"

	"github.com/xth0m/misen/internal/recipe"
)



type (
	PantryIngredient struct {
		Name string
		Loc string
		Qnty float64
		Unit recipe.IngredientUnit
		ExpDate time.Time
		Tags []string
	}

	PantryDB struct {
		Data	[]*PantryIngredient
		Ids		map[string][]int
	}
)

var DB *PantryDB = nil

