package foosbot_test

import (
	"github.com/alevinval/foosbot"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterTeam(t *testing.T) {
	foosbot.Reset()

	p1 := foosbot.NewPlayer("p1")
	p2 := foosbot.NewPlayer("p2")
	team := foosbot.NewTeam(p1, p2)
	foosbot.AddTeam(team)

	team, ok := foosbot.TeamByPlayers(p1, p2)
	assert.True(t, ok)
	team, ok = foosbot.TeamByID(team.ID)
	assert.True(t, ok)
	assert.Equal(t, 2, len(team.Players))
}
