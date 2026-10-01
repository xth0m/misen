package parser

import (
	"testing"

	"github.com/xth0m/misen/internal/recipe"
)

func TestParseIngredient(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected recipe.RecipeIngredient
	}{
		{
			name:  "name with quantity and unit",
			input: "{flour: 500gr}",
			expected: recipe.RecipeIngredient{
				Name: "flour",
				Qnty: recipe.Quantity{Kind: recipe.Whole, Value: 500},
				Unit: recipe.GR,
			},
		},
		{
			name:  "name with note, quantity and unit",
			input: "{water `lukewarm`: 300lt}",
			expected: recipe.RecipeIngredient{
				Name: "water",
				Note: "lukewarm",
				Qnty: recipe.Quantity{Kind: recipe.Whole, Value: 300},
				Unit: recipe.LT,
			},
		},
		{
			name:  "unitless defaults to piece",
			input: "{eggs: 4}",
			expected: recipe.RecipeIngredient{
				Name: "eggs",
				Qnty: recipe.Quantity{Kind: recipe.Whole, Value: 4},
				Unit: recipe.PIECE,
			},
		},
		{
			name:  "standalone ingredient defaults to 1 piece",
			input: "{onion}",
			expected: recipe.RecipeIngredient{
				Name: "onion",
				Qnty: recipe.Quantity{Kind: recipe.Whole, Value: 1},
				Unit: recipe.PIECE,
			},
		},
		{
			name:  "decimal quantity",
			input: "{salt: 0.5tsp}",
			expected: recipe.RecipeIngredient{
				Name: "salt",
				Qnty: recipe.Quantity{Kind: recipe.Decimal, Value: 0.5},
				Unit: recipe.TSP,
			},
		},
		{
			name:  "fraction quantity",
			input: "{olive oil: 1/2tbsp}",
			expected: recipe.RecipeIngredient{
				Name: "olive oil",
				Qnty: recipe.Quantity{Kind: recipe.Fraction, Num: 1, Den: 2},
				Unit: recipe.TBSP,
			},
		},
		{
			name:  "mixed fraction quantity",
			input: "{flour: 2 1/2cup}",
			expected: recipe.RecipeIngredient{
				Name: "flour",
				Qnty: recipe.Quantity{Kind: recipe.WholeFrac, Value: 2, Num: 1, Den: 2},
				Unit: recipe.CUP,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			p := NewParser("test.misen")
			p.Iter.Source = tt.input
			p.recipe = &recipe.Recipe{}

			// Consume opening '{' as parseSection does
			p.Iter.next()

			got := p.parseIngredient()

			if got.Name != tt.expected.Name {
				t.Errorf("expected name %q, got %q", tt.expected.Name, got.Name)
			}

			if got.Note != tt.expected.Note {
				t.Errorf("expected note %q, got %q", tt.expected.Note, got.Note)
			}

			if got.Qnty != tt.expected.Qnty {
				t.Errorf("expected qnty %+v, got %+v", tt.expected.Qnty, got.Qnty)
			}

			if got.Unit != tt.expected.Unit {
				t.Errorf("expected unit %v, got %v", tt.expected.Unit, got.Unit)
			}
		})
	}
}

func TestParseIngredientErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedMsg string
	}{
		{
			name:        "denominator zero",
			input:       "{water: 1/0lt}",
			expectedMsg: "invalid fraction, denominator cannot be '0'",
		},
		{
			name:        "mixed fraction denominator zero",
			input:       "{water: 2 1/0lt}",
			expectedMsg: "invalid fraction, denominator cannot be '0'",
		},
		{
			name:        "missing quantity with unit",
			input:       "{salt: tsp}",
			expectedMsg: "expected quantity",
		},
		{
			name:        "empty quantity after colon",
			input:       "{salt: }",
			expectedMsg: "expected quantity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser("test.misen")
			p.Iter.Source = tt.input
			p.recipe = &recipe.Recipe{}

			p.Iter.next()
			p.parseIngredient()

			if len(p.Errors) == 0 {
				t.Fatalf("expected error containing %q, got no errors", tt.expectedMsg)
			}

			if p.Errors[0].Msg != tt.expectedMsg {
				t.Errorf("expected error %q, got %q", tt.expectedMsg, p.Errors[0].Msg)
			}
		})
	}
}

func TestParseSectionInstructions(t *testing.T) {

	input := "In a bowl mix {water: 300lt} and {sugar: 25gr}."

	p := NewParser("test.misen")
	p.Iter.Source = input
	p.recipe = &recipe.Recipe{}

	section, ok := p.parseSection()
	if !ok {
		t.Fatalf("expected section to parse successfully")
	}

	if len(section.Ingredients) != 2 {
		t.Errorf("expected 2 ingredients in section, got %d", len(section.Ingredients))
	}

	if len(section.Instructions) != 5 {
		t.Errorf("expected 5 instruction nodes (text, ing, text, ing, text), got %d", len(section.Instructions))
	}
}
