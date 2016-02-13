package parsing

import (
	"bytes"
	"strings"
)

const (
	eof = rune(0)
)

type Scanner struct {
	r *bytes.Reader
}

func NewScanner(r *bytes.Reader) *Scanner {
	s := new(Scanner)
	s.r = r
	return s
}

func (s *Scanner) Scan() (Token, string) {
	ch, err := s.read()
	if err != nil {
		return TokenEOF, ""
	}
	s.unread()
	if isWhitespace(ch) {
		_, literal := s.scanIdentifier()
		return TokenWhitespace, literal
	} else if isDigit(ch) {
		_, literal := s.scanIdentifier()
		return TokenDigit, literal
	} else if isLetter(ch) {
		return s.scanIdentifier()
	} else {
		return TokenEOF, ""
	}
}

func (s *Scanner) scanIdentifier() (Token, string) {
	var buf bytes.Buffer
	ch, err := s.read()
	if err != nil {
		return TokenEOF, ""
	}
	buf.WriteRune(ch)
	for {
		ch, err := s.read()
		if ch == eof || err != nil {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			break
		}
		buf.WriteRune(ch)
	}
	switch strings.ToLower(buf.String()) {
	case TokenCommandMatch.String():
		return TokenCommandMatch, buf.String()
	case TokenCommandStats.String():
		return TokenCommandStats, buf.String()
	case TokenKeywordFoosbot.String(), "foosball", "fb":
		return TokenKeywordFoosbot, buf.String()
	case TokenKeywordVS.String():
		return TokenKeywordVS, buf.String()
	default:
		return TokenIdentifier, buf.String()
	}
}

func (s *Scanner) read() (rune, error) {
	ch, _, err := s.r.ReadRune()
	return ch, err
}

func (s *Scanner) unread() {
	s.r.UnreadRune()
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }
func isLetter(ch rune) bool     { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }
func isDigit(ch rune) bool      { return (ch >= '0' && ch <= '9') }
