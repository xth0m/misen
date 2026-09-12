package recipe


type IngredientUnit int
const (
	PIECE IngredientUnit = iota
	GR
	G
	KG
	ML
	L
	LT
	TSP
	TBSP
	CUP
)

var unitMap = map[string]IngredientUnit{
	"piece": PIECE,
	"gr":    GR, 	"g": G,
	"kg":    KG,
	"ml":    ML,
	"lt":    LT,	"l": L,
	"tsp":   TSP,
	"tbsp":  TBSP,
	"cup":   CUP,
}

var unitStr = func() map[IngredientUnit]string {
	m := make(map[IngredientUnit]string, len(unitMap))

	for s, u := range unitMap {
		m[u] = s
	}

	return m
}()


func (u IngredientUnit) String() string {
	return unitStr[u]
}


func IsUnit(str string) (IngredientUnit, bool) {
	u, ok := unitMap[string(str)]
	return u, ok
}



