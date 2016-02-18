package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/alevinval/foosbot"
	"github.com/alevinval/foosbot/parsing"
	"github.com/dustin/go-humanize"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

func addMatchCommand(ctx *foosbot.Context, outcomes []*foosbot.Outcome, teams []*foosbot.Team) string {
	for i := range teams {
		ctx.AddTeam(teams[i])
	}
	for i := range outcomes {
		ctx.AddMatchWithOutcome(outcomes[i])
	}
	ids := []string{}
	for i := range outcomes {
		ids = append(ids, outcomes[i].ShortID())
	}
	idsStr := strings.Join(ids, ",")
	return fmt.Sprintf("%d matches %q registered to history.", len(outcomes), idsStr)

}

func getTeamStatsCommand(ctx *foosbot.Context, team *foosbot.Team) string {
	stats := ctx.TeamStats(team)
	if stats.PlayedGames == 0 {
		return fmt.Sprintf("%s hasn't played any match yet.", ctx.Print(team))
	}
	response := fmt.Sprintf("*Team %s*\n", ctx.Print(team))
	response += fmt.Sprintf("Played %d matches (%d wins - %d defeats)\n", stats.PlayedGames, stats.Wins, stats.Defeats)
	response += fmt.Sprintf("```Recent match history for team %s:\n", team.ShortID())
	for i := range stats.Matches {
		idx := len(stats.Outcomes) - 1 - i
		outcome := stats.Outcomes[idx]
		match := stats.Matches[idx]
		outcomeStr := "Won"
		if outcome.IsLooser(team) {
			outcomeStr = "Lost"
		}
		response += fmt.Sprintf("%s: %s %s (%s)\n", match.ShortID(), outcomeStr, ctx.Print(outcome),
			humanize.Time(match.PlayedAt))
		if i >= 10 {
			break
		}
	}
	response += "```"
	return response
}

func getPlayerStatsCommand(ctx *foosbot.Context, player *foosbot.Player) string {
	stats := ctx.PlayerStats(player)
	if stats.PlayedGames == 0 {
		return fmt.Sprintf("%s hasn't played any match yet.", ctx.Print(player))
	}
	response := fmt.Sprintf("*Player %s*\n", ctx.Print(player))
	response += fmt.Sprintf("Played %d matches (%d wins - %d defeats)\n", stats.PlayedGames, stats.Wins, stats.Defeats)
	response += fmt.Sprintf("```Recent match history for %s:\n", player.Name)
	for i := range stats.Matches {
		idx := len(stats.Outcomes) - 1 - i
		outcome := stats.Outcomes[idx]
		match := stats.Matches[idx]
		wt, _ := ctx.Query.TeamByID(outcome.WinnerID)
		lt, _ := ctx.Query.TeamByID(outcome.LooserID)
		team := ctx.Query.TeamWithPlayer([]*foosbot.Team{wt, lt}, player)
		outcomeStr := "Won"
		if outcome.IsLooser(team) {
			outcomeStr = "Lost"
		}
		response += fmt.Sprintf("%s: %s %s (%s)\n", match.ShortID(), outcomeStr, ctx.Print(outcome),
			humanize.Time(match.PlayedAt))
		if i >= 10 {
			break
		}
	}
	response += "```"
	return response
}

func process(ctx *foosbot.Context, msg *slack.MessageEvent) (response string) {
	in := []byte(msg.Text)
	r := bytes.NewReader(in)
	p := parsing.NewParser(r)

	token, err := p.ParseCommand()
	if err == parsing.ErrNotFoosbotCommand {
		return
	} else if err != nil {
		response = fmt.Sprintf("%s", err)
		return
	}
	fmt.Println(time.Now().String(), msg.Text)
	switch token.Type {
	case parsing.TokenCommandMatch:
		outcomes, teams, err := p.ParseMatch()
		if err != nil {
			response = err.Error()
			return
		}
		response = addMatchCommand(ctx, outcomes, teams)
	case parsing.TokenCommandStats:
		iface, err := p.ParseStats()
		if err != nil {
			response = err.Error()
			return
		}
		switch obj := iface.(type) {
		case *foosbot.Team:
			response = getTeamStatsCommand(ctx, obj)
		case *foosbot.Player:
			response = getPlayerStatsCommand(ctx, obj)
		}

	}
	return
}

func loadToken() (string, error) {
	tBytes, err := ioutil.ReadFile(".access_token")
	tBytes = bytes.Trim(tBytes, " \n\r\t")
	return string(tBytes), err
}

func backup(c *foosbot.Context) {
	for {
		time.Sleep(1 * time.Hour)
		c.Store()
	}
}

func exit() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	return ch
}

func main() {
	flag.Parse()

	ctx := foosbot.NewContext()
	err := ctx.Load()
	if err != nil {
		log.Printf("cannot load repository")
		return
	}
	go backup(ctx)

	token, err := loadToken()
	if err != nil {
		log.Printf("cannot open slack access token: %s", err)
		return
	}

	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	incomingExit := exit()
Loop:
	for {
		select {
		case <-incomingExit:
			break Loop
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
			case *slack.ConnectedEvent:
			case *slack.PresenceChangeEvent:
			case *slack.LatencyReport:
			case *slack.RTMError:
			case *slack.MessageEvent:
				response := process(ctx, ev)
				rtm.SendMessage(rtm.NewOutgoingMessage(response, ev.Channel))
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop
			default:
				// Ignore other events.
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}

	err = ctx.Store()
	if err != nil {
		log.Printf("cannot store repository")
		return
	}
}
