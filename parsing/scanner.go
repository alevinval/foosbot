package parsing

import (
	"bytes"
	"strings"
)

const (
	// Tokens
	EOF Token = iota
	WHITESPACE
	IDENTIFIER
	DIGIT
	COMMAND_MATCH
	COMMAND_STATS
	KEYWORD_FOOSBOT
	KEYWORD_VS

	eof = rune(0)
)

type (
	Token   int
	Scanner struct {
		r *bytes.Reader
	}
)

func NewScanner(r *bytes.Reader) *Scanner {
	s := new(Scanner)
	s.r = r
	return s
}

func (s *Scanner) read() (rune, error) {
	ch, _, err := s.r.ReadRune()
	return ch, err
}

func (s *Scanner) unread() {
	s.r.UnreadRune()
}

func (s *Scanner) Scan() (Token, string) {
	ch, _ := s.read()
	s.unread()
	if isWhitespace(ch) {
		_, literal := s.scanIdentifier()
		return WHITESPACE, literal
	} else if isDigit(ch) {
		_, literal := s.scanIdentifier()
		return DIGIT, literal
	} else if isLetter(ch) {
		return s.scanIdentifier()
	} else {
		return EOF, "EOF"
	}
}

func (s *Scanner) scanIdentifier() (tok Token, literal string) {
	var buf bytes.Buffer
	ch, err := s.read()
	if err != nil {
		return EOF, "EOF"
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
	switch strings.ToUpper(buf.String()) {
	case "MATCH":
		return COMMAND_MATCH, buf.String()
	case "STATS":
		return COMMAND_STATS, buf.String()
	case "FOOSBOT", "FOOSBALL", "FB":
		return KEYWORD_FOOSBOT, buf.String()
	case "VS":
		return KEYWORD_VS, buf.String()
	default:
		return IDENTIFIER, buf.String()
	}
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }
func isLetter(ch rune) bool     { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }
func isDigit(ch rune) bool      { return (ch >= '0' && ch <= '9') }
