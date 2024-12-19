package game

import (
	"github.com/stevezaluk/mtgjson-models/card"
	"github.com/stevezaluk/mtgjson-models/user"
)

/*
CardObject Represents the card played or generated for the game. This needs to differ
from the protobuf models as there are additional values that need to be tracked like
ownership, its parent zone, and the state of the card
*/
type CardObject struct {
	Metadata   *card.CardSet
	Owner      *user.User
	ParentZone *Zone

	IsTapped          bool
	IsFaceDown        bool
	WasPlayedThisTurn bool
}

/*
NewCardObject Create a new pointer to a card object. Its Metadata, Owner, and ParentZone are required
*/
func NewCardObject(metadata *card.CardSet, owner *user.User, zone *Zone) *CardObject {
	return &CardObject{
		Metadata:          metadata,
		Owner:             owner,
		ParentZone:        zone,
		WasPlayedThisTurn: true,
	}
}

/*
TapCard Set IsTapped to true, and consider the card tapped out
*/
func (card *CardObject) TapCard() {
	if !card.IsTapped {
		card.IsTapped = true
	}
}

/*
UnTapCard Set IsTapped to false, and consider the card untapped
*/
func (card *CardObject) UnTapCard() {
	if card.IsTapped {
		card.IsTapped = false
	}
}
