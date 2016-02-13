package foosbot

type teamStats struct {
	Team        *Team    `json:"team"`
	Matches     []*Match `json:"matches"`
	PlayedGames int      `json:"played_games"`
	Wins        int      `json:"wins"`
	Defeats     int      `json:"defeats"`
}

func TeamStats(team *Team) *teamStats {
	ts := new(teamStats)
	ts.Team = team
	matches := MatchesWithTeam(team)
	for _, match := range matches {
		if match.WinnerID == team.ID {
			ts.Wins += match.N
		} else {
			ts.Defeats += match.N
		}
	}
	ts.PlayedGames = ts.Wins + ts.Defeats
	ts.Matches = matches
	return ts
}
