package game

import (
	"github.com/stevezaluk/mtgjson-models/deck"
	"github.com/stevezaluk/mtgjson-models/user"
)

type DeckObject struct {
	Metadata   *deck.Deck
	Owner      *user.User
	Controller *user.User
	Zone       *Zone

	IsTopCardRevealed bool
}

/*
NewDeck Create a pointer to a new DeckObject and initialize its related Zone
*/
func NewDeck(deck *deck.Deck, owner *user.User) *DeckObject {
	zone, _ := NewZone(DeckZoneId, owner, false, false, true)

	return &DeckObject{
		Metadata:   deck,
		Owner:      owner,
		Controller: owner,
		Zone:       zone,
	}
}
