// Stores player info
// Start the database
package database

type Player struct {
	ID       string
	Username string
	Email    string
	Score    int
}

type MySQL interface {
	AddPlayer(username string, password string, email string) (string, error)
	UpdatePlayerScore(id string, score int) error
	ListPlayers() []string //list of player ids
	GetPlayer(id string) *Player
	ValidateLogin(username, password string) error
}
