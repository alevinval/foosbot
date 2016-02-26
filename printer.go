package foosbot

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"strings"
)

func (ctx *Context) ReportStats(stats *Stats) string {
	if stats.PlayedGames == 0 {
		return fmt.Sprintf("%s hasn't played any match yet.", ctx.Print(stats.who))
	}
	response := fmt.Sprintf("*%s*\n", ctx.Print(stats.who))
	response += fmt.Sprintf("Played %d matches (%d wins - %d defeats) - %.2f%% winrate\n", stats.PlayedGames,
		stats.Wins, stats.Defeats, stats.WinRate)
	response += fmt.Sprintf("```Recent match history:\n")
	response += ctx.reportStatsHistory(stats)
	response += "```"
	return response
}

func (ctx *Context) ReportLeaderBoard(stats PlayerStatsSlice) string {
	response := "Top leaderboard:\n```"
	for i, stat := range stats {
		response += fmt.Sprintf("%d.- %-12s w: %-3d l: %-3d (%-5.2f%%)\n", i+1, stat.Player.Name, stat.Wins,
			stat.Defeats, stat.WinRate)
		if i >= 10 {
			break
		}
	}
	response += "```"
	return response
}

func (ctx *Context) reportStatsHistory(stats *Stats) string {
	response := ""
	for i, result := range stats.Matches {
		response += ctx.reportMatchResult(result)
		if i >= 10 {
			break
		}
	}
	return response
}

func (ctx *Context) reportMatchResult(result *MatchResult) string {
	return fmt.Sprintf("%s: %-4s against (%s) (%s)\n", result.Match.ShortID(), result.Status(),
		ctx.Print(result.Opponent.Players), humanize.Time(result.Match.PlayedAt))
}

func (ctx *Context) Print(i interface{}) (out string) {
	switch obj := i.(type) {
	case *Match:
		out = fmt.Sprintf("%s (%s)", obj.ShortID(), humanize.Time(obj.PlayedAt))
	case *Team:
		out = fmt.Sprintf("%s (%s)", obj.ShortID(), namesFromPlayers(obj.Players))
	case *Player:
		out = fmt.Sprintf("%s (%s)", obj.ShortID(), obj.Name)
	case []*Player:
		out = namesFromPlayers(obj)
	default:
		b, _ := json.Marshal(obj)
		out = string(b)
	}
	return
}

func namesFromPlayers(players []*Player) string {
	names := []string{}
	for _, player := range players {
		names = append(names, player.Name)
	}
	return strings.Join(names, ",")
}
