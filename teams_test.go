package foosbot_test

import (
	"github.com/alevinval/foosbot"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterTeam(t *testing.T) {
	c := foosbot.NewContext()

	p1 := foosbot.NewPlayer("p1")
	p2 := foosbot.NewPlayer("p2")
	team, _ := foosbot.NewTeam(p1, p2)
	c.AddTeam(team)

	team, ok := c.Query.TeamByPlayers(p1, p2)
	assert.True(t, ok)
	team, ok = c.Query.TeamByID(team.ID)
	assert.True(t, ok)
	assert.Equal(t, 2, len(team.Players))
}
