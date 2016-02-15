package foosbot_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatchWinners(t *testing.T) {
	w, l := newTeam("a"), newTeam("b")
	m := newMatch(w, l)

	assert.Equal(t, m.WinnerID, w.ID)
	assert.True(t, m.IsWinner(w))
	assert.True(t, m.IsLooser(l))
	assert.Equal(t, m.Winner(), w)
	assert.Contains(t, m.Loosers(), l)
}
