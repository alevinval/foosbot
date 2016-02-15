package foosbot

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Team struct {
	ID      string    `json:"team_id"`
	Players []*Player `json:"players"`
}

func buildTeamId(players ...*Player) string {
	ids := []string{}
	for _, user := range players {
		ids = append(ids, user.ID)
	}
	sort.Strings(ids)
	return hash(ids...)
}

func NewTeam(players ...*Player) (*Team, error) {
	if len(players) == 0 {
		return nil, errors.New("provide at least 1 player")
	}
	teamID := buildTeamId(players...)
	team := new(Team)
	team.ID = teamID
	team.Players = players
	return team, nil
}

func (t *Team) ShortID() string {
	return strings.ToUpper(t.ID[:8])
}

func (t *Team) String() string {
	p1 := t.Players[0].Name
	p2 := t.Players[1].Name
	return fmt.Sprintf("Team %s (%s %s)", t.ShortID(), p1, p2)
}
