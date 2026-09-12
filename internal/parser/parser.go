package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/xth0m/misen/internal/recipe"
)



type (
	Cursor struct {
		pos, loc, col int
	}

	ParserIter struct {
		Source string
		Cursor
		currLineStartIdx int
	}

	Parser struct {
		Filename string

		Iter ParserIter
		Errors []ParserErr
	}

	ParserErr struct {
		Filename string
		Cursor
		line string

		msg string 
	}
)

func (e ParserErr) String() string {
	return e.msg
}


func NewParser(filename, content string) Parser {
	return Parser{
		Filename: filename,
		Iter: ParserIter{
			Source: content,
			Cursor: Cursor{
				loc: 1, col: 1,
			},
		},
		Errors: make([]ParserErr, 0),
	}
}

func (p *Parser) parse() (*recipe.Recipe, error) {
	return nil, nil
}





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

