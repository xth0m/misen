package parser

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/xth0m/misen/internal/recipe"
	"github.com/xth0m/misen/internal/util"
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

		hasError bool
		Errors []util.Diagnostic

		recipe *recipe.Recipe
	}
)


func NewParser(filename string) Parser {
	return Parser{
		Filename: filename,
		Iter: ParserIter{
			Source: "",
			Cursor: Cursor{
				loc: 1, col: 1,
			},
		},
		hasError: false,
		Errors: make([]util.Diagnostic, 0),
	}
}

func (p *Parser) reportErrorAt(line string, c Cursor, msg string) {
	if p.hasError {
		return
	}

	err := util.Diagnostic{
		Filename: p.Filename,
		Pos: c.pos,
		Loc: c.loc,
		Col: c.col,
		Line: line,
		Msg: msg,
	}

	p.Errors = append(p.Errors, err)
	p.hasError = true
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

func (p *Parser) parseSection() (*recipe.RecipeSection, bool) {
	node := &recipe.RecipeSection{
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
					&recipe.RecipeText{
						Content: p.Iter.Source[start.pos:p.Iter.pos],
					})
				return node, true
			} else if len(node.Instructions) > 0 {
				return node, true
			} else {
				return nil, false
			}

		case ch == '\n' && p.Iter.match('\n'):
			if start.pos < p.Iter.pos - 2 {
				node.Instructions = append(node.Instructions,
					&recipe.RecipeText{
						Content: p.Iter.Source[start.pos : p.Iter.pos-2],
					})
			}

			return node, true

		case ch == '{':
			if start.pos < p.Iter.pos-1 {
				node.Instructions = append(node.Instructions,
					&recipe.RecipeText{
						Content: p.Iter.Source[start.pos : p.Iter.pos-1],
					})
			}

			ingredient := p.parseIngredient()
			if p.hasError {
				break
			}

			p.recipe.Ingredients = append(p.recipe.Ingredients, ingredient)
			idx := len(p.recipe.Ingredients) - 1

			node.Ingredients = append(node.Ingredients, idx)
			node.Instructions = append(node.Instructions, &recipe.RecipeIngredientDecl{ Idx: idx })

			start = p.Iter.Cursor
		}


		if p.hasError {
			for !p.Iter.matchStr("\n\n") && p.Iter.peek() != 0 {
				p.Iter.next()
			}

			p.hasError = false
			return nil, true
		}
	}
}


func ParseQuantity(raw string) (recipe.Quantity, error) {
	p := NewParser("")
	p.Iter.Source = raw

	qnty, err := p.parseQuantityAmount()
	if err != nil {
		return recipe.Quantity{}, err
	}

	p.Iter.skipWhitespace()

	unit, err := p.parseQuantityUnit(p.Iter.Source[p.Iter.pos:])
	if err != nil {
		return recipe.Quantity{}, err
	}

	qnty.Unit = unit

	return qnty, nil
}

func (p *Parser) parseQuantityAmount() (recipe.Quantity, error) {
	n, isInt, err := p.Iter.lexNum()
	if err != nil {
		return recipe.Quantity{}, err
	}
	
	lexInt := func() (int64, error) {
		n, isInt, err := p.Iter.lexNum()

		if err != nil {
			return 0, err
		}

		if !isInt {
			return 0, fmt.Errorf("expected integer")
		}

		return int64(n), nil
	}

	if isInt {
		if r := p.Iter.peekFar(); r >= '0' && r <= '9' {

			p.Iter.skipWhitespace()

			// 2 1/2 Case
			num, err := lexInt()
			if err != nil {
				return recipe.Quantity{}, err
			}

			if !p.Iter.match('/') {
				return recipe.Quantity{}, fmt.Errorf("expected '/'")
			}

			den, err := lexInt()
			if err != nil {
				return recipe.Quantity{}, err
			}

			if den == 0 {
				return recipe.Quantity{}, fmt.Errorf("invalid fraction, denominator cannot be '0'")
			}

			return recipe.Quantity{
				Kind: recipe.WholeFrac,
				Value: n,
				Num: num,
				Den: den,
			}, nil
		} else if p.Iter.peek() == '/' {
			p.Iter.match('/')

			den, err := lexInt()
			if err != nil {
				return recipe.Quantity{}, err
			}

			if den == 0 {
				return recipe.Quantity{}, fmt.Errorf("invalid fraction, denominator cannot be '0'")
			}

			return recipe.Quantity{ Kind: recipe.Fraction, Num: int64(n), Den: den }, nil
		}
	
		return recipe.Quantity{ Kind: recipe.Whole, Value: n }, nil
	} else {
		return recipe.Quantity{ Kind: recipe.Decimal, Value: n }, nil
	}
}

func (p *Parser) parseQuantityUnit(unit string) (recipe.IngredientUnit, error) {
	unit = strings.TrimSpace(unit)

	if len(unit) == 0 {
		return recipe.PIECE, nil
	} else if u, ok := recipe.IsUnit(unit); ok {
		return u, nil
	} else {
		// p.reportError(start, "unknown unit")
		return recipe.NONE, fmt.Errorf("unknown unit")
	}
}


func (p *Parser) parseIngredient() *recipe.RecipeIngredient {
	ingredient := &recipe.RecipeIngredient{}

	start := p.Iter.Cursor

	for {
		ch := p.Iter.peek()

		if ch != 0 && ch != '\n' && ch != '`' && ch != ':' && ch != '}' {
			p.Iter.next()
		} else {
			break
		}
	}

	ingredient.Name = strings.TrimSpace(p.Iter.Source[start.pos:p.Iter.pos])

	if len(ingredient.Name) == 0 {
		p.reportError(start, "no ingredient name")
	}


	start = p.Iter.Cursor

	if p.Iter.consume('`') {
		note, ok := p.Iter.until('`')
		if !ok {
			p.reportError(start, "no matching ` found in the line")
		} else {
			ingredient.Note = strings.TrimSpace(note)
			p.Iter.next()
		}
	}

	if p.Iter.consume(':') {
		p.Iter.skipWhitespace()

		qnty, err := p.parseQuantityAmount()
		if err != nil {
			p.reportError(p.Iter.Cursor, err.Error())
		}

		ingredient.Qnty = qnty

		p.Iter.skipWhitespace()

		start := p.Iter.Cursor
		unit, ok := p.Iter.until('}')
		if !ok {
			p.reportError(p.Iter.Cursor, "no } found")
		}

		u, err := p.parseQuantityUnit(unit)
		if err != nil {
			p.reportError(start, err.Error())
		}

		ingredient.Qnty.Unit = u

		if !p.Iter.consume('}') {
			p.reportError(p.Iter.Cursor, "no } found")
		}
	} else if p.Iter.consume('}') {
		ingredient.Qnty = recipe.Quantity{ Kind: recipe.Whole, Value: 1 }
		ingredient.Qnty.Unit = recipe.PIECE
	} else {
		p.reportError(p.Iter.Cursor, "no } found")
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

		if section != nil {
			p.recipe.Sections = append(p.recipe.Sections, section)
		}

		if p.Iter.peek() == 0 {
			break
		}
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
