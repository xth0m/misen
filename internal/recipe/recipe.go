package recipe



type RecipeNode interface {
	Node()
}

type RecipeText struct {
	Content string
}

func (RecipeText) Node() {}

type RecipeIngredient struct {
	Name string
	Qnty float64
	Unit IngredientUnit
	Note string
}

func (RecipeIngredient) Node() {}



type RecipeSection struct {
	Ingredients []int // store the index for Recipe:Ingredients
	Instructions []RecipeNode
}


type Recipe struct {
	Ingredients []RecipeIngredient

	Metadata	map[string]string
	Sections	[]RecipeSection
}

