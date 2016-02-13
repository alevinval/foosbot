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

func AddPlayer(player *Player) {
	_, ok := PlayersMap[player.ID]
	if ok {
		return
	}
	Players = append(Players, player)
	PlayersMap[player.ID] = player
	PlayersNameMap[player.Name] = player
	return
}

func PlayerByID(playerID string) (player *Player, ok bool) {
	player, ok = PlayersMap[playerID]
	return
}

func PlayerByName(name string) (player *Player, ok bool) {
	player, ok = PlayersNameMap[name]
	return
}
