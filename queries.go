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

func (q *Queries) MatchesWithTeam(t *Team) (matches []*Match, outcomes []*Outcome) {
	outcomesMap := map[string]*Outcome{}
	for _, outcome := range q.ctx.Outcomes {
		if outcome.IsWinner(t) || outcome.IsLooser(t) {
			outcomesMap[outcome.ID] = outcome
		}
	}
	for _, match := range q.ctx.Matches {
		outcome, ok := outcomesMap[match.OutcomeID]
		if !ok {
			continue
		}
		outcomes = append(outcomes, outcome)
		matches = append(matches, match)
	}
	return matches, outcomes
}

func (q *Queries) MatchesWithPlayer(p *Player) (matches []*Match, outcomes []*Outcome, teams []*Team) {
	qTeams := q.TeamsWithPlayer(p)
	qTeamsMap := map[string]*Team{}
	for _, team := range qTeams {
		qTeamsMap[team.ID] = team
	}
	outcomesMap := map[string]*Outcome{}
	for _, outcome := range q.ctx.Outcomes {
		_, okW := qTeamsMap[outcome.WinnerID]
		_, okL := qTeamsMap[outcome.LooserID]
		if !okW && !okL {
			continue
		} else {
			outcomesMap[outcome.ID] = outcome
		}
	}
	for _, match := range q.ctx.Matches {
		outcome, ok := outcomesMap[match.OutcomeID]
		if !ok {
			continue
		}
		outcomes = append(outcomes, outcome)
		matches = append(matches, match)
		if w, ok := qTeamsMap[outcome.WinnerID]; ok {
			teams = append(teams, w)
		} else {
			l := qTeamsMap[outcome.LooserID]
			teams = append(teams, l)
		}
	}
	return
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
