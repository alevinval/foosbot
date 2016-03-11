package foosbot

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"strings"
)

func (ctx *Context) ReportStats(stats *Stats) string {
	if stats.PlayedGames == 0 {
		return fmt.Sprintf("%s hasn't played any match yet.", Print(stats.who))
	}
	response := fmt.Sprintf("*%s*\n", Print(stats.who))
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
		Print(result.Opponent.Players), humanize.Time(result.Match.PlayedAt))
}

func Print(i interface{}) (out string) {
	switch obj := i.(type) {
	case *Match:
		out = fmt.Sprintf("%s (%s)", obj.ShortID(), humanize.Time(obj.PlayedAt))
	case *Team:
		out = fmt.Sprintf("%s (%s)", obj.ShortID(), namesFromPlayers(obj.Players))
	case *Player:
		out = fmt.Sprintf("%s (%s)", obj.ShortID(), obj.Name)
	case []*Player:
		out = namesFromPlayers(obj)
	case *MatchResult:
		pad := func(field, value interface{}, pad int) string {
			f := fmt.Sprintf("%%%ds: %%s", pad)
			return fmt.Sprintf(f, field, value)
		}
		match := pad("Match", obj.Match.ID, 8)
		date := pad("Date", obj.Match.PlayedAt, 8)
		winner := pad("Winner", Print(obj.Team), 8)
		looser := pad("Looser", Print(obj.Opponent), 8)
		out = fmt.Sprintf("\n%s\n%s\n%s\n%s\n", match, date, winner, looser)
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
