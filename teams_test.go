package foosbot_test

import (
	"github.com/alevinval/foosbot"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTeam(t *testing.T) {
	p1, p2 := foosbot.NewPlayer("p1"), foosbot.NewPlayer("p2")
	team, err := foosbot.NewTeam(p1, p2)
	assert.Nil(t, err)
	assert.Equal(t, foosbot.BuildTeamID(p2, p1), team.ID)
	assert.Equal(t, foosbot.BuildTeamID(p1, p2)[:8], team.ShortID())
	assert.Equal(t, []*foosbot.Player{p1, p2}, team.Players)
}

func TestNewTeamWithoutPlayers(t *testing.T) {
	team, err := foosbot.NewTeam()
	assert.Nil(t, team)
	assert.Equal(t, foosbot.ErrTeamNoPlayers, err)
}

func TestNewTeamWithDuplicatedPlayers(t *testing.T) {
	p1, p2 := foosbot.NewPlayer("a"), foosbot.NewPlayer("a")
	team, err := foosbot.NewTeam(p1, p2)
	assert.Nil(t, team)
	assert.Equal(t, err, foosbot.ErrTeamDuplicatePlayer)
}
