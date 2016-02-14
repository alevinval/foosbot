package foosbot

type teamStats struct {
	Team        *Team           `json:"team"`
	Matches     []*Match        `json:"matches"`
	History     []*HistoryEntry `json:"history"`
	PlayedGames int             `json:"played_games"`
	Wins        int             `json:"wins"`
	Defeats     int             `json:"defeats"`
}

func (c *Context) TeamStats(team *Team) *teamStats {
	ts := new(teamStats)
	ts.Team = team
	matches, history := c.MatchesWithTeam(team)
	for _, match := range matches {
		if match.WinnerID == team.ID {
			ts.Wins++
		} else {
			ts.Defeats++
		}
	}
	ts.PlayedGames = ts.Wins + ts.Defeats
	ts.Matches = matches
	ts.History = history
	return ts
}
