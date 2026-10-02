package recipe

import (
	"fmt"
)


type RecipeNode interface {
	Node()
}

type RecipeText struct {
	Content string
}

func (*RecipeText) Node() {}



type QuantityKind int
const (
	Whole QuantityKind = iota
	Decimal
	Fraction
	WholeFrac
)

type Quantity struct {
	Kind QuantityKind

	Value float64
	Num, Den int64

	Unit IngredientUnit
}

func (q Quantity) StringAmount() string {
	switch q.Kind {
	case Whole:
		return fmt.Sprintf("%d", int64(q.Value))

	case Decimal:
		return fmt.Sprintf("%g", q.Value)

	case Fraction:
		return fmt.Sprintf("%d/%d", q.Num, q.Den)

	case WholeFrac:
		return fmt.Sprintf("%d %d/%d", int64(q.Value), q.Num, q.Den)
	}

	return ""
}



type RecipeIngredient struct {
	Name string
	Qnty Quantity
	Note string
}


type RecipeIngredientDecl struct {
	Idx int
}

func (*RecipeIngredientDecl) Node() {}


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
	Ingredients []*RecipeIngredient

	Metadata RecipeMetadata
	Sections []*RecipeSection

	Scale float64
}

