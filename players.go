package foosbot

type Player struct {
	ID   string `json:"player_id"`
	Name string `json:"name"`
}

func NewPlayer(name string) (p *Player) {
	p = &Player{
		ID:   BuildPlayerID(name),
		Name: name,
	}
	return
}

func (p *Player) ShortID() string {
	return p.ID[:8]
}

func BuildPlayerID(name string) string {
	return hash(name)
}
