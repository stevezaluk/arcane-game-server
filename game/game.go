package game

import (
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk-client/api"
	"log/slog"
	"net"
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

/*
AddPlayer Adds a player to the associated game object, and fetches there user data and deck metadata
*/
func (game *Game) AddPlayer(email string, deckCode string, conn *net.Conn) error {
	user, err := game.API.User.GetUser(email)
	if err != nil {
		slog.Error("Error while requesting user data", "email", email, "err", err.Error())
		return err
	}

	deck, err := game.API.Deck.GetDeck(deckCode, user.Email) // this might not properly validate user ownership here
	if err != nil {
		slog.Error("Error while requesting user deck", "email", email, "deckCode", deckCode, "err", err.Error())
		return err
	}

	if deck.MtgjsonApiMeta.Owner != email || deck.MtgjsonApiMeta.Owner != "system" {
		slog.Error("User does not have ownership of the requested deck", "email", email, "deckCode", deckCode)
		return sdkErrors.ErrInvalidPermissions
	}

	player, err := NewPlayer(user, NewDeck(deck, user), conn)
	if err != nil {
		slog.Error("Error while creating deck object for player", "email", email, "deckCode", deckCode)
		return err
	}

	game.Players = append(game.Players, player)

	return nil
}
