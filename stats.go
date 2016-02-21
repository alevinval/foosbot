package foosbot

type Stats struct {
	Results     []*MatchResult
	PlayedGames int32   `json:"played_games"`
	Wins        int32   `json:"wins"`
	Defeats     int32   `json:"defeats"`
	WinRate     float32 `json:"win_rate"`
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
	stats.Results = Query(ctx).MatchesWithTeam(team)
	for _, result := range stats.Results {
		computeStats(&stats.Stats, result)
	}
	stats.WinRate = float32(stats.Wins) / float32(stats.PlayedGames) * 100
	return stats
}

func (ctx *Context) PlayerStats(player *Player) *playerStats {
	stats := new(playerStats)
	stats.Player = player
	stats.Results = Query(ctx).MatchesWithPlayer(player)
	for _, result := range stats.Results {
		computeStats(&stats.Stats, result)
	}
	stats.WinRate = float32(stats.Wins) / float32(stats.PlayedGames) * 100
	return stats
}

func computeStats(stats *Stats, result *MatchResult) {
	stats.PlayedGames++
	if result.Status == StatusWon {
		stats.Wins++
	} else {
		stats.Defeats++
	}
}
