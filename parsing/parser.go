package parsing

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/alevinval/foosbot"
	"strconv"
)

var (
	ErrNotFoosbotCommand = errors.New("invalid foosbot command")
)

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token
		lit string
		n   int
	}
}

func NewParser(r *bytes.Reader) *Parser {
	p := new(Parser)
	p.s = NewScanner(r)
	return p
}

func (p *Parser) ParseCommand() (Token, error) {
	token, _ := p.scan()
	if token != TokenKeywordFoosbot {
		return token, ErrNotFoosbotCommand
	}
	token, literal := p.scan()
	if token != TokenCommandMatch && token != TokenCommandStats {
		return token, newParseError(token.String(), "a valid command (match, stats)", literal)
	}
	return token, nil
}

func (p *Parser) ParseMatch() ([]*foosbot.Match, error) {
	var t1Score, t2Score int64
	p1name, err := p.parsePlayerName()
	if err != nil {
		return nil, err
	}
	p2name, err := p.parsePlayerName()
	if err != nil {
		return nil, err
	}
	t1Score, err = p.parseScore()
	if err != nil {
		return nil, err
	}
	err = p.parseVs()
	if err != nil {
		return nil, err
	}
	t2Score, err = p.parseScore()
	if err != nil {
		return nil, err
	}
	p3name, err := p.parsePlayerName()
	if err != nil {
		return nil, err
	}
	p4name, err := p.parsePlayerName()
	if err != nil {
		return nil, err
	}
	err = p.parseEof()
	if err != nil {
		return nil, err
	}
	t1Players := []string{p1name, p2name}
	t2Players := []string{p3name, p4name}

	if t1Players[0] == t1Players[1] {
		return nil, newCommandError(fmt.Sprintf("player %q found twice in team 1", t1Players[0]))
	}
	if t2Players[0] == t2Players[1] {
		return nil, newCommandError(fmt.Sprintf("player %q found twice in team 2", t2Players[0]))
	}
	if in(t1Players[0], t2Players) {
		return nil, newCommandError(fmt.Sprintf("player %q cannot be in both teams", t1Players[0]))
	}
	if in(t1Players[1], t2Players) {
		return nil, newCommandError(fmt.Sprintf("player %q cannot be in both teams", t1Players[1]))
	}

	// Parsing correct, re-create match history
	p1, p2 := foosbot.NewPlayer(t1Players[0]), foosbot.NewPlayer(t1Players[1])
	t1 := foosbot.NewTeam(p1, p2)
	p3, p4 := foosbot.NewPlayer(t2Players[0]), foosbot.NewPlayer(t2Players[1])
	t2 := foosbot.NewTeam(p3, p4)

	matches := []*foosbot.Match{}
	for t1Score > 0 {
		match := foosbot.NewMatch(t1, t2)
		matches = append(matches, match)
		t1Score--
	}
	for t2Score > 0 {
		match := foosbot.NewMatch(t2, t1)
		matches = append(matches, match)
		t2Score--
	}
	return matches, nil
}

func (p *Parser) ParseStats() (*foosbot.Team, error) {
	p1name, err := p.parsePlayerName()
	if err != nil {
		return nil, err
	}
	p2name, err := p.parsePlayerName()
	if err != nil {
		return nil, err
	}
	err = p.parseEof()
	if err != nil {
		return nil, err
	}
	p1 := foosbot.NewPlayer(p1name)
	p2 := foosbot.NewPlayer(p2name)
	team := foosbot.NewTeam(p1, p2)
	return team, nil
}

func (p *Parser) parsePlayerName() (string, error) {
	token, literal := p.scan()
	if token != TokenIdentifier {
		return "", newParseError(token.String(), fmt.Sprintf("player name %s", TokenIdentifier), literal)
	}
	return literal, nil
}

func (p *Parser) parseScore() (int64, error) {
	token, literal := p.scan()
	if token != TokenDigit {
		return 0, newParseError(token.String(), fmt.Sprintf("team score %s", TokenDigit), literal)
	}
	value, _ := strconv.ParseInt(literal, 10, 0)
	return value, nil
}

func (p *Parser) parseVs() error {
	token, literal := p.scan()
	if token != TokenKeywordVS {
		return newParseError(token.String(), fmt.Sprintf("%s keyword", TokenKeywordVS), literal)
	}
	return nil
}
func (p *Parser) parseEof() error {
	token, literal := p.scan()
	if token != TokenEOF {
		return newParseError(token.String(), TokenEOF.String(), literal)
	}
	return nil
}

func (p *Parser) scan() (tok Token, lit string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	tok, lit = p.s.Scan()
	if tok == TokenWhitespace {
		tok, lit = p.s.Scan()
	}

	p.buf.tok, p.buf.lit = tok, lit
	return tok, lit
}

func (p *Parser) unscan() {
	p.buf.n = 1
}

func in(match string, arr []string) bool {
	for _, el := range arr {
		if el == match {
			return true
		}
	}
	return false
}

func newParseError(found, expected, literal string) error {
	var msg string
	if literal == "" {
		msg = fmt.Sprintf("Syntax error: found %s, expected %s.", found, expected)
	} else {
		msg = fmt.Sprintf("Syntax error: found %s %s, expected %s.", literal, found, expected)
	}
	return errors.New(msg)
}

func newCommandError(message string) error {
	msg := fmt.Sprintf("Invalid command: %s", message)
	return errors.New(msg)
}
