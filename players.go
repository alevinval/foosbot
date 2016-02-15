package foosbot

import "strings"

type Player struct {
	ID   string `json:"player_id"`
	Name string `json:"name"`
}

func NewPlayer(name string) *Player {
	player := new(Player)
	player.ID = hash(name)
	player.Name = name
	return player
}

func (p *Player) ShortID() string {
	return strings.ToUpper(p.ID[:8])
}
