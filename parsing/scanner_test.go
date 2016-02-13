package parsing_test

import (
	"bytes"
	"github.com/alevinval/foosbot/parsing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newScanner(input string) *parsing.Scanner {
	r := bytes.NewReader([]byte(input))
	return parsing.NewScanner(r)
}

func TestScanEof(t *testing.T) {
	s := newScanner("")
	token := s.Scan()
	assert.Equal(t, parsing.TokenEOF, token.Type)
	assert.Equal(t, "", token.Literal)
}

func TestScanWhitespaces(t *testing.T) {
	s := newScanner("    ")
	token := s.Scan()
	assert.Equal(t, parsing.TokenWhitespace, token.Type)
	assert.Equal(t, "    ", token.Literal)
}

func TestScanDigits(t *testing.T) {
	s := newScanner("1234")
	token := s.Scan()
	assert.Equal(t, parsing.TokenDigit, token.Type)
	assert.Equal(t, "1234", token.Literal)
}

func TestScanIdentifier(t *testing.T) {
	s := newScanner("abcdef")
	token := s.Scan()
	assert.Equal(t, parsing.TokenIdentifier, token.Type)
	assert.Equal(t, "abcdef", token.Literal)
}

func TestScanLanguage(t *testing.T) {
	s := newScanner("foosbot match stats vs")
	token := s.Scan()
	assert.Equal(t, parsing.TokenKeywordFoosbot, token.Type)
	assert.Equal(t, "foosbot", token.Literal)

	s.Scan()
	token = s.Scan()
	assert.Equal(t, parsing.TokenCommandMatch, token.Type)
	assert.Equal(t, "match", token.Literal)

	s.Scan()
	token = s.Scan()
	assert.Equal(t, parsing.TokenCommandStats, token.Type)
	assert.Equal(t, "stats", token.Literal)

	s.Scan()
	token = s.Scan()
	assert.Equal(t, parsing.TokenKeywordVS, token.Type)
	assert.Equal(t, "vs", token.Literal)
}
