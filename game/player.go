package game

import (
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
