package pantry

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xth0m/misen/internal/parser"
	"github.com/xth0m/misen/internal/util"
)


func parseKeyValue(tk string, ingredient *PantryIngredient) (string, error) {
	parts := strings.SplitN(tk, ":", 2)

	if len(parts) != 2 || parts[1] == "" {
		return parts[0], fmt.Errorf("forgot something?")
	}

	key, value := parts[0], parts[1]

	switch key {
	case "qnty":
		qnty, err := parser.ParseQuantity(value)
		if err != nil {
			return value, err
		}

		ingredient.Qnty = qnty

	case "expiration-date":
		date, err := time.Parse(time.DateOnly, value)
		if err != nil {
			return value, err
		}

		ingredient.ExpDate = date
	}

	return "", nil
}

func parseIngredient(it []string) (*PantryIngredient, string, error) {
	ingredient := &PantryIngredient{}

	var sb strings.Builder
	idx := 0

	for idx < len(it) {
		tk := it[idx]

		if strings.HasPrefix(tk, "+") || strings.Contains(tk, ":") {
			break
		}

		sb.WriteString(tk)
		sb.WriteByte(' ')
		idx++
	}

	ingredient.Name = strings.TrimSpace(sb.String())

	for idx < len(it) {
		tk := it[idx]

		if tk[0] == '+' {
			if len(tk) == 1 {
				return nil, tk, fmt.Errorf("+ but no tag")
			}

			ingredient.Tags = append(ingredient.Tags, tk[1:])
		} else {
			errloc, err := parseKeyValue(tk, ingredient)
			if err != nil {
				return nil, errloc, err
			}
		}

		idx++
	}

	return ingredient, "", nil
}

func parse(filepath string, scanner *bufio.Scanner) []util.Diagnostic {
	var location string
	loc := 0

	var errors []util.Diagnostic

	for scanner.Scan() {
		line := scanner.Text()
		loc++

		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}

		if tokens[0][0] == '@' {
			if len(tokens[0]) > 1 {
				location = tokens[0][1:]
			} else {
				errors = append(errors, util.Diagnostic{
					Filename: filepath,
					Loc: loc,
					Col: 1,
					Line: line,
					Msg: "Did you forget to add the location?",
				})
			}

			continue
		}

		ingredient, errloc, err := parseIngredient(tokens)
		if err != nil {
			errors = append(errors, util.Diagnostic{
				Filename: filepath,
				Loc: loc,
				Line: line,
				Col: strings.Index(line, errloc) + 1,
				Msg: err.Error(),
			})

			continue
		}

		ingredient.Loc = location

		DB.Data = append(DB.Data, ingredient)
		DB.Ids[ingredient.Name] = append(DB.Ids[ingredient.Name], len(DB.Data) - 1)
	}

	return errors
}

func LoadDB(filepath string) error {
	if DB != nil {
		return fmt.Errorf("pantry db already exists")
	}

	DB = &PantryDB{
		Data: make([]*PantryIngredient, 0),
		Ids: make(map[string][]int),
	}

	f, err := os.OpenFile(filepath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	
	errors := parse(filepath, scanner)
	if len(errors) != 0 {
		for _, err := range errors {
			fmt.Fprintln(os.Stderr, err)
		}

		return fmt.Errorf("failed to parse pantry file")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}


