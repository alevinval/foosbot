package parsing

import (
	"github.com/alevinval/foosbot"
)

type (
	MatchStatement struct {
		TeamOne      []string
		TeamOneScore int
		TeamTwo      []string
		TeamTwoScore int
		TotalMatches int
	}
)

func (m *MatchStatement) Execute(ctx *foosbot.Context) error {
	t1Players := []*foosbot.Player{}
	for _, name := range m.TeamOne {
		t1Players = append(t1Players, foosbot.NewPlayer(name))
	}
	team1, err := foosbot.NewTeam(t1Players...)
	if err != nil {
		return err
	}

	t2Players := []*foosbot.Player{}
	for _, name := range m.TeamTwo {
		t2Players = append(t2Players, foosbot.NewPlayer(name))
	}
	team2, err := foosbot.NewTeam(t2Players...)
	if err != nil {
		return err
	}
	outcomes := []*foosbot.Outcome{}
	var i = m.TeamOneScore
	for i > 0 {
		outcome, err := foosbot.NewOutcome(team1, team2)
		if err != nil {
			return newCommandError(err.Error())
		}
		outcomes = append(outcomes, outcome)
		i--
	}
	i = m.TeamTwoScore
	for i > 0 {
		outcome, err := foosbot.NewOutcome(team2, team1)
		if err != nil {
			return newCommandError(err.Error())
		}
		outcomes = append(outcomes, outcome)
		i--
	}
	ctx.AddTeam(team1)
	ctx.AddTeam(team2)
	for i := range outcomes {
		ctx.AddMatchWithOutcome(outcomes[i])
	}
	return nil
}
