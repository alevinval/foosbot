package foosbot

import "github.com/alevinval/foosbot/parsing"

func (ctx *Context) ExecuteMatch(statement parsing.MatchStatement) error {
	t1Players := []*Player{}
	for _, name := range statement.TeamOne {
		t1Players = append(t1Players, NewPlayer(name))
	}
	team1, err := NewTeam(t1Players...)
	if err != nil {
		return err
	}

	t2Players := []*Player{}
	for _, name := range statement.TeamTwo {
		t2Players = append(t2Players, NewPlayer(name))
	}
	team2, err := NewTeam(t2Players...)
	if err != nil {
		return err
	}
	outcomes := []*Outcome{}
	var i = statement.TeamOneScore
	for i > 0 {
		outcome, err := NewOutcome(team1, team2)
		if err != nil {
			return err
			//return newCommandError(err.Error())
		}
		outcomes = append(outcomes, outcome)
		i--
	}
	i = statement.TeamTwoScore
	for i > 0 {
		outcome, err := NewOutcome(team2, team1)
		if err != nil {
			return err
			//return newCommandError(err.Error())
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

func (ctx *Context) ExecuteStats(statement parsing.StatStatement) (*Stats, error) {
	if len(statement.Names) == 1 {
		player := NewPlayer(statement.Names[0])
		stats := ctx.PlayerStats(player)
		return &stats.Stats, nil

	} else {
		players := []*Player{}
		for _, name := range statement.Names {
			player := NewPlayer(name)
			players = append(players, player)
		}
		team, err := NewTeam(players...)
		if err != nil {
			return nil, err
		}
		stats := ctx.TeamStats(team)
		return &stats.Stats, nil
	}
}
