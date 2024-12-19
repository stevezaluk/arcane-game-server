package game

import (
	arcaneErrors "github.com/stevezaluk/arcane-game-server/errors"
	"github.com/stevezaluk/mtgjson-models/user"
)

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

/*
NewZone Create a pointer to a new Zone. Owner can be nil provided that the zone is not shared
*/
func NewZone(zoneId string, owner *user.User, isPublic bool, isShared bool, isOrdered bool) (*Zone, error) {
	if owner != nil && isShared {
		return nil, arcaneErrors.ErrZoneCannotBeShared
	}

	return &Zone{
		ZoneId:    zoneId,
		Owner:     owner,
		IsPublic:  isPublic,
		IsShared:  isShared,
		IsOrdered: isOrdered,
	}, nil
}
