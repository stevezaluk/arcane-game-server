package game

import (
	"github.com/stevezaluk/mtgjson-sdk-client/api"
)

const (
	CommanderGameMode = "gamemode:commander"
	ModernGameMode    = "gamemode:modern"
	StandardGameMode  = "gamemode:standard"
)

/*
Game A representation of a single MTG Game. Process's game commands sent from the client
*/
type Game struct {
	Name     string
	GameMode string

	Players []*Player // ordered

	Battlefield *Zone
	Exile       *Zone
	Command     *Zone
	API         *api.MtgjsonApi
}

/*
NewGame Initialize the zones of a new Game and return a pointer to it
*/
func NewGame(lobbyName string, gameMode string) (*Game, error) {
	battlefield, _ := NewZone(BattlefieldZoneId, nil, true, true, false)
	exile, _ := NewZone(ExileZoneId, nil, true, true, false)

	var commandZone *Zone
	if gameMode == CommanderGameMode {
		commandZone, _ = NewZone(CommanderZoneId, nil, true, true, false)
	}

	return &Game{
		Name:        lobbyName,
		GameMode:    gameMode,
		Battlefield: battlefield,
		Exile:       exile,
		Command:     commandZone,
		API:         api.New(),
	}, nil
}
