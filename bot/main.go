package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/alevinval/foosbot"
	"github.com/alevinval/foosbot/parsing"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"
)

func addMatchCommand(ctx *foosbot.Context, statement *parsing.MatchStatement) string {
	statement.Execute(ctx)
	total := statement.TeamOneScore + statement.TeamTwoScore
	return fmt.Sprintf("%d matches registered to history.", total)
}

func getLeaderboard(ctx *foosbot.Context) string {
	stats := ctx.PlayersStatsFromMatches(10, 10)
	response := ctx.ReportLeaderBoard(stats)
	return response
}
func getTeamStatsCommand(ctx *foosbot.Context, team *foosbot.Team) string {
	stats := ctx.TeamStats(team)
	response := ctx.ReportStats(&stats.Stats, team)
	return response
}

func getPlayerStatsCommand(ctx *foosbot.Context, player *foosbot.Player) string {
	stats := ctx.PlayerStats(player)
	response := ctx.ReportStats(&stats.Stats, player)
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
		matchStatement, err := p.ParseMatch()
		if err != nil {
			response = err.Error()
			return
		}
		response = addMatchCommand(ctx, matchStatement)
	case parsing.TokenCommandLeaderboard:
		response = getLeaderboard(ctx)
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
