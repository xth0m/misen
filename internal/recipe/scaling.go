package recipe

import "math"



func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a % b
	}

	return a
}



func (q *Quantity) scaleFraction(factor float64) {
	if factor == math.Trunc(factor) {
		q.Num *= int64(factor)
	} else {
		q.Kind = Decimal
		q.Value = (float64(q.Num) / float64(q.Den)) * factor
		return
	}

	if q.Num == q.Den {
		q.Kind = Whole
		q.Value = 1
	} else if q.Num > q.Den {
		if (q.Num % q.Den) == 0 {
			q.Kind = Whole
			q.Value = float64(q.Num / q.Den)
		} else {
			q.Kind = WholeFrac
			q.Value = float64(q.Num / q.Den)
			q.Num = q.Num - (q.Den * int64(q.Value))
		}
	}

	if t := gcd(int(q.Num), int(q.Den)); t > 1 {
		q.Num /= int64(t)
		q.Den /= int64(t)
	}
}

func (q *Quantity) scaleWholeFrac(factor float64) {
	if factor == math.Trunc(factor) {
		q.Value *= factor
		q.Num *= int64(factor)
	} else {
		q.Kind = Decimal
		q.Value = (float64(q.Num) / float64(q.Den) + q.Value) * factor
		return
	}

	if q.Num == q.Den {
		q.Kind = Whole
		q.Value += 1
	} else if q.Num > q.Den {
		if (q.Num % q.Den) == 0 {
			q.Kind = Whole
			q.Value += float64(q.Num / q.Den)
		} else {
			q.Kind = WholeFrac
			q.Value += float64(q.Num / q.Den)
			q.Num %= q.Den
		}
	}

	if t := gcd(int(q.Num), int(q.Den)); t > 1 {
		q.Num /= int64(t)
		q.Den /= int64(t)
	}
}

func (q Quantity) Scale(factor float64) Quantity {
	new := q

	switch new.Kind {
	case Whole, Decimal:
		new.Value *= factor

		if new.Value == math.Trunc(new.Value) {
			new.Kind = Whole
		} else {
			new.Kind = Decimal
		}

	case Fraction:
		new.scaleFraction(factor)

	case WholeFrac:
		new.scaleWholeFrac(factor)
	}

	return new
}


func NormUnit(q Quantity, u IngredientUnit) (Quantity, IngredientUnit) {
	switch u {
	case G, GR:
		if (q.Kind == Whole || q.Kind == Decimal) && q.Value >= 1000 {
			q.Value /= 1000
			u = KG
		}

	case ML:
		if (q.Kind == Whole || q.Kind == Decimal) && q.Value >= 1000 {
			q.Value /= 1000
			u = LT
		}

	case KG:
		if q.Kind == Decimal && q.Value < 1 {
			q.Value *= 1000
			u = GR
		}

	case LT:
		if q.Kind == Decimal && q.Value < 1 {
			q.Value *= 1000
			u = ML
		}
	}

	return q, u
}
