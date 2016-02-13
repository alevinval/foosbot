package foosbot_test

import (
	"github.com/alevinval/foosbot"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterMatch(t *testing.T) {
	foosbot.Reset()

	p1 := foosbot.NewPlayer("p1")
	p2 := foosbot.NewPlayer("p2")
	p3 := foosbot.NewPlayer("p3")
	p4 := foosbot.NewPlayer("p4")
	winner := foosbot.NewTeam(p1, p2)
	looser := foosbot.NewTeam(p3, p4)
	match := foosbot.NewMatch(winner, looser)
	foosbot.AddMatchWithHistory(match)

	match, ok := foosbot.MatchByTeams(winner, looser)
	assert.True(t, ok)
	assert.Equal(t, match.WinnerID, winner.ID)
}
