package foosbot_test

import (
	"github.com/alevinval/foosbot"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterMatch(t *testing.T) {
	c := foosbot.NewContext()

	p1 := foosbot.NewPlayer("p1")
	p2 := foosbot.NewPlayer("p2")
	p3 := foosbot.NewPlayer("p3")
	p4 := foosbot.NewPlayer("p4")
	winner := foosbot.NewTeam(p1, p2)
	looser := foosbot.NewTeam(p3, p4)
	match := foosbot.NewMatch(winner, looser)
	c.AddMatchWithHistory(match)

	match, ok := c.MatchByTeams(winner, looser)
	assert.True(t, ok)
	assert.Equal(t, match.WinnerID, winner.ID)
}
