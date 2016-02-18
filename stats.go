package foosbot

type stats struct {
	Matches     []*Match   `json:"matches"`
	Outcomes    []*Outcome `json:"outcomes"`
	PlayedGames int        `json:"played_games"`
	Wins        int        `json:"wins"`
	Defeats     int        `json:"defeats"`
}

type teamStats struct {
	stats
	Team *Team `json:"team"`
}

type playerStats struct {
	stats
	Player *Player `json:"player"`
}

func (ctx *Context) TeamStats(team *Team) *teamStats {
	stats := new(teamStats)
	stats.Team = team
	matches, outcomes := ctx.Query.MatchesWithTeam(team)
	for i := range matches {
		computeStats(&stats.stats, team, matches[i], outcomes[i])
	}
	return stats
}

func (ctx *Context) PlayerStats(player *Player) *playerStats {
	stats := new(playerStats)
	stats.Player = player
	matches, outcomes, teams := ctx.Query.MatchesWithPlayer(player)
	for i := range matches {
		computeStats(&stats.stats, teams[i], matches[i], outcomes[i])
	}
	return stats
}

func computeStats(s *stats, team *Team, match *Match, outcome *Outcome) {
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
