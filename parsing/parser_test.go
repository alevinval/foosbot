package parsing_test

import (
	"bytes"
	"github.com/alevinval/foosbot/parsing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCorrectCommand(t *testing.T) {
	cmd := "match alex joaquin 2 vs 1 samuel jordi"
	r := bytes.NewReader([]byte(cmd))
	p := parsing.NewParser(r)
	_, err := p.ParseMatch()
	assert.Nil(t, err)
}

func TestIncorrectCommandMatch(t *testing.T) {
	cmd := "foosbot matchs alex joaquin 2 vs 1 samuel jordi"
	r := bytes.NewReader([]byte(cmd))
	p := parsing.NewParser(r)
	_, err := p.ParseMatch()
	assert.NotNil(t, err)
}

func TestIncorrectCommandT1Players(t *testing.T) {
	cmd := "match alex jordi joaquin 2 vs 1 samuel jordi"
	r := bytes.NewReader([]byte(cmd))
	p := parsing.NewParser(r)
	_, err := p.ParseMatch()
	assert.NotNil(t, err)
}

func TestIncorrectCommandT1Score(t *testing.T) {
	cmd := "match alex joaquin vs 1 samuel jordi"
	r := bytes.NewReader([]byte(cmd))
	p := parsing.NewParser(r)
	_, err := p.ParseMatch()
	assert.NotNil(t, err)
}

func TestIncorrectCommandMissingVS(t *testing.T) {
	cmd := "match alex joaquin 2 1 samuel jordi"
	r := bytes.NewReader([]byte(cmd))
	p := parsing.NewParser(r)
	_, err := p.ParseMatch()
	assert.NotNil(t, err)
}

func TestIncorrectCommandT2Score(t *testing.T) {
	cmd := "match alex joaquin 2 vs _ samuel jordi"
	r := bytes.NewReader([]byte(cmd))
	p := parsing.NewParser(r)
	_, err := p.ParseMatch()
	assert.NotNil(t, err)
}

func TestIncorrectCommandT2Players(t *testing.T) {
	cmd := "match alex joaquin 2 vs 2 jordi"
	r := bytes.NewReader([]byte(cmd))
	p := parsing.NewParser(r)
	_, err := p.ParseMatch()
	assert.NotNil(t, err)
}

func TestIncorrectCommandDuplicatePlayerSameTeam(t *testing.T) {
	cmd := "match alex alex 2 vs 2 joan jordi"
	r := bytes.NewReader([]byte(cmd))
	p := parsing.NewParser(r)
	_, err := p.ParseMatch()
	assert.NotNil(t, err)

	cmd = "match alex joaquin 2 vs 2 jordi jordi"
	r = bytes.NewReader([]byte(cmd))
	p = parsing.NewParser(r)
	_, err = p.ParseMatch()
	assert.NotNil(t, err)
}

func TestIncorrectCommandDuplicatePlayersCrossTeam(t *testing.T) {
	cmd := "match alex joaquin 2 vs 2 jordi alex"
	r := bytes.NewReader([]byte(cmd))
	p := parsing.NewParser(r)
	_, err := p.ParseMatch()
	assert.NotNil(t, err)
}
