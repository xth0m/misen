package recipe

import (
	"fmt"
	"os"
)

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

type RecipeIngredientDecl struct {
	Idx int
}

func (RecipeIngredientDecl) Node() {}

type RecipeSection struct {
	Ingredients  []int // store the index for Recipe:Ingredients
	Instructions []RecipeNode
}

type RecipeMetadata struct {
	Title  string         `yaml:"title"`
	Source string         `yaml:"source"`
	Tags   []string       `yaml:"tags"`
	Custom map[string]any `yaml:",inline"`
}

type Recipe struct {
	Ingredients []RecipeIngredient

	Metadata RecipeMetadata
	Sections []RecipeSection
}

func (r *Recipe) Print(ofd *os.File) {

	writeIngredient := func(ofd *os.File, ingredient *RecipeIngredient) {

		if ingredient.Unit != PIECE {

			fmt.Fprintf(ofd, "%g %s %s", ingredient.Qnty, ingredient.Unit, ingredient.Name)
		} else {

			fmt.Fprintf(ofd, "%g %s", ingredient.Qnty, ingredient.Name)
		}

		if ingredient.Note != "" {
			fmt.Fprintf(ofd, " '%s'", ingredient.Note)
		}

		fmt.Fprintln(ofd)
	}

	fmt.Fprintf(ofd, "%s\n\n", r.Metadata.Title)

	for _, ingredient := range r.Ingredients {
		writeIngredient(ofd, &ingredient)
	}

	fmt.Fprintln(ofd)

	for _, instruction := range r.Sections {

		if len(instruction.Ingredients) > 0 {

			for _, ingredient := range instruction.Ingredients {
				writeIngredient(ofd, &r.Ingredients[ingredient])
			}

			fmt.Fprintf(ofd, "\n")
		}

		for _, node := range instruction.Instructions {
			switch node := node.(type) {
			case RecipeText:
				fmt.Fprintf(ofd, "%s", node.Content)

			case RecipeIngredientDecl:
				fmt.Fprintf(ofd, "%s", r.Ingredients[node.Idx].Name)
			}
		}

		fmt.Printf("\n\n")
	}

}
