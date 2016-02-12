package foosbot

import (
	"sort"
)

type Team struct {
	ID      string   `json:"team_id"`
	Players []Player `json:"players"`
}

var (
	teams    = []*Team{}
	teamsMap = map[string]*Team{}
)

func (t *Team) ShortID() string {
	return t.ID[:8]
}

func buildTeamId(players ...Player) string {
	ids := []string{}
	for _, user := range players {
		ids = append(ids, user.ID)
	}
	sort.Strings(ids)
	return hash(ids...)
}

func NewTeam(players ...Player) *Team {
	teamID := buildTeamId(players...)
	team := new(Team)
	team.ID = teamID
	team.Players = players
	return team
}

func AddTeam(team *Team) {
	_, ok := teamsMap[team.ID]
	if ok {
		return
	}
	teams = append(teams, team)
	teamsMap[team.ID] = team

	for _, player := range team.Players {
		AddPlayer(player)
	}
	return
}

func TeamByID(id string) (team *Team, ok bool) {
	team, ok = teamsMap[id]
	return
}

func TeamByPlayers(players ...Player) (team *Team, ok bool) {
	teamID := buildTeamId(players...)
	team, ok = teamsMap[teamID]
	return
}
