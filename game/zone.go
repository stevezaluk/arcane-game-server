package game

import "github.com/stevezaluk/mtgjson-models/user"

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
	Owner  *user.User
	Cards  []*CardObject

	IsPublic  bool
	IsShared  bool
	IsOrdered bool
}
