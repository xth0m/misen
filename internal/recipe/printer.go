package recipe

import (
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)



var Au *aurora.Aurora


func printIngredient(ofd *os.File, ingredient *RecipeIngredient) {
	fmt.Fprintf(ofd, " - ")

	if ingredient.Unit != PIECE {
		fmt.Fprintf(ofd, "%g %s %s", Au.BrightYellow(ingredient.Qnty), Au.Faint(ingredient.Unit), ingredient.Name)
	} else {
		fmt.Fprintf(ofd, "%g %s", Au.BrightYellow(ingredient.Qnty), ingredient.Name)
	}

	if ingredient.Note != "" {
		fmt.Fprintf(ofd, " '%s'", Au.Italic(Au.Faint(ingredient.Note)))
	}
}

func printMetadata(ofd *os.File, metadata *RecipeMetadata) {

	fmt.Fprintf(ofd, "Title: %s\n", metadata.Title)
	fmt.Fprintf(ofd, "Source: %s\n", metadata.Source)
}


func buildInstruction(r *Recipe, instr *RecipeSection) string {
	var sb strings.Builder

	for _, node := range instr.Instructions {
		switch node := node.(type) {
		case RecipeText:
			sb.WriteString(node.Content)

		case RecipeIngredientDecl:
			fmt.Fprintf(&sb, "%s", Au.Bold(Au.Cyan(r.Ingredients[node.Idx].Name)))
		}
	}

	return sb.String()
}


func (r *Recipe) Print(ofd *os.File) {
	if Au == nil {
		Au = aurora.New(aurora.WithColors(false))
	}

	printMetadata(ofd, &r.Metadata)
	
	fmt.Fprintln(ofd)

	for _, ingredient := range r.Ingredients {
		printIngredient(ofd, &ingredient)
		fmt.Fprintln(ofd)
	}

	fmt.Fprintf(ofd, "\n-----------------------\n\n")

	for idx, instruction := range r.Sections {
		output := buildInstruction(r, &instruction)

		numStr := fmt.Sprintf("%*d", 2, idx+1)
		fmt.Fprintf(ofd, "%s %s ", Au.BrightBlue(numStr), Au.Faint("|"))

		for _, ch := range output {
			fmt.Fprintf(ofd, "%c", ch)

			if ch == '\n' {
				fmt.Fprintf(ofd, "%*s %s ", 2, "", Au.Faint("|"))
			}
		}

		fmt.Fprintf(ofd, "\n\n")
	}

}
