package foosbot

type (
	MatchResult struct {
		Match    *Match
		Outcome  *Outcome
		Team     *Team
		Opponent *Team
		Winner   bool
	}
	MatchFilterFunc     func(*MatchResult) bool
	MatchesQueryBuilder struct {
		ctx     *Context
		qb      *QueryBuilder
		filter  MatchFilterFunc
		limit   int
		results []*MatchResult
	}
)

var (
	FilterNoMatches MatchFilterFunc = nil
)

func QueryMatches(ctx *Context) *MatchesQueryBuilder {
	return &MatchesQueryBuilder{
		ctx:     ctx,
		qb:      Query(ctx),
		filter:  FilterNoMatches,
		limit:   0,
		results: []*MatchResult{},
	}
}

func (s MatchResult) Status() string {
	if s.Winner {
		return "won"
	} else {
		return "lost"
	}
}

func (mqb *MatchesQueryBuilder) Limit(limit int) *MatchesQueryBuilder {
	mqb.limit = limit
	return mqb
}

func (mqb *MatchesQueryBuilder) FilterByTeam(team *Team) *MatchesQueryBuilder {
	mqb.filter = func(result *MatchResult) bool {
		if team.ID != result.Team.ID && team.ID != result.Opponent.ID {
			return true
		}
		if team.ID == result.Team.ID {
			result.Winner = true
		} else {
			result.Winner = false
			result.Team, result.Opponent = result.Opponent, result.Team
		}
		return false
	}
	return mqb
}

func (mqb *MatchesQueryBuilder) FilterByPlayer(player *Player) *MatchesQueryBuilder {
	playerTeams := mqb.qb.TeamsWithPlayer(player)
	playerTeamsMap := map[string]*Team{}
	for _, team := range playerTeams {
		playerTeamsMap[team.ID] = team
	}
	mqb.filter = func(result *MatchResult) bool {
		_, isWinner := playerTeamsMap[result.Team.ID]
		_, isLooser := playerTeamsMap[result.Opponent.ID]
		if !isWinner && !isLooser {
			return true
		}
		if isWinner {
			result.Winner = true
		} else if isLooser {
			result.Winner = false
			result.Team, result.Opponent = result.Opponent, result.Team
		}
		return false
	}
	return mqb
}

func (mqb *MatchesQueryBuilder) Get() []*MatchResult {
	for i := range mqb.ctx.Matches {
		if mqb.limit > 0 && i >= mqb.limit {
			break
		}
		idx := len(mqb.ctx.Matches) - i - 1
		match := mqb.ctx.Matches[idx]
		outcome, ok := mqb.qb.OutcomeByID(match.OutcomeID)
		if !ok {
			continue
		}
		winner, ok := mqb.qb.TeamByID(outcome.WinnerID)
		if !ok {
			continue
		}
		looser, ok := mqb.qb.TeamByID(outcome.LooserID)
		if !ok {
			continue
		}
		result := &MatchResult{
			Match:    match,
			Outcome:  outcome,
			Team:     winner,
			Opponent: looser,
		}
		if mqb.filter != nil && mqb.filter(result) {
			continue
		}
		mqb.results = append(mqb.results, result)
	}
	return mqb.results
}
