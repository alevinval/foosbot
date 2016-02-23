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

func (t Token) String() string {
	return t.Type.String()
}

type Scanner struct {
	r    *bytes.Reader
	pos  int
	last int
}

func NewScanner(r *bytes.Reader) *Scanner {
	s := new(Scanner)
	s.r = r
	return s
}

func (s *Scanner) Scan() Token {
	ch, err := s.read()
	if err != nil {
		return Token{Type: TokenEOF, Literal: "", Pos: s.pos}
	}
	s.unread()
	if isWhitespace(ch) {
		token, literal := s.scanWhitespaces()
		return Token{Type: token, Literal: literal, Pos: s.pos}
	} else if isDigit(ch) {
		token, literal := s.scanDigit()
		return Token{Type: token, Literal: literal, Pos: s.pos}
	} else if isLetter(ch) {
		token, literal := s.scanIdentifier()
		return Token{Type: token, Literal: literal, Pos: s.pos}
	} else {
		return Token{Type: TokenIllegal, Literal: string(ch), Pos: s.pos}
	}
}

func (s *Scanner) scanDigit() (tokenType, string) {
	var buf bytes.Buffer
	for {
		ch, err := s.read()
		if ch == eof || err != nil {
			break
		} else if isWhitespace(ch) {
			s.unread()
			break
		} else if !isDigit(ch) {
			s.unread()
			buf.WriteRune(ch)
			return TokenIllegal, buf.String()
		}
		buf.WriteRune(ch)
	}
	return TokenDigit, buf.String()
}

func (s *Scanner) scanWhitespaces() (tokenType, string) {
	var buf bytes.Buffer
	for {
		ch, err := s.read()
		if ch == eof || err != nil {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		}
		buf.WriteRune(ch)
	}
	return TokenWhitespace, buf.String()
}

func (s *Scanner) scanIdentifier() (tokenType, string) {
	var buf bytes.Buffer
	for {
		ch, err := s.read()
		if ch == eof || err != nil {
			break
		} else if isWhitespace(ch) {
			s.unread()
			break
		} else if !isLetter(ch) && !isDigit(ch) {
			s.unread()
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
	case TokenCommandLeaderboard.String():
		return TokenCommandLeaderboard, buf.String()
	case TokenKeywordFoosbot.String(), "foosball", "fb":
		return TokenKeywordFoosbot, buf.String()
	case TokenKeywordVS.String():
		return TokenKeywordVS, buf.String()
	default:
		return TokenIdentifier, buf.String()
	}
}

func (s *Scanner) read() (rune, error) {
	ch, n, err := s.r.ReadRune()
	s.pos += n
	s.last = n
	return ch, err
}

func (s *Scanner) unread() {
	s.r.UnreadRune()
	s.pos -= s.last
	s.last = 0
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }
func isLetter(ch rune) bool     { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }
func isDigit(ch rune) bool      { return (ch >= '0' && ch <= '9') }
