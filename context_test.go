package foosbot_test

import (
	"github.com/alevinval/foosbot"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newTeam(names ...string) (team *foosbot.Team) {
	players := make([]*foosbot.Player, len(names))
	for i := range players {
		players[i] = foosbot.NewPlayer(names[i])
	}
	team, _ = foosbot.NewTeam(players...)
	return
}

func TestRegisterMatch(t *testing.T) {
	winner, looser := newTeam("a", "b"), newTeam("c", "d")
	outcome, _ := foosbot.NewOutcome(winner, looser)
	ctx := foosbot.NewContext()
	ctx.AddMatchWithOutcome(outcome)

	outcomeByID, ok := ctx.Query.OutcomeByID(outcome.ID)
	assert.True(t, ok)
	assert.Equal(t, outcome, outcomeByID)

	outcomeByTeams, ok := ctx.Query.OutcomeByTeams(winner, looser)
	assert.True(t, ok)
	assert.Equal(t, outcome, outcomeByTeams)
}

func TestRegisterTeam(t *testing.T) {
	ctx := foosbot.NewContext()
	team := newTeam("p1", "p2")
	ctx.AddTeam(team)

	teamByID, ok := ctx.Query.TeamByID(team.ID)
	assert.True(t, ok)
	teamByPlayers, ok := ctx.Query.TeamByPlayers(team.Players...)
	assert.True(t, ok)

	assert.Equal(t, team, teamByID)
	assert.Equal(t, team, teamByPlayers)
}
