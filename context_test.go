package foosbot_test

import (
	"github.com/alevinval/foosbot"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newTeam(names ...string) *foosbot.Team {
	players := make([]*foosbot.Player, len(names))
	for i := range players {
		players[i] = foosbot.NewPlayer(names[i])
	}
	t, _ := foosbot.NewTeam(players...)
	return t
}
func newMatch(winner *foosbot.Team, loosers *foosbot.Team) *foosbot.Match {
	return foosbot.NewMatch(winner, loosers)
}

func TestRegisterMatch(t *testing.T) {
	w, l := newTeam("a", "b"), newTeam("c", "d")
	m := newMatch(w, l)

	c := foosbot.NewContext()
	c.AddMatchWithHistory(m)

	match, ok := c.Query.MatchByTeams(w, l)
	assert.True(t, ok)
	assert.Equal(t, match.WinnerID, w.ID)
}
