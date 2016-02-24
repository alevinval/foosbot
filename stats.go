package foosbot

import (
	"sort"
)

type Stats struct {
	Matches     []*MatchResult
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

type playerStatsSlice []*playerStats

func (p playerStatsSlice) Len() int {
	return len(p)
}
func (p playerStatsSlice) Less(i, j int) bool {
	// Reverse order ( we want higher win-rates first )
	return p[i].WinRate > p[j].WinRate
}
func (p playerStatsSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (ctx *Context) TeamStats(team *Team) *teamStats {
	stats := new(teamStats)
	stats.Team = team
	stats.Matches = Query(ctx).FilterByTeam(team).Matches().Results()
	for _, result := range stats.Matches {
		computeStats(&stats.Stats, result)
	}
	stats.WinRate = float32(stats.Wins) / float32(stats.PlayedGames) * 100
	return stats
}

func (ctx *Context) PlayerStats(player *Player) *playerStats {
	stats := new(playerStats)
	stats.Player = player
	stats.Matches = Query(ctx).FilterByPlayer(player).Matches().Results()
	for _, result := range stats.Matches {
		computeStats(&stats.Stats, result)
	}
	stats.WinRate = float32(stats.Wins) / float32(stats.PlayedGames) * 100
	return stats
}

func (ctx *Context) PlayersStatsFromMatches() playerStatsSlice {
	statsMap := make(map[string]*playerStats)
	for _, player := range ctx.Players {
		statsMap[player.Name] = &playerStats{Player: player}
	}
	matches := Query(ctx).Matches().Results()
	for _, match := range matches {
		winnerPlayers := match.Team.Players
		looserPlayers := match.Opponent.Players
		for _, player := range winnerPlayers {
			statsMap[player.Name].Wins++
			statsMap[player.Name].PlayedGames++
		}
		for _, player := range looserPlayers {
			statsMap[player.Name].Defeats++
			statsMap[player.Name].PlayedGames++
		}
	}
	for _, s := range statsMap {
		s.WinRate = float32(s.Wins) / float32(s.PlayedGames) * 100
	}

	stats := playerStatsSlice{}
	for _, s := range statsMap {
		stats = append(stats, s)
	}
	sort.Sort(stats)
	var reverseStats playerStatsSlice
	reverseStats = append(reverseStats, stats...)
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
