package recipe

import "testing"

func TestQuantityScale(t *testing.T) {
	tests := []struct {
		name     string
		input    Quantity
		factor   float64
		expected Quantity
	}{
		// Whole
		{
			name:     "whole scale by 1",
			input:    Quantity{Kind: Whole, Value: 4},
			factor:   1,
			expected: Quantity{Kind: Whole, Value: 4},
		},
		{
			name:     "whole scale by 3",
			input:    Quantity{Kind: Whole, Value: 4},
			factor:   3,
			expected: Quantity{Kind: Whole, Value: 12},
		},

		// Decimal
		{
			name:     "decimal scale by 2",
			input:    Quantity{Kind: Decimal, Value: 0.5},
			factor:   2,
			expected: Quantity{Kind: Whole, Value: 1},
		},
		{
			name:     "decimal scale by 4",
			input:    Quantity{Kind: Decimal, Value: 1.25},
			factor:   4,
			expected: Quantity{Kind: Whole, Value: 5},
		},
		{
			name:     "decimal scale to decimal",
			input:    Quantity{Kind: Decimal, Value: 0.5},
			factor:   3,
			expected: Quantity{Kind: Decimal, Value: 1.5},
		},

		// Fraction -> Fraction
		{
			name:     "fraction scaled and simplified",
			input:    Quantity{Kind: Fraction, Num: 1, Den: 4},
			factor:   2,
			expected: Quantity{Kind: Fraction, Num: 1, Den: 2},
		},

		// Fraction -> Whole
		{
			name:     "fraction to whole equal",
			input:    Quantity{Kind: Fraction, Num: 1, Den: 2},
			factor:   2,
			expected: Quantity{Kind: Whole, Value: 1},
		},
		{
			name:     "fraction to whole multiple",
			input:    Quantity{Kind: Fraction, Num: 1, Den: 2},
			factor:   4,
			expected: Quantity{Kind: Whole, Value: 2},
		},

		// Fraction -> WholeFrac
		{
			name:     "fraction to mixed number",
			input:    Quantity{Kind: Fraction, Num: 3, Den: 4},
			factor:   2,
			expected: Quantity{Kind: WholeFrac, Value: 1, Num: 1, Den: 2},
		},
		{
			name:     "fraction to mixed number 2",
			input:    Quantity{Kind: Fraction, Num: 3, Den: 4},
			factor:   3,
			expected: Quantity{Kind: WholeFrac, Value: 2, Num: 1, Den: 4},
		},

		// WholeFrac -> WholeFrac
		{
			name:     "mixed to mixed number",
			input:    Quantity{Kind: WholeFrac, Value: 1, Num: 3, Den: 4},
			factor:   2,
			expected: Quantity{Kind: WholeFrac, Value: 3, Num: 1, Den: 2},
		},
		{
			name:     "mixed to mixed number 5 1/2 * 3",
			input:    Quantity{Kind: WholeFrac, Value: 5, Num: 1, Den: 2},
			factor:   3,
			expected: Quantity{Kind: WholeFrac, Value: 16, Num: 1, Den: 2},
		},

		// WholeFrac -> Whole
		{
			name:     "mixed to whole equal",
			input:    Quantity{Kind: WholeFrac, Value: 5, Num: 1, Den: 2},
			factor:   2,
			expected: Quantity{Kind: Whole, Value: 11},
		},
		{
			name:     "mixed to whole multiple",
			input:    Quantity{Kind: WholeFrac, Value: 2, Num: 1, Den: 4},
			factor:   4,
			expected: Quantity{Kind: Whole, Value: 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Scale(tt.factor)

			if got.Kind != tt.expected.Kind {
				t.Fatalf("expected kind %v, got %v", tt.expected.Kind, got.Kind)
			}

			switch got.Kind {
			case Whole:
				if got.Value != tt.expected.Value {
					t.Errorf("expected whole %d, got %d", int64(tt.expected.Value), int64(got.Value))
				}

			case Decimal:
				if got.Value != tt.expected.Value {
					t.Errorf("expected decimal %g, got %g", tt.expected.Value, got.Value)
				}

			case Fraction:
				if got.Num != tt.expected.Num || got.Den != tt.expected.Den {
					t.Errorf("expected fraction %d/%d, got %d/%d", tt.expected.Num, tt.expected.Den, got.Num, got.Den)
				}

			case WholeFrac:
				if got.Value != tt.expected.Value || got.Num != tt.expected.Num || got.Den != tt.expected.Den {
					t.Errorf("expected %d %d/%d, got %d %d/%d", int64(tt.expected.Value), tt.expected.Num, tt.expected.Den, int64(got.Value), got.Num, got.Den)
				}
			}
		})
	}
}
