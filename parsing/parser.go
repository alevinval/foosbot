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
	s            *Scanner
	lastToken    Token
	useLastToken bool
}

func NewParser(r *bytes.Reader) *Parser {
	p := new(Parser)
	p.s = NewScanner(r)
	return p
}

func (p *Parser) ParseCommand() (Token, error) {
	token := p.scan()
	if token.Type != TokenKeywordFoosbot {
		return token, ErrNotFoosbotCommand
	}
	token = p.scan()
	if token.Type != TokenCommandMatch && token.Type != TokenCommandStats {
		return token, newParseError(token, "a valid command (match, stats)")
	}
	return token, nil
}

func (p *Parser) ParseMatch() (outcomes []*foosbot.Outcome, teams []*foosbot.Team, err error) {
	var t1Score, t2Score int64
	p1name, err := p.parsePlayerName()
	if err != nil {
		return
	}
	p2name, err := p.parsePlayerName()
	if err != nil {
		return
	}
	t1Score, err = p.parseScore()
	if err != nil {
		return
	}
	err = p.parseVs()
	if err != nil {
		return
	}
	t2Score, err = p.parseScore()
	if err != nil {
		return
	}
	p3name, err := p.parsePlayerName()
	if err != nil {
		return
	}
	p4name, err := p.parsePlayerName()
	if err != nil {
		return
	}
	err = p.parseEof()
	if err != nil {
		return
	}
	t1Players := []string{p1name, p2name}
	t2Players := []string{p3name, p4name}

	// Parsing correct, create outcomes
	p1, p2 := foosbot.NewPlayer(t1Players[0]), foosbot.NewPlayer(t1Players[1])
	p3, p4 := foosbot.NewPlayer(t2Players[0]), foosbot.NewPlayer(t2Players[1])
	t1, perr := foosbot.NewTeam(p1, p2)
	if perr != nil {
		err = newCommandError(perr.Error())
		return
	}
	t2, perr := foosbot.NewTeam(p3, p4)
	if perr != nil {
		err = newCommandError(perr.Error())
		return
	}

	for t1Score > 0 {
		outcome, oerr := foosbot.NewOutcome(t1, t2)
		if oerr != nil {
			err = newCommandError(oerr.Error())
			return
		}
		outcomes = append(outcomes, outcome)
		t1Score--
	}
	for t2Score > 0 {
		outcome, oerr := foosbot.NewOutcome(t2, t1)
		if err != nil {
			err = newCommandError(oerr.Error())
			return
		}
		outcomes = append(outcomes, outcome)
		t2Score--
	}
	return outcomes, []*foosbot.Team{t1, t2}, nil
}

func (p *Parser) ParseStats() (interface{}, error) {
	p1name, err := p.parsePlayerName()
	if err != nil {
		return nil, err
	}
	err = p.parseEof()
	if err != nil {
		p.unscan()
	} else {
		p1 := foosbot.NewPlayer(p1name)
		return p1, nil
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
	team, _ := foosbot.NewTeam(p1, p2)
	return team, nil
}

func (p *Parser) parsePlayerName() (string, error) {
	token := p.scan()
	if token.Type != TokenIdentifier {
		return "", newParseError(token, fmt.Sprintf("player name %s", TokenIdentifier))
	}
	return token.Literal, nil
}

func (p *Parser) parseScore() (int64, error) {
	token := p.scan()
	if token.Type != TokenDigit {
		return 0, newParseError(token, fmt.Sprintf("team score %s", TokenDigit))
	}
	value, _ := strconv.ParseInt(token.Literal, 10, 0)
	return value, nil
}

func (p *Parser) parseVs() error {
	token := p.scan()
	if token.Type != TokenKeywordVS {
		return newParseError(token, fmt.Sprintf("%s keyword", TokenKeywordVS))
	}
	return nil
}
func (p *Parser) parseEof() error {
	token := p.scan()
	if token.Type != TokenEOF {
		return newParseError(token, TokenEOF.String())
	}
	return nil
}

func (p *Parser) scan() Token {
	if p.useLastToken {
		p.useLastToken = false
		return p.lastToken
	}

	token := p.s.Scan()
	if token.Type == TokenWhitespace {
		token = p.s.Scan()
	}
	p.lastToken = token
	return token
}

func (p *Parser) unscan() {
	p.useLastToken = true
}

func newParseError(found Token, expected string) error {
	var msg string
	if found.Literal == "" {
		msg = fmt.Sprintf("Syntax error: found %s, expected %s on position %d.", found.Type, expected, found.Pos)
	} else {
		msg = fmt.Sprintf("Syntax error: found %s %q, expected %s on position %d.", found.Type, found.Literal, expected, found.Pos)
	}
	return errors.New(msg)
}

func newCommandError(message string) error {
	msg := fmt.Sprintf("Invalid command: %s.", message)
	return errors.New(msg)
}
