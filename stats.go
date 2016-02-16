package foosbot

type teamStats struct {
	Team        *Team      `json:"team"`
	Matches     []*Match   `json:"matches"`
	Outcomes    []*Outcome `json:"outcomes"`
	PlayedGames int        `json:"played_games"`
	Wins        int        `json:"wins"`
	Defeats     int        `json:"defeats"`
}

func (ctx *Context) TeamStats(team *Team) *teamStats {
	stats := new(teamStats)
	stats.Team = team
	matches, outcomes := ctx.Query.MatchesWithTeam(team)
	for _, outcome := range outcomes {
		if outcome.IsWinner(team) {
			stats.Wins++
		} else if outcome.IsLooser(team) {
			stats.Defeats++
		} else {
			// That should never happen if repository/behaviour is correct
		}
	}
	stats.PlayedGames = stats.Wins + stats.Defeats
	stats.Outcomes = outcomes
	stats.Matches = matches
	return stats
}
