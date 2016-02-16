package foosbot_test

import (
	"github.com/alevinval/foosbot"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOutcome(t *testing.T) {
	winner, looser := newTeam("team a"), newTeam("team b")
	outcome, err := foosbot.NewOutcome(winner, looser)
	assert.Nil(t, err)
	assert.Equal(t, foosbot.BuildOutcomeID(winner, looser), outcome.ID)
	assert.Equal(t, foosbot.BuildOutcomeID(winner, looser)[:8], outcome.ShortID())
	assert.Equal(t, outcome.WinnerID, winner.ID)
	assert.Equal(t, outcome.LooserID, looser.ID)
	assert.Equal(t, 0, outcome.Occurrences)
	assert.True(t, outcome.IsWinner(winner))
	assert.True(t, outcome.IsLooser(looser))
}

func TestNewOutcomeDuplicated(t *testing.T) {
	winner, looser := newTeam("team a"), newTeam("team a")
	outcome, err := foosbot.NewOutcome(winner, looser)
	assert.Nil(t, outcome)
	assert.Equal(t, err, foosbot.ErrOutcomeImpossible)
}
