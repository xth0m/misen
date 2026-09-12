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



type RecipeMetadata struct {
	Title string `yaml:"title"`
	Source string `yaml:"source"`
	Tags []string `yaml:"tags"`
	Custom map[string]any `yaml:",inline"`
}

type Recipe struct {
	Ingredients []RecipeIngredient

	Metadata	RecipeMetadata
	Sections	[]RecipeSection
}

