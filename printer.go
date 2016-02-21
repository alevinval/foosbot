package foosbot

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"strings"
)

func (ctx *Context) ReportStats(status *Stats, obj interface{}) string {
	if status.PlayedGames == 0 {
		return fmt.Sprintf("%s hasn't played any match yet.", ctx.Print(obj))
	}
	response := fmt.Sprintf("*%s*\n", ctx.Print(obj))
	response += fmt.Sprintf("Played %d matches (%d wins - %d defeats) - %.2f%% winrate\n", status.PlayedGames,
		status.Wins, status.Defeats, status.WinRate)
	response += fmt.Sprintf("```Recent match history\n")
	response += ctx.reportHistory(status)
	response += "```"
	return response
}

func (ctx *Context) reportHistory(stats *Stats) string {
	response := ""
	response += string(len(stats.Results))

	for i, result := range stats.Results {
		response += ctx.reportHistoryLine(result)
		if i >= 10 {
			break
		}
	}
	return response
}

func (ctx *Context) reportHistoryLine(result *MatchResult) string {
	return fmt.Sprintf("%s: %-4s against (%s) (%s)\n", result.Match.ShortID(), result.Status,
		ctx.Print(result.Opponent.Players), humanize.Time(result.Match.PlayedAt))
}

func (ctx *Context) Print(i interface{}) (out string) {
	switch obj := i.(type) {
	case *Match:
		out = fmt.Sprintf("%s (%s)", obj.ShortID(), humanize.Time(obj.PlayedAt))
	case *Outcome:
		w, _ := ctx.Query.TeamByID(obj.WinnerID)
		l, _ := ctx.Query.TeamByID(obj.LooserID)
		out = fmt.Sprintf("(%s) vs (%s)", namesFromTeam(w), namesFromTeam(l))
	case *Team:
		out = fmt.Sprintf("%s (%s)", obj.ShortID(), namesFromTeam(obj))
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

func namesFromTeam(t *Team) string {
	names := []string{}
	for _, player := range t.Players {
		names = append(names, player.Name)
	}
	return strings.Join(names, ",")
}

func namesFromPlayers(players []*Player) string {
	names := []string{}
	for _, player := range players {
		names = append(names, player.Name)
	}
	return strings.Join(names, ",")
}
