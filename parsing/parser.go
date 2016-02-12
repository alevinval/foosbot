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
	if token != KEYWORD_FOOSBOT {
		return token, ErrNotFoosbotCommand
	}
	token, literal := p.scan()
	p.unscan()
	if token != COMMAND_MATCH && token != COMMAND_STATS {
		msg := fmt.Sprintf("Invalid command: %q is not a valid command, try with: %q, %q", literal, "match", "stats")
		return token, errors.New(msg)
	}
	return token, nil
}

func (p *Parser) ParseMatch() ([]*foosbot.Match, error) {
	var t1Players, t2Players []string
	var t1Score, t2Score int64

	// Expect match keyword
	tok, literal := p.scan()
	if tok != COMMAND_MATCH {
		msg := fmt.Sprintf("Invalid command: expected %q keyword, instead found %q", "match", literal)
		return nil, errors.New(msg)
	}

	// Expect T1 identifiers
	tok, literal = p.scan()
	for tok == IDENTIFIER {
		t1Players = append(t1Players, literal)
		tok, literal = p.scan()
	}
	if len(t1Players) != 2 {
		msg := fmt.Sprintf("Invalid command: provide 2 opponent names for team 1")
		return nil, errors.New(msg)
	}

	// Expect T1 score
	if tok == DIGIT {
		t1Score, _ = strconv.ParseInt(literal, 10, 0)
	} else {
		msg := fmt.Sprintf("Invalid command: expected team 1 score, instead found %q", literal)
		return nil, errors.New(msg)
	}

	// Expect VS keyword
	tok, literal = p.scan()
	if tok != KEYWORD_VS {
		msg := fmt.Sprintf("Invalid command: expected %q keyword, instead found %q", "vs", literal)
		return nil, errors.New(msg)
	}

	// Expect T2 score
	tok, literal = p.scan()
	if tok == DIGIT {
		t2Score, _ = strconv.ParseInt(literal, 10, 0)
	} else {
		msg := fmt.Sprintf("Invalid command: expected team 2 score, instead found %q", literal)
		return nil, errors.New(msg)
	}

	// Expect T2 identifier
	tok, literal = p.scan()
	for tok == IDENTIFIER {
		t2Players = append(t2Players, literal)
		tok, literal = p.scan()
	}
	if len(t2Players) != 2 {
		msg := fmt.Sprintf("Invalid command: provide 2 opponent names for team 2")
		return nil, errors.New(msg)
	}

	// Expect end of command
	tok, literal = p.scan()
	if tok != EOF {
		msg := fmt.Sprintf("Invalid command: expected end of command, instead found %q", literal)
		return nil, errors.New(msg)
	}

	if t1Players[0] == t1Players[1] {
		msg := fmt.Sprintf("Invalid command: player %q found twice in team 1", t1Players[0])
		return nil, errors.New(msg)
	}
	if t2Players[0] == t2Players[1] {
		msg := fmt.Sprintf("Invalid command: player %q found twice in team 2", t1Players[0])
		return nil, errors.New(msg)
	}
	if in(t1Players[0], t2Players) {
		msg := fmt.Sprintf("Invalid command: player %q cannot be in both teams", t1Players[0])
		return nil, errors.New(msg)
	}
	if in(t1Players[1], t2Players) {
		msg := fmt.Sprintf("Invalid command: player %q cannot be in both teams", t1Players[1])
		return nil, errors.New(msg)
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

func (p *Parser) ParseStats() ([]string, error) {
	var p1, p2 string

	tok, literal := p.scan()
	if tok != COMMAND_STATS {
		msg := fmt.Sprintf("Invalid command: expected %q keyword, instead found %q", "stats", literal)
		return nil, errors.New(msg)
	}

	tok, literal = p.scan()
	if tok != IDENTIFIER {
		msg := fmt.Sprintf("Invalid command: expected 1st player name, instead found %q", literal)
		return nil, errors.New(msg)
	}
	p1 = literal

	tok, literal = p.scan()
	if tok != IDENTIFIER {
		msg := fmt.Sprintf("Invalid command: expected 2nd player name, instead found %q", literal)
		return nil, errors.New(msg)
	}
	p2 = literal

	tok, literal = p.scan()
	if tok != EOF {
		msg := fmt.Sprintf("Invalid command: expected EOF, instead found %q", literal)
		return nil, errors.New(msg)
	}

	return []string{p1, p2}, nil
}

func (p *Parser) scan() (tok Token, lit string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	tok, lit = p.s.Scan()
	if tok == WHITESPACE {
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
