package parser

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora/v4"
	"github.com/xth0m/misen/internal/recipe"
	"go.yaml.in/yaml/v4"
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

		Iter   ParserIter
		Errors []ParserErr

		recipe *recipe.Recipe
	}

	ParserErr struct {
		Filename string
		Cursor
		line string

		msg string
	}
)

func (e ParserErr) String() string {
	var sb strings.Builder

	gutterWidth := 4

	fmt.Fprintf(&sb, "%s at %s\n", aurora.Red("error"),
		aurora.Sprintf(aurora.Underline("%s:%d:%d"), e.Filename, e.loc, e.col))

	fmt.Fprintf(&sb, "%*s |\n", gutterWidth, "")
	fmt.Fprintf(&sb, "%*d | %s\n", gutterWidth, e.loc, e.line)
	fmt.Fprintf(&sb, "%*s | %*s\n", gutterWidth, "", e.col, "^")
	fmt.Fprintf(&sb, "%*s = %s\n", gutterWidth, "", e.msg)

	return sb.String()
}

func NewParser(filename string) Parser {
	return Parser{
		Filename: filename,
		Iter: ParserIter{
			Source: "",
			Cursor: Cursor{
				loc: 1, col: 1,
			},
		},
		Errors: make([]ParserErr, 0),
	}
}

func (p *Parser) reportErrorAt(line string, c Cursor, msg string) {
	err := ParserErr{
		Filename: p.Filename,
		Cursor:   c,
		line:     line,
		msg:      msg,
	}

	p.Errors = append(p.Errors, err)
}

func (p *Parser) reportError(c Cursor, msg string) {
	idx := strings.Index(p.Iter.Source[p.Iter.currLineStartIdx:], "\n")

	var line string

	if idx == -1 {
		line = p.Iter.Source[p.Iter.currLineStartIdx:]
	} else {
		line = p.Iter.Source[p.Iter.currLineStartIdx : p.Iter.currLineStartIdx+idx]
	}

	p.reportErrorAt(line, c, msg)
}

func (p *Parser) parseFrontmatter() {
	p.Iter.skipWhitespace()

	start := p.Iter.Cursor

	if !p.Iter.matchStr("---\n") {
		// Should warn to insert a frontmatter?
		return
	}

	index := strings.Index(p.Iter.Source[p.Iter.pos:], "---")
	if index == -1 {
		p.reportErrorAt("", start, "end of frontmatter not found")
		return
	}

	frontmatter := p.Iter.Source[p.Iter.pos : p.Iter.pos+index]

	err := yaml.Unmarshal([]byte(frontmatter), &p.recipe.Metadata)
	if err != nil {
		var (
			loadErr  *yaml.LoadError
			loadErrs *yaml.LoadErrors
		)

		reportYAMLError := func(err *yaml.LoadError) {

			actualLine := p.Iter.loc + err.Mark.Line - 1
			actualCol := err.Mark.Column

			lines := strings.Split(frontmatter, "\n")
			var line string

			if err.Mark.Line >= 1 && err.Mark.Line <= len(lines) {
				line = lines[err.Mark.Line-1]
			}

			p.reportErrorAt(line,
				Cursor{
					pos: p.Iter.pos + err.Mark.Index,
					loc: actualLine,
					col: actualCol,
				},
				err.Message)
		}

		switch {
		case errors.As(err, &loadErr):
			reportYAMLError(loadErr)

		case errors.As(err, &loadErrs):
			for _, le := range loadErrs.Errors {
				reportYAMLError(le)
			}
		}
	}

	newlinec := strings.Count(frontmatter, "\n")
	p.Iter.loc += newlinec
	p.Iter.col = 1
	p.Iter.pos += index

	p.Iter.matchStr("---\n")
}

func (p *Parser) parseSection() (recipe.RecipeSection, bool) {
	node := recipe.RecipeSection{
		Ingredients:  make([]int, 0),
		Instructions: make([]recipe.RecipeNode, 0),
	}

	p.Iter.skipWhitespace()

	start := p.Iter.Cursor
	for {
		ch := p.Iter.next()

		switch {
		case ch == 0:
			if start.pos < p.Iter.pos {
				node.Instructions = append(node.Instructions,
					recipe.RecipeText{
						Content: p.Iter.Source[start.pos:p.Iter.pos],
					})
				return node, true
			} else {
				return recipe.RecipeSection{}, false
			}

		case ch == '\n' && p.Iter.match('\n'):
			if start.pos < p.Iter.pos {
				node.Instructions = append(node.Instructions,
					recipe.RecipeText{
						Content: p.Iter.Source[start.pos : p.Iter.pos-2],
					})
			}

			return node, true

		case ch == '{':
			if start.pos < p.Iter.pos-1 {
				node.Instructions = append(node.Instructions,
					recipe.RecipeText{
						Content: p.Iter.Source[start.pos : p.Iter.pos-1],
					})
			}

			ingredient := p.parseIngredient()

			p.recipe.Ingredients = append(p.recipe.Ingredients, ingredient)
			idx := len(p.recipe.Ingredients) - 1

			node.Ingredients = append(node.Ingredients, idx)
			node.Instructions = append(node.Instructions, recipe.RecipeIngredientDecl{Idx: idx})

			start = p.Iter.Cursor
		}
	}
}

func (p *Parser) parseIngredient() recipe.RecipeIngredient {

	ingredient := recipe.RecipeIngredient{}

	p.Iter.skipWhitespace()

	start := p.Iter.Cursor

	for {
		ch := p.Iter.peek()

		if ch != '\'' && ch != ':' && ch != '}' {
			p.Iter.next()
		} else {
			break
		}
	}

	ingredient.Name = strings.TrimSpace(p.Iter.Source[start.pos:p.Iter.pos])

	if len(ingredient.Name) == 0 {
	}

	if p.Iter.consume('\'') {
		note, ok := p.Iter.until('\'')
		if !ok {
		}

		ingredient.Note = strings.TrimSpace(note)
		p.Iter.next()
	}

	if p.Iter.consume(':') {
		p.Iter.skipWhitespace()

		qnty, ok := p.Iter.lexNum()
		if !ok {
		}

		ingredient.Qnty = qnty

		p.Iter.skipWhitespace()

		start := p.Iter.Cursor
		unit, ok := p.Iter.until('}')

		unit = strings.TrimSpace(unit)

		if len(unit) == 0 {

			ingredient.Unit = recipe.PIECE

		} else if u, ok := recipe.IsUnit(unit); ok {

			ingredient.Unit = u

		} else {
			p.reportError(start, "unknown unit")
		}

		p.Iter.consume('}')

	} else if p.Iter.consume('}') {
		ingredient.Qnty = 1
		ingredient.Unit = recipe.PIECE
	}

	return ingredient
}

func (p *Parser) Parse() (*recipe.Recipe, error) {
	if !strings.HasSuffix(p.Filename, ".misen") {
		return nil, fmt.Errorf("recipe filename must end with '.misen' extension")
	}

	source, err := readRecipe(p.Filename)
	if err != nil {
		return nil, err
	}

	p.Iter.Source = source
	p.recipe = &recipe.Recipe{}

	p.parseFrontmatter()

	for {
		p.Iter.skipWhitespace()

		section, ok := p.parseSection()
		if !ok {
			break
		}

		p.recipe.Sections = append(p.recipe.Sections, section)
	}

	return p.recipe, nil
}

func readRecipe(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
