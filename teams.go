package foosbot

import (
	"fmt"
	"sort"
	"strings"
)

type Team struct {
	ID      string    `json:"team_id"`
	Players []*Player `json:"players"`
}

func (t *Team) ShortID() string {
	return strings.ToUpper(t.ID[:8])
}

func (t *Team) String() string {
	p1 := t.Players[0].Name
	p2 := t.Players[1].Name
	return fmt.Sprintf("Team %s (%s %s)", t.ShortID(), p1, p2)
}

func buildTeamId(players ...*Player) string {
	ids := []string{}
	for _, user := range players {
		ids = append(ids, user.ID)
	}
	sort.Strings(ids)
	return hash(ids...)
}

func NewTeam(players ...*Player) *Team {
	teamID := buildTeamId(players...)
	team := new(Team)
	team.ID = teamID
	team.Players = players
	return team
}

func (c *Context) AddTeam(team *Team) {
	_, ok := c.TeamsMap[team.ID]
	if ok {
		return
	}
	c.Teams = append(c.Teams, team)
	c.TeamsMap[team.ID] = team

	for _, player := range team.Players {
		c.AddPlayer(player)
	}
	return
}

func (c *Context) TeamByID(id string) (team *Team, ok bool) {
	team, ok = c.TeamsMap[id]
	return
}

func (c *Context) TeamByPlayers(players ...*Player) (team *Team, ok bool) {
	teamID := buildTeamId(players...)
	team, ok = c.TeamsMap[teamID]
	return
}
