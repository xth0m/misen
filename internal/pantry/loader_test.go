package pantry

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/xth0m/misen/internal/recipe"
)

func TestParseIngredient(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected PantryIngredient
	}{
		{
			name:  "standard with quantity, expiration date and tag",
			input: "milk qnty:4lt expiration-date:2026-10-01 +closed",
			expected: PantryIngredient{
				Name:    "milk",
				Qnty:    4,
				Unit:    recipe.LT,
				ExpDate: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
				Tags:    []string{"closed"},
			},
		},
		{
			name:  "multi-word ingredient name",
			input: "extra virgin olive oil qnty:500ml",
			expected: PantryIngredient{
				Name: "extra virgin olive oil",
				Qnty: 500,
				Unit: recipe.ML,
			},
		},
		{
			name:  "unitless quantity defaults to piece",
			input: "eggs qnty:12",
			expected: PantryIngredient{
				Name: "eggs",
				Qnty: 12,
				Unit: recipe.PIECE,
			},
		},
		{
			name:  "ingredient with only tags",
			input: "sourdough bread +opened +homemade",
			expected: PantryIngredient{
				Name: "sourdough bread",
				Tags: []string{"opened", "homemade"},
			},
		},
		{
			name:  "bare ingredient without attributes or tags",
			input: "garlic",
			expected: PantryIngredient{
				Name: "garlic",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := strings.Fields(tt.input)
			got, _, err := parseIngredient(tokens)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Name != tt.expected.Name {
				t.Errorf("got name %q, want %q", got.Name, tt.expected.Name)
			}

			if got.Qnty != tt.expected.Qnty {
				t.Errorf("got qnty %f, want %f", got.Qnty, tt.expected.Qnty)
			}

			if got.Unit != tt.expected.Unit {
				t.Errorf("got unit %v, want %v", got.Unit, tt.expected.Unit)
			}

			if !got.ExpDate.Equal(tt.expected.ExpDate) {
				t.Errorf("got expiration date %v, want %v", got.ExpDate, tt.expected.ExpDate)
			}

			if len(got.Tags) != len(tt.expected.Tags) {
				t.Fatalf("got %d tags, want %d", len(got.Tags), len(tt.expected.Tags))
			}

			for i := range got.Tags {
				if got.Tags[i] != tt.expected.Tags[i] {
					t.Errorf("got tag[%d] %q, want %q", i, got.Tags[i], tt.expected.Tags[i])
				}
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
			name:        "unknown unit",
			input:       "flour qnty:500xyz",
			expectedMsg: "unknown unit",
		},
		{
			name:        "unit with no quantity number",
			input:       "milk qnty:lt",
			expectedMsg: "there is unit but no quantity",
		},
		{
			name:        "plus with no tag name",
			input:       "bread +",
			expectedMsg: "+ but no tag",
		},
		{
			name:        "key without value",
			input:       "eggs qnty:",
			expectedMsg: "forgot something?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := strings.Fields(tt.input)
			_, _, err := parseIngredient(tokens)

			if err == nil {
				t.Fatalf("expected error containing %q, got none", tt.expectedMsg)
			}

			if err.Error() != tt.expectedMsg {
				t.Errorf("got error message %q, want %q", err.Error(), tt.expectedMsg)
			}
		})
	}
}

func TestParsePantryFile(t *testing.T) {

	input := `@fridge
milk qnty:2lt +opened
eggs qnty:6

@pantry
olive oil qnty:1lt
`

	DB = &PantryDB{
		Data: make([]*PantryIngredient, 0),
		Ids:  make(map[string][]int),
	}

	scanner := bufio.NewScanner(strings.NewReader(input))
	errors := parse("test.pantry", scanner)

	if len(errors) != 0 {
		t.Fatalf("expected 0 errors, got %d: %v", len(errors), errors)
	}

	if len(DB.Data) != 3 {
		t.Fatalf("expected 3 ingredients in DB, got %d", len(DB.Data))
	}

	if DB.Data[0].Loc != "fridge" || DB.Data[0].Name != "milk" {
		t.Errorf("expected first item to be milk in fridge, got %s in %s", DB.Data[0].Name, DB.Data[0].Loc)
	}

	if DB.Data[2].Loc != "pantry" || DB.Data[2].Name != "olive oil" {
		t.Errorf("expected third item to be olive oil in pantry, got %s in %s", DB.Data[2].Name, DB.Data[2].Loc)
	}

	if len(DB.Ids["eggs"]) != 1 || DB.Ids["eggs"][0] != 1 {
		t.Errorf("expected eggs to index at Data[1], got %v", DB.Ids["eggs"])
	}
}
