package foosbot

import (
	"sort"
)

type (
	Stats struct {
		PlayedGames int32   `json:"played_games"`
		Wins        int32   `json:"wins"`
		Defeats     int32   `json:"defeats"`
		WinRate     float32 `json:"win_rate"`
		Matches     []*MatchResult
		who         interface{}
	}
	TeamStats struct {
		Stats
		Team *Team `json:"team"`
	}
	PlayerStats struct {
		Stats
		Player *Player `json:"player"`
	}
	PlayerStatsSlice []*PlayerStats
)

func (p PlayerStatsSlice) Len() int {
	return len(p)
}

func (p PlayerStatsSlice) Less(i, j int) bool {
	return p[i].WinRate > p[j].WinRate // Reverse order ( we want higher win-rates first )
}

func (p PlayerStatsSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (ctx *Context) TeamStats(team *Team) *TeamStats {
	teamstats := new(TeamStats)
	teamstats.who = team
	teamstats.Team = team
	teamstats.Matches = QueryMatches(ctx).FilterByTeam(team).Get()
	for _, result := range teamstats.Matches {
		computeStats(&teamstats.Stats, result)
	}
	computeWinrate(&teamstats.Stats)
	return teamstats
}

func (ctx *Context) PlayerStats(player *Player) *PlayerStats {
	playerstats := new(PlayerStats)
	playerstats.who = player
	playerstats.Player = player
	playerstats.Matches = QueryMatches(ctx).FilterByPlayer(player).Get()
	for _, result := range playerstats.Matches {
		computeStats(&playerstats.Stats, result)
	}
	computeWinrate(&playerstats.Stats)
	return playerstats
}

func (ctx *Context) PlayersStatsFromMatches(minPlayedGames int32, maxResults int) PlayerStatsSlice {
	playerstatsMap := make(map[string]*PlayerStats)
	for _, player := range ctx.Players {
		playerstatsMap[player.Name] = &PlayerStats{Player: player}
	}
	matches := QueryMatches(ctx).Get()
	for _, match := range matches {
		winnerPlayers := match.Team.Players
		looserPlayers := match.Opponent.Players
		for _, player := range winnerPlayers {
			playerstatsMap[player.Name].Wins++
			playerstatsMap[player.Name].PlayedGames++
		}
		for _, player := range looserPlayers {
			playerstatsMap[player.Name].Defeats++
			playerstatsMap[player.Name].PlayedGames++
		}
	}
	for _, playerstats := range playerstatsMap {
		computeWinrate(&playerstats.Stats)
	}

	stats := PlayerStatsSlice{}
	for _, s := range playerstatsMap {
		if s.PlayedGames >= minPlayedGames {
			stats = append(stats, s)
		}
	}
	sort.Sort(stats)

	// Trim the stats
	numResults := len(stats)
	if numResults > maxResults {
		numResults = maxResults
	}
	return stats[:numResults]
}

func computeStats(stats *Stats, result *MatchResult) {
	stats.PlayedGames++
	if result.Winner {
		stats.Wins++
	} else {
		stats.Defeats++
	}
}

func computeWinrate(stats *Stats) {
	stats.WinRate = float32(stats.Wins) / float32(stats.PlayedGames) * 100
}
