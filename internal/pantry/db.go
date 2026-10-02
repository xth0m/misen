package pantry

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/xth0m/misen/internal/parser"
	"github.com/xth0m/misen/internal/recipe"
)



type (
	PantryIngredient struct {
		Name string
		Loc string
		Qnty recipe.Quantity
		ExpDate time.Time
		Tags []string
	}

	PantryDB struct {
		Data	[]*PantryIngredient
		Ids		map[string][]int
	}
)

var DB *PantryDB = nil


type Seq[T any] func(yield func(T) bool)

func (db *PantryDB) Where(predicate ...CriterionFunc) Seq[*PantryIngredient] {
	return func(yield func(*PantryIngredient) bool) {
		for _, v := range db.Data {
			res := true
			for _, fn := range predicate {
				fnRes, err := fn(v)
				if err != nil {}

				res = res && fnRes
			}

			if res && !yield(v) {
				return
			}

		}
	}
}

func (seq Seq[T]) Count() int {
	count := 0

	for range seq {
		count++
	}

	return count
}


type CriterionFunc func(*PantryIngredient) (bool, error)

func criterion(k, v string, eq bool) CriterionFunc {
	return func(pi *PantryIngredient) (bool, error) {
		switch k {
		case "loc":
			return (pi.Loc == v) == eq, nil

		case "qnty":
			qnty, err := parser.ParseQuantity(v)
			if err != nil {
				return false, err
			}

			return (pi.Qnty == qnty) == eq, nil

		case "tag":
			return slices.Contains(pi.Tags, v) == eq, nil
		}

		return false, nil
	}
}

func (db *PantryDB) Query(query string) Seq[*PantryIngredient] {
	var funcs []CriterionFunc

	for tk := range strings.FieldsSeq(query) {
		if len(tk) == 1 {
			break
		}

		eq := true
		switch tk[0] {
		case '-':
			eq = false
			tk = tk[1:]
		
		case '+':
			eq = true
			tk = tk[1:]
		}

		var cr CriterionFunc

		if parts := strings.SplitN(tk, ":", 4); len(parts) == 2 {
			if len(parts) != 2 || parts[1] == "" {
				continue
			}

			cr = criterion(parts[0], parts[1], eq)
		} else {
			cr = criterion("tag", tk, eq)
		}

		funcs = append(funcs, cr)
	}

	return db.Where(funcs...)
}


func PrintQueryRes(seq Seq[*PantryIngredient]) {
	grouped := make(map[string][]*PantryIngredient)
	for it := range seq {
		grouped[it.Loc] = append(grouped[it.Loc], it)
	}

	for k, v := range grouped {
		fmt.Printf("@%s\n", k)

		for _, ing := range v {
			fmt.Printf(" - %s %s %s", ing.Qnty.StringAmount(), ing.Qnty.Unit, ing.Name)

			if len(ing.Tags) != 0 {
				fmt.Printf(" %v", ing.Tags)
			}

			fmt.Println()
		}

		fmt.Println()
	}
}

