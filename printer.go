package foosbot

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"strings"
)

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

func (ctx *Context) ReportHistoryLine(match *Match, outcome *Outcome, team *Team, opponent *Team) string {
	outcomeStr := "won"
	if outcome.IsLooser(team) {
		outcomeStr = "lost"
	}
	return fmt.Sprintf("%s: %-4s against (%s) (%s)\n", match.ShortID(), outcomeStr, ctx.Print(opponent.Players),
		humanize.Time(match.PlayedAt))
}

func (ctx *Context) ReportStats(s Stats, obj interface{}) string {
	if s.PlayedGames == 0 {
		return fmt.Sprintf("%s hasn't played any match yet.", ctx.Print(obj))
	}
	response := fmt.Sprintf("*%s*\n", ctx.Print(obj))
	response += fmt.Sprintf("Played %d matches (%d wins - %d defeats) - %.2f%% winrate\n", s.PlayedGames,
		s.Wins, s.Defeats, s.WinRate)
	return response
}

func (ctx *Context) ReportTeamHistory(s Stats, team *Team) string {
	report := ""
	for i := range s.Matches {
		idx := len(s.Outcomes) - 1 - i
		outcome := s.Outcomes[idx]
		match := s.Matches[idx]
		winner, _ := ctx.Query.TeamByID(outcome.WinnerID)
		looser, _ := ctx.Query.TeamByID(outcome.LooserID)

		opponent := looser
		if outcome.IsLooser(team) {
			opponent = winner
		}
		report += ctx.ReportHistoryLine(match, outcome, team, opponent)
		if i >= 10 {
			break
		}
	}
	return report
}

func (ctx *Context) ReportPlayerHistory(s Stats, player *Player) string {
	report := ""
	for i := range s.Matches {
		idx := len(s.Outcomes) - 1 - i
		outcome := s.Outcomes[idx]
		match := s.Matches[idx]

		winner, _ := ctx.Query.TeamByID(outcome.WinnerID)
		looser, _ := ctx.Query.TeamByID(outcome.LooserID)
		team := ctx.Query.TeamWithPlayer([]*Team{winner, looser}, player)
		opponent := winner
		if outcome.IsWinner(team) {
			opponent = looser
		}
		report += ctx.ReportHistoryLine(match, outcome, team, opponent)
		if i >= 10 {
			break
		}
	}
	return report
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
