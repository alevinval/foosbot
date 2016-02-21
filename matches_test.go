package foosbot_test

import (
	"github.com/alevinval/foosbot"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMatch(t *testing.T) {
	winner, looser := newTeam("a"), newTeam("b")
	outcome, _ := foosbot.NewOutcome(winner, looser)
	match := foosbot.NewMatch(outcome)
	assert.Equal(t, foosbot.BuildMatchID(outcome, match.PlayedAt), match.ID)
	assert.Equal(t, foosbot.BuildMatchID(outcome, match.PlayedAt)[:8], match.ShortID())
	assert.Equal(t, outcome.ID, match.OutcomeID)
}
