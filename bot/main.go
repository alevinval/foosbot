package main

import (
	"bytes"
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

func addMatchCommand(ctx *foosbot.Context, matches []*foosbot.Match) string {
	for i := range matches {
		ctx.AddMatchWithHistory(matches[i])
	}
	ids := []string{}
	for i := range matches {
		ids = append(ids, matches[i].ShortID())
	}
	idsStr := strings.Join(ids, ",")
	return fmt.Sprintf("%d matches %q registered to history.", len(matches), idsStr)

}

func getStatsCommand(ctx *foosbot.Context, team *foosbot.Team) string {
	stats := ctx.TeamStats(team)
	if stats.PlayedGames == 0 {
		return fmt.Sprintf("%s hasn't played any match yet.", team)
	}
	response := fmt.Sprintf(
		"%s has played %d matches, with a stunning record of %d wins and "+
			"%d defeats.\n", team, stats.PlayedGames, stats.Wins, stats.Defeats)
	response += fmt.Sprintf("Recent match history:\n")
	for i := range stats.Matches {
		idx := len(stats.Matches) - 1 - i
		match := stats.Matches[idx]
		history := stats.History[idx]
		if match.WinnerID == team.ID {
			response += fmt.Sprintf("- Won against %s (%s)\n",
				match.Loosers()[0], humanize.Time(history.PlayedAt))
		} else {
			response += fmt.Sprintf("- Lost against %s (%s)\n",
				match.Winner(), humanize.Time(history.PlayedAt))
		}
		if i >= 4 {
			break
		}
	}
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
	switch token.Type {
	case parsing.TokenCommandMatch:
		matches, err := p.ParseMatch()
		if err != nil {
			response = err.Error()
			return
		}
		response = addMatchCommand(ctx, matches)
	case parsing.TokenCommandStats:
		team, err := p.ParseStats()
		if err != nil {
			response = err.Error()
			return
		}
		response = getStatsCommand(ctx, team)
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
				fmt.Printf("Message: %v\n", ev.Text)
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
