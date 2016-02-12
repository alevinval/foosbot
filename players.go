package foosbot

type Player struct {
	ID   string `json:"player_id"`
	Name string `json:"name"`
}

var (
	players       = []Player{}
	playerMap     = map[string]Player{}
	playerNameMap = map[string]Player{}
)

func (p *Player) ShortID() string {
	return p.ID[:8]
}

func NewPlayer(name string) Player {
	playerID := hash(name)
	player := Player{
		ID:   playerID,
		Name: name,
	}
	return player
}

func AddPlayer(player Player) {
	_, ok := playerMap[player.ID]
	if ok {
		return
	}
	players = append(players, player)
	playerMap[player.ID] = player
	playerNameMap[player.Name] = player
	return
}

func PlayerByID(playerID string) (player Player, ok bool) {
	player, ok = playerMap[playerID]
	return
}

func PlayerByName(name string) (player Player, ok bool) {
	player, ok = playerNameMap[name]
	return
}
