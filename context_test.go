package foosbot_test

import (
	"github.com/alevinval/foosbot"
	"github.com/alevinval/foosbot/parsing"
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

	outcomeByID, ok := foosbot.Query(ctx).OutcomeByID(outcome.ID)
	assert.True(t, ok)
	assert.Equal(t, outcome, outcomeByID)

	outcomeByTeams, ok := foosbot.Query(ctx).OutcomeByTeams(winner, looser)
	assert.True(t, ok)
	assert.Equal(t, outcome, outcomeByTeams)
}

func TestRegisterTeam(t *testing.T) {
	ctx := foosbot.NewContext()
	team := newTeam("p1", "p2")
	ctx.AddTeam(team)

	teamByID, ok := foosbot.Query(ctx).TeamByID(team.ID)
	assert.True(t, ok)
	teamByPlayers, ok := foosbot.Query(ctx).TeamByPlayers(team.Players...)
	assert.True(t, ok)

	assert.Equal(t, team, teamByID)
	assert.Equal(t, team, teamByPlayers)

	// Try registering twice, should be ok.
	ctx.AddTeam(team)
	teamByIDSecond, ok := foosbot.Query(ctx).TeamByID(team.ID)
	assert.True(t, ok)
	assert.Equal(t, teamByID, teamByIDSecond)
}

func TestAddOutcomeManyTimes(t *testing.T) {
	ctx := foosbot.NewContext()
	winner, looser := newTeam("a"), newTeam("b")
	outcome, _ := foosbot.NewOutcome(winner, looser)
	ctx.AddMatchWithOutcome(outcome)
	ctx.AddMatchWithOutcome(outcome)
	ctx.AddMatchWithOutcome(outcome)

	outcomeByID, ok := foosbot.Query(ctx).OutcomeByID(outcome.ID)
	assert.True(t, ok)
	assert.Equal(t, 3, outcomeByID.Occurrences)
}

func TestReset(t *testing.T) {
	ctx := foosbot.NewContext()
	winner, looser := newTeam("a"), newTeam("b")
	outcome, _ := foosbot.NewOutcome(winner, looser)
	ctx.AddTeam(winner)
	ctx.AddTeam(looser)
	ctx.AddMatchWithOutcome(outcome)

	_, ok := foosbot.Query(ctx).OutcomeByID(outcome.ID)
	assert.True(t, ok)
	_, ok = foosbot.Query(ctx).TeamByID(winner.ID)
	assert.True(t, ok)
	_, ok = foosbot.Query(ctx).TeamByID(looser.ID)
	assert.True(t, ok)
	_, ok = foosbot.Query(ctx).OutcomeByID(outcome.ID)
	assert.True(t, ok)
	player, ok := foosbot.Query(ctx).PlayerByName("a")
	assert.True(t, ok)
	_, ok = foosbot.Query(ctx).PlayerByID(player.ID)
	assert.True(t, ok)
	ctx.Reset()
	_, ok = foosbot.Query(ctx).OutcomeByID(outcome.ID)
	assert.False(t, ok)
	_, ok = foosbot.Query(ctx).TeamByID(winner.ID)
	assert.False(t, ok)
	_, ok = foosbot.Query(ctx).TeamByID(looser.ID)
	assert.False(t, ok)
	_, ok = foosbot.Query(ctx).OutcomeByID(outcome.ID)
	assert.False(t, ok)
	player, ok = foosbot.Query(ctx).PlayerByName("a")
	assert.False(t, ok)
	assert.Nil(t, player)
}

func TestExecuteMatchStatement(t *testing.T) {
	ctx := foosbot.NewContext()

	statement := parsing.MatchStatement{
		TeamOne:      []string{"p1", "p2"},
		TeamOneScore: 3,
		TeamTwo:      []string{"p3", "p4"},
		TeamTwoScore: 5,
	}
	err := ctx.ExecuteMatch(statement)
	assert.Nil(t, err)

	statement = parsing.MatchStatement{
		TeamOne:      []string{"p1", "p1"},
		TeamOneScore: 2,
		TeamTwo:      []string{"p3", "p4"},
		TeamTwoScore: 1,
	}
	err = ctx.ExecuteMatch(statement)
	assert.NotNil(t, err)

	statement = parsing.MatchStatement{
		TeamOne:      []string{"p1", "p2"},
		TeamOneScore: 2,
		TeamTwo:      []string{"p3", "p3"},
		TeamTwoScore: 1,
	}
	err = ctx.ExecuteMatch(statement)
	assert.NotNil(t, err)
}

func TestExecuteStatsStatement(t *testing.T) {
	ctx := foosbot.NewContext()
	statement := parsing.StatStatement{
		Names: []string{"p1", "p2"},
	}
	_, err := ctx.ExecuteStats(statement)
	assert.Nil(t, err)
}
