package foosbot

import (
	"errors"
	"sort"
)

var (
	ErrTeamNoPlayers       = errors.New("a team must have at least 1 player")
	ErrTeamDuplicatePlayer = errors.New("a team cannot contain duplicated players")
)

type Team struct {
	ID      string    `json:"team_id"`
	Players []*Player `json:"players"`
}

func NewTeam(players ...*Player) (team *Team, err error) {
	if len(players) == 0 {
		err = ErrTeamNoPlayers
		return
	}
	playerNames := make([]string, len(players))
	for i := range players {
		playerNames[i] = players[i].Name
	}
	if repeated(playerNames, []string{}) {
		err = ErrTeamDuplicatePlayer
		return
	}
	team = &Team{
		ID:      BuildTeamID(players...),
		Players: players,
	}
	return
}

func (team *Team) ShortID() string {
	return team.ID[:8]
}

func BuildTeamID(players ...*Player) string {
	playerIds := []string{}
	for _, player := range players {
		playerIds = append(playerIds, player.ID)
	}
	sort.Strings(playerIds)
	return hash(playerIds...)
}
