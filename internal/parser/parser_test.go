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
				Qnty: 500,
				Unit: recipe.GR,
			},
		},
		{
			name:  "name with note, quantity and unit",
			input: "{water 'lukewarm': 300lt}",
			expected: recipe.RecipeIngredient{
				Name: "water",
				Note: "lukewarm",
				Qnty: 300,
				Unit: recipe.LT,
			},
		},
		{
			name:  "unitless defaults to piece",
			input: "{eggs: 4}",
			expected: recipe.RecipeIngredient{
				Name: "eggs",
				Qnty: 4,
				Unit: recipe.PIECE,
			},
		},
		{
			name:  "standalone ingredient defaults to 1 piece",
			input: "{onion}",
			expected: recipe.RecipeIngredient{
				Name: "onion",
				Qnty: 1,
				Unit: recipe.PIECE,
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
				t.Errorf("expected qnty %g, got %g", tt.expected.Qnty, got.Qnty)
			}

			if got.Unit != tt.expected.Unit {
				t.Errorf("expected unit %v, got %v", tt.expected.Unit, got.Unit)
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
