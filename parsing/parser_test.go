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
	p := newParser("p1 p2 24 vs 32 p3 p4")
	statement, err := p.ParseMatch()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(statement.TeamOne))
	assert.Equal(t, 24, statement.TeamOneScore)
	assert.Equal(t, 2, len(statement.TeamTwo))
	assert.Equal(t, 32, statement.TeamTwoScore)

	p = newParser("1 p2 2 vs 1 p3 p4")
	_, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 2 2 vs 1 p3 p4")
	_, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p2 2* vs 1 p3 p4")
	_, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p2 2 _ 1 p3 p4")
	_, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p2 2 vs x1 p3 p4")
	_, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p2 2 vs 1 p4")
	_, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p2 2 vs 1 p3 2")
	_, err = p.ParseMatch()
	assert.NotNil(t, err)

	p = newParser("p1 p2 2 vs 1 p1 p2")
	s, err := p.ParseMatch()
	assert.Nil(t, err)
	err = s.Execute(foosbot.NewContext())
	assert.NotNil(t, err)

	p = newParser("p1 p1 2 vs 1 p3 p4")
	s, err = p.ParseMatch()
	assert.Nil(t, err)
	err = s.Execute(foosbot.NewContext())
	assert.NotNil(t, err)

	p = newParser("p1 p2 2 vs 1 p3 p3")
	s, err = p.ParseMatch()
	assert.Nil(t, err)
	err = s.Execute(foosbot.NewContext())
	assert.NotNil(t, err)
}

func TestParseStatsCommand(t *testing.T) {
	p := newParser("p1 p2")
	s, err := p.ParseStats()
	assert.Nil(t, err)
	assert.Equal(t, []string{"p1", "p2"}, s.Names)

	p = newParser("p1")
	s, err = p.ParseStats()
	assert.Nil(t, err)
	assert.Equal(t, []string{"p1"}, s.Names)

	p = newParser("1 p2")
	_, err = p.ParseStats()
	assert.NotNil(t, err)

	p = newParser("1a p2")
	_, err = p.ParseStats()
	assert.NotNil(t, err)

	p = newParser("a1 2")
	_, err = p.ParseStats()
	assert.NotNil(t, err)
}
