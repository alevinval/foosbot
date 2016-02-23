package foosbot

func Query(ctx *Context) *QueryBuilder {
	return &QueryBuilder{ctx: ctx, indexes: ctx.indexes}
}

type QueryBuilder struct {
	ctx     *Context
	indexes *Indexing
}

func (qb *QueryBuilder) OutcomeByID(id string) (*Outcome, bool) {
	m, ok := qb.indexes.outcomesMap[id]
	return m, ok
}

func (qb *QueryBuilder) OutcomeByTeams(a, b *Team) (outcome *Outcome, ok bool) {
	id := BuildOutcomeID(a, b)
	return qb.OutcomeByID(id)
}

func (qb *QueryBuilder) TeamByID(id string) (team *Team, ok bool) {
	team, ok = qb.indexes.teamsMap[id]
	return
}

func (qb *QueryBuilder) TeamByPlayers(players ...*Player) (team *Team, ok bool) {
	id := BuildTeamID(players...)
	team, ok = qb.TeamByID(id)
	return
}

func (qb *QueryBuilder) PlayerByID(playerID string) (player *Player, ok bool) {
	player, ok = qb.indexes.playersMap[playerID]
	return
}

func (qb *QueryBuilder) PlayerByName(name string) (player *Player, ok bool) {
	player, ok = qb.indexes.playersNameMap[name]
	return
}

func (qb *QueryBuilder) TeamsWithPlayer(p *Player) (teams []*Team) {
	return qb.indexes.playersTeamMap[p.ID]
}

func (qb *QueryBuilder) TeamWithPlayer(teams []*Team, player *Player) *Team {
	for _, team := range teams {
		_, ok := qb.indexes.teamsPlayerMap[team.ID][player.ID]
		if ok {
			return team
		}
	}
	return nil
}

func (qb *QueryBuilder) MatchesWithTeam(team *Team) []*MatchResult {
	rs := []*MatchResult{}
	for i := range qb.ctx.Matches {
		match := qb.ctx.Matches[len(qb.ctx.Matches)-i-1]
		outcome, ok := qb.OutcomeByID(match.OutcomeID)
		if !ok {
			continue
		}
		isWinner, isLooser := outcome.IsWinner(team), outcome.IsLooser(team)
		if !isWinner && !isLooser {
			continue
		}
		var opponent *Team
		if isWinner {
			opponent, ok = qb.TeamByID(outcome.LooserID)
		} else {
			opponent, ok = qb.TeamByID(outcome.WinnerID)
		}
		if !ok {
			continue
		}
		result := &MatchResult{
			Match:    match,
			Status:   MatchStatus(isWinner),
			Outcome:  outcome,
			Team:     team,
			Opponent: opponent,
		}
		rs = append(rs, result)
	}
	return rs
}

func (qb *QueryBuilder) MatchesWithPlayer(p *Player) []*MatchResult {
	qTeams := qb.TeamsWithPlayer(p)
	qTeamsMap := map[string]*Team{}
	for _, team := range qTeams {
		qTeamsMap[team.ID] = team
	}

	rs := []*MatchResult{}
	for i := range qb.ctx.Matches {
		match := qb.ctx.Matches[len(qb.ctx.Matches)-i-1]
		outcome, ok := qb.OutcomeByID(match.OutcomeID)
		if !ok {
			continue
		}
		var team, opponent *Team
		winner, isWinner := qTeamsMap[outcome.WinnerID]
		looser, isLooser := qTeamsMap[outcome.LooserID]
		if !isWinner && !isLooser {
			continue
		} else if isWinner {
			team = winner
			opponent, ok = qb.TeamByID(outcome.LooserID)
			if !ok {
				continue
			}
		} else {
			team = looser
			opponent, ok = qb.TeamByID(outcome.WinnerID)
			if !ok {
				continue
			}
		}
		result := &MatchResult{
			Match:    match,
			Status:   MatchStatus(isWinner),
			Outcome:  outcome,
			Team:     team,
			Opponent: opponent,
		}
		rs = append(rs, result)
	}
	return rs
}

func (qb *QueryBuilder) Matches() []*MatchResult {
	rs := []*MatchResult{}
	for i := range qb.ctx.Matches {
		match := qb.ctx.Matches[len(qb.ctx.Matches)-i-1]
		outcome, ok := qb.OutcomeByID(match.OutcomeID)
		if !ok {
			continue
		}
		winner, ok := qb.TeamByID(outcome.WinnerID)
		if !ok {
			continue
		}
		looser, ok := qb.TeamByID(outcome.LooserID)
		if !ok {
			continue
		}
		result := &MatchResult{
			Match:    match,
			Outcome:  outcome,
			Team:     winner,
			Opponent: looser,
		}
		rs = append(rs, result)
	}
	return rs
}
