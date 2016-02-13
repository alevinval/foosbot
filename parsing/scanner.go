package parsing

import (
	"bytes"
	"strings"
)

const (
	eof = rune(0)
)

type Token struct {
	Type    tokenType
	Literal string
	Pos     int
}

type scanner struct {
	r    *bytes.Reader
	pos  int
	last int
}

func newScanner(r *bytes.Reader) *scanner {
	s := new(scanner)
	s.r = r
	return s
}

func (t Token) String() string {
	return t.Type.String()
}

func (s *scanner) Scan() Token {
	ch, err := s.read()
	if err != nil {
		return Token{Type: TokenEOF, Literal: "", Pos: s.pos}
	}
	s.unread()
	if isWhitespace(ch) {
		_, literal := s.scanIdentifier()
		return Token{Type: TokenWhitespace, Literal: literal, Pos: s.pos}
	} else if isDigit(ch) {
		_, literal := s.scanIdentifier()
		return Token{Type: TokenDigit, Literal: literal, Pos: s.pos}
	} else if isLetter(ch) {
		token, literal := s.scanIdentifier()
		return Token{Type: token, Literal: literal, Pos: s.pos}
	} else {
		return Token{Type: TokenIllegal, Literal: string(ch), Pos: s.pos}
	}
}

func (s *scanner) scanIdentifier() (tokenType, string) {
	var buf bytes.Buffer
	ch, err := s.read()
	if err != nil {
		return TokenEOF, string(ch)
	}
	buf.WriteRune(ch)
	for {
		ch, err := s.read()
		if ch == eof || err != nil {
			break
		} else if isWhitespace(ch) {
			break
		} else if !isLetter(ch) && !isDigit(ch) {
			buf.WriteRune(ch)
			return TokenIllegal, buf.String()
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

func (s *scanner) read() (rune, error) {
	ch, n, err := s.r.ReadRune()
	s.pos += n
	s.last = n
	return ch, err
}

func (s *scanner) unread() {
	s.r.UnreadRune()
	s.pos -= s.last
	s.last = 0
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }
func isLetter(ch rune) bool     { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }
func isDigit(ch rune) bool      { return (ch >= '0' && ch <= '9') }
