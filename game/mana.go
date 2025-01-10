package game

const (
	GreenMana     = "mana:green"
	BlueMana      = "mana:blue"
	RedMana       = "mana:red"
	WhiteMana     = "mana:white"
	BlackMana     = "mana:black"
	ColorlessMana = "mana:colorless"
)

/*
Mana Abstraction of a in-game mana. The count indicates the number of mana we have for an individual color,
this ensures that we don't have to do additional operations on the players mana pool
*/
type Mana struct {
	Color string
	Count int
}

/*
NewMana Create a new mana object and return it
*/
func NewMana(color string, count int) Mana {
	return Mana{
		Color: color,
		Count: count,
	}
}
