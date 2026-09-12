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

		Iter ParserIter
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

func (p *Parser) reportErrorAt(line string, pos, loc, col int, msg string) {
	err := ParserErr{
		Filename: p.Filename,
		Cursor: Cursor{
			pos: pos,
			loc: loc,
			col: col,
		},
		line: line,
		msg: msg,
	}

	p.Errors = append(p.Errors, err)
}



func (p *Parser) parseFrontmatter() {
	p.Iter.skipWhitespace()

	start := p.Iter.Cursor

	if !p.Iter.matchStr("---\n") {
		// Should warn to insert a fronmatter?
		return
	}

	index := strings.Index(p.Iter.Source[p.Iter.pos:], "---")
	if index == -1 {
		p.reportErrorAt("", start.pos, start.loc, start.col, "end of frontmatter not found")
		return
	}

	frontmatter := p.Iter.Source[p.Iter.pos:p.Iter.pos+index]

	err := yaml.Unmarshal([]byte(frontmatter), &p.recipe.Metadata)
	if err != nil {
		var (
			loadErr *yaml.LoadError
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

			p.reportErrorAt(line, p.Iter.pos+err.Mark.Index,
							actualLine, actualCol,
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

	return p.recipe, nil
}

func readRecipe(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

