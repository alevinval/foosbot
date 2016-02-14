package foosbot

import "strings"

type Player struct {
	ID   string `json:"player_id"`
	Name string `json:"name"`
}

func (p *Player) ShortID() string {
	return strings.ToUpper(p.ID[:8])
}

func NewPlayer(name string) *Player {
	playerID := hash(name)
	player := new(Player)
	player.ID = playerID
	player.Name = name
	return player
}

func (c *Context) AddPlayer(player *Player) {
	_, ok := c.PlayersMap[player.ID]
	if ok {
		return
	}
	c.Players = append(c.Players, player)
	c.PlayersMap[player.ID] = player
	c.PlayersNameMap[player.Name] = player
	return
}

func (c *Context) PlayerByID(playerID string) (player *Player, ok bool) {
	player, ok = c.PlayersMap[playerID]
	return
}

func (c *Context) PlayerByName(name string) (player *Player, ok bool) {
	player, ok = c.PlayersNameMap[name]
	return
}
