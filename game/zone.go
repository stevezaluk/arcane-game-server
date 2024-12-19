package game

const (
	BattlefieldZoneId = "zone:battlefield"
	ExileZoneId       = "zone:exile"
	GraveyardZoneId   = "zone:graveyard"
	HandZoneId        = "zone:hand"
	DeckZoneId        = "zone:deck"
	CommanderZoneId   = "zone:commander"
)

type Zone struct {
	ZoneId string
	Cards  []*CardObject

	IsPublic  bool
	IsShared  bool
	IsOrdered bool
}
