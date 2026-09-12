package parser

import (
	"fmt"
	"os"
	"strings"
)



func readRecipe(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func Parse(recipeFilename string) error {
	if !strings.HasSuffix(recipeFilename, ".misen") {
		return fmt.Errorf("recipe filename must end with '.misen' extension")
	}

	content, err := readRecipe(recipeFilename)
	if err != nil {
		return err
	}

	fmt.Println(content)

	return nil
}

