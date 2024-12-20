package game

import (
	"github.com/stevezaluk/mtgjson-models/deck"
	"github.com/stevezaluk/mtgjson-models/user"
	"net"
)

type Player struct {
	User *user.User
	Conn *net.Conn

	Library   *DeckObject
	Graveyard *Zone
	Hand      *Zone

	ManaPool []Mana

	LifeTotal          int
	CommanderDamage    int
	PoisonCounters     int
	EnergyCounters     int
	ExperienceCounters int

	IsMonarch   bool
	IsGameOwner bool
}

func NewPlayer(user *user.User, deck *deck.Deck, conn *net.Conn) (*Player, error) {
	deckObject := NewDeck(deck, user)

	player := &Player{
		User:    user,
		Conn:    conn,
		Library: deckObject,
	}

	graveyard, err := NewZone(GraveyardZoneId, player.User, true, false, true)
	if err != nil {
		return nil, err
	}

	hand, err := NewZone(HandZoneId, player.User, false, false, false)
	if err != nil {
		return nil, err
	}

	player.Hand = hand
	player.Graveyard = graveyard

	return player, nil
}
