package game

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
}
