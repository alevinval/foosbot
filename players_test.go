package foosbot_test

import (
	"github.com/alevinval/foosbot"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPlayer(t *testing.T) {
	name := "player 1"
	p := foosbot.NewPlayer(name)
	assert.Equal(t, foosbot.BuildPlayerID(name), p.ID)
	assert.Equal(t, foosbot.BuildPlayerID(name)[:8], p.ShortID())
	assert.Equal(t, name, p.Name)
}
