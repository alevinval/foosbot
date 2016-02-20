package foosbot

type Stats struct {
	Matches     []*Match   `json:"matches"`
	Outcomes    []*Outcome `json:"outcomes"`
	PlayedGames int32      `json:"played_games"`
	Wins        int32      `json:"wins"`
	Defeats     int32      `json:"defeats"`
	WinRate     float32    `json:"win_rate"`
}

type teamStats struct {
	Stats
	Team *Team `json:"team"`
}

type playerStats struct {
	Stats
	Player *Player `json:"player"`
}

func (ctx *Context) TeamStats(team *Team) *teamStats {
	stats := new(teamStats)
	stats.Team = team
	matches, outcomes := ctx.Query.MatchesWithTeam(team)
	for i := range matches {
		computeStats(&stats.Stats, team, matches[i], outcomes[i])
	}
	stats.WinRate = float32(stats.Wins) / float32(stats.PlayedGames) * 100
	return stats
}

func (ctx *Context) PlayerStats(player *Player) *playerStats {
	stats := new(playerStats)
	stats.Player = player
	matches, outcomes, teams := ctx.Query.MatchesWithPlayer(player)
	for i := range matches {
		computeStats(&stats.Stats, teams[i], matches[i], outcomes[i])
	}
	stats.WinRate = float32(stats.Wins) / float32(stats.PlayedGames) * 100
	return stats
}

func computeStats(s *Stats, team *Team, match *Match, outcome *Outcome) {
	s.Matches = append(s.Matches, match)
	s.Outcomes = append(s.Outcomes, outcome)
	if outcome.IsWinner(team) {
		s.Wins++
		s.PlayedGames++
	} else if outcome.IsLooser(team) {
		s.Defeats++
		s.PlayedGames++
	} else {
		// Ignore history entry then
	}
}
