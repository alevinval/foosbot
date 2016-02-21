package foosbot

type Queries struct {
	ctx *Context
}

func (q *Queries) OutcomeByID(id string) (*Outcome, bool) {
	m, ok := q.ctx.outcomesMap[id]
	return m, ok
}

func (q *Queries) OutcomeByTeams(a, b *Team) (outcome *Outcome, ok bool) {
	id := BuildOutcomeID(a, b)
	return q.OutcomeByID(id)
}

func (q *Queries) TeamByID(id string) (team *Team, ok bool) {
	team, ok = q.ctx.teamsMap[id]
	return
}

func (q *Queries) TeamByPlayers(players ...*Player) (team *Team, ok bool) {
	id := BuildTeamID(players...)
	team, ok = q.TeamByID(id)
	return
}

func (q *Queries) PlayerByID(playerID string) (player *Player, ok bool) {
	player, ok = q.ctx.playersMap[playerID]
	return
}

func (q *Queries) PlayerByName(name string) (player *Player, ok bool) {
	player, ok = q.ctx.playersNameMap[name]
	return
}

func (q *Queries) MatchesWithTeam(team *Team) []*MatchResult {
	rs := []*MatchResult{}
	for i := range q.ctx.Matches {
		match := q.ctx.Matches[len(q.ctx.Matches)-i-1]
		outcome, ok := q.OutcomeByID(match.OutcomeID)
		if !ok {
			continue
		}
		isWinner, isLooser := outcome.IsWinner(team), outcome.IsLooser(team)
		if !isWinner && !isLooser {
			continue
		}
		var opponent *Team
		if isWinner {
			opponent, ok = q.TeamByID(outcome.LooserID)
		} else {
			opponent, ok = q.TeamByID(outcome.WinnerID)
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

func (q *Queries) MatchesWithPlayer(p *Player) []*MatchResult {
	qTeams := q.TeamsWithPlayer(p)
	qTeamsMap := map[string]*Team{}
	for _, team := range qTeams {
		qTeamsMap[team.ID] = team
	}

	rs := []*MatchResult{}
	for i := range q.ctx.Matches {
		match := q.ctx.Matches[len(q.ctx.Matches)-i-1]
		outcome, ok := q.OutcomeByID(match.OutcomeID)
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
			opponent, ok = q.TeamByID(outcome.LooserID)
			if !ok {
				continue
			}
		} else {
			team = looser
			opponent, ok = q.TeamByID(outcome.WinnerID)
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

func (q *Queries) TeamsWithPlayer(p *Player) (teams []*Team) {
	return q.ctx.playersTeamMap[p.ID]
}

func (q *Queries) TeamWithPlayer(teams []*Team, player *Player) *Team {
	for _, team := range teams {
		_, ok := q.ctx.teamsPlayerMap[team.ID][player.ID]
		if ok {
			return team
		}
	}
	return nil
}
