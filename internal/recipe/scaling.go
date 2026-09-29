package recipe

import "math"


func (q Quantity) Scale(factor float64) Quantity {
	new := q

	gcd := func(a, b int) int {
		for b != 0 {
			a, b = b, a % b
		}

		return a
	}

	switch new.Kind {
	case Whole, Decimal:
		new.Value *= factor

		if new.Value == math.Trunc(new.Value) {
			new.Kind = Whole
		} else {
			new.Kind = Decimal
		}

	case WholeFrac, Fraction:
		if factor != math.Trunc(factor) {
			new.Kind = Decimal
			new.Value = (float64(new.Num) / float64(new.Den) + new.Value) * factor
			return new
		}

		new.Value *= factor
		new.Num *= int64(factor)

		if new.Num >= new.Den {
			new.Value += float64(new.Num / new.Den)
			new.Num %= new.Den

			switch {
			case new.Num == 0:
				new.Kind = Whole

			case new.Num > 0 && new.Value > 0:
				new.Kind = WholeFrac

			case new.Num > 0 && new.Value == 0:
				new.Kind = Fraction
			}
		}

		if t := gcd(int(new.Num), int(new.Den)); t > 1 {
			new.Num /= int64(t)
			new.Den /= int64(t)
		}
	}

	return new
}


func NormUnit(q Quantity, u IngredientUnit) (Quantity, IngredientUnit) {
	if q.Kind == Fraction || q.Kind == WholeFrac {
		return q, u
	}

	switch u {
	case G, GR, ML:
		if (q.Kind == Whole || q.Kind == Decimal) && q.Value >= 1000 {
			q.Value /= 1000

			if u == GR || u == G {
				u = KG
			} else {
				u = LT
			}
		}

	case KG, LT:
		if q.Kind == Decimal && q.Value < 1 {
			q.Value *= 1000

			if u == KG {
				u = GR
			} else {
				u = ML
			}
		}

	case TSP:
		if q.Kind == Whole && (int64(q.Value) % 3 == 0) {
			q.Value /= 3
			u = TBSP
		}
	}

	if q.Value == math.Trunc(q.Value) {
		q.Kind = Whole
	} else {
		q.Kind = Decimal
	}

	return q, u
}
