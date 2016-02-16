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
	outcomeIds := []string{}
	outcomesMap := map[string]*Outcome{}
	for _, outcome := range q.ctx.Outcomes {
		if outcome.IsLooser(t) || outcome.IsWinner(t) {
			outcomeIds = append(outcomeIds, outcome.ID)
			outcomesMap[outcome.ID] = outcome
		}
	}
	for _, entry := range q.ctx.Matches {
		outcome, ok := outcomesMap[entry.OutcomeID]
		if !ok {
			continue
		}
		outcomes = append(outcomes, outcome)
		matches = append(matches, entry)
	}
	return matches, outcomes
}
