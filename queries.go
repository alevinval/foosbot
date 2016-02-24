package foosbot

type (
	QueryBuilder struct {
		ctx     *Context
		indexes *Indexing
		results []*MatchResult
		filter  FilterFunc
	}

	FilterFunc func(result *MatchResult) bool
)

var (
	FilterNothing FilterFunc = nil
)

func Query(ctx *Context) *QueryBuilder {
	return &QueryBuilder{
		ctx:     ctx,
		indexes: ctx.indexes,
		filter:  FilterNothing,
		results: []*MatchResult{},
	}
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

func (qb *QueryBuilder) Matches() *QueryBuilder {
	for i := range qb.ctx.Matches {
		idx := len(qb.ctx.Matches) - i - 1
		match := qb.ctx.Matches[idx]
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
		if qb.filter != nil && qb.filter(result) {
			continue
		}
		qb.results = append(qb.results, result)
	}
	return qb
}

func (qb *QueryBuilder) FilterByTeam(team *Team) *QueryBuilder {
	qb.filter = func(result *MatchResult) bool {
		if team.ID != result.Team.ID && team.ID != result.Opponent.ID {
			return true
		}
		if team.ID == result.Team.ID {
			result.Status = StatusWon
		} else {
			result.Status = StatusLost
			swap(result.Team, result.Opponent)
		}
		return false
	}
	return qb
}

func (qb *QueryBuilder) FilterByPlayer(player *Player) *QueryBuilder {
	playerTeams := qb.TeamsWithPlayer(player)
	playerTeamsMap := map[string]*Team{}
	for _, team := range playerTeams {
		playerTeamsMap[team.ID] = team
	}
	qb.filter = func(result *MatchResult) bool {
		_, isWinner := playerTeamsMap[result.Team.ID]
		_, isLooser := playerTeamsMap[result.Opponent.ID]
		if !isWinner && !isLooser {
			return true
		}
		if isWinner {
			result.Status = StatusWon
		} else if isLooser {
			result.Status = StatusLost
			swap(result.Team, result.Opponent)
		}
		return false
	}
	return qb
}

func (qb *QueryBuilder) Results() []*MatchResult {
	return qb.results
}

func swap(a, b interface{}) {
	a, b = b, a
}
