package parsing_test

import (
	"bytes"
	"github.com/alevinval/foosbot"
	"github.com/alevinval/foosbot/parsing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newParser(input string) *parsing.Parser {
	r := bytes.NewReader([]byte(input))
	p := parsing.NewParser(r)
	return p
}

func TestParseCommands(t *testing.T) {
	p := newParser("foosbot match foosbot stats")
	token, err := p.ParseCommand()
	assert.Nil(t, err)
	assert.Equal(t, parsing.TokenCommandMatch, token.Type)

	token, err = p.ParseCommand()
	assert.Nil(t, err)
	assert.Equal(t, parsing.TokenCommandStats, token.Type)
}

func TestParseMatchCommand(t *testing.T) {
	p := newParser("p1 p2 2 vs 1 p3 p4")
	outcomes, _, err := p.ParseMatch()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(outcomes))

	p = newParser("1 p2 2 vs 1 p3 p4")
	outcomes, _, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 2 2 vs 1 p3 p4")
	outcomes, _, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p2 2* vs 1 p3 p4")
	outcomes, _, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p2 2 _ 1 p3 p4")
	outcomes, _, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p2 2 vs x1 p3 p4")
	outcomes, _, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p2 2 vs 1 p4")
	outcomes, _, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p2 2 vs 1 p3 2")
	outcomes, _, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p2 2 vs 1 p1 p2")
	outcomes, _, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p1 2 vs 1 p3 p4")
	outcomes, _, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p2 2 vs 1 p3 p3")
	outcomes, _, err = p.ParseMatch()
	assert.NotNil(t, err)
}

func TestParseStatsCommand(t *testing.T) {
	p := newParser("p1 p2")
	obj, err := p.ParseStats()
	assert.Nil(t, err)
	assert.IsType(t, obj, new(foosbot.Team))
	assert.Equal(t, 2, len(obj.(*foosbot.Team).Players))

	p = newParser("p1")
	obj, err = p.ParseStats()
	assert.Nil(t, err)
	assert.IsType(t, obj, new(foosbot.Player))

	p = newParser("1 p2")
	obj, err = p.ParseStats()
	assert.NotNil(t, err)

	p = newParser("1a p2")
	obj, err = p.ParseStats()
	assert.NotNil(t, err)

	p = newParser("a1 2")
	obj, err = p.ParseStats()
	assert.NotNil(t, err)
}
