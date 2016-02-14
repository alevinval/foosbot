package main

import (
	"bytes"
	"fmt"
	"github.com/alevinval/foosbot"
	"github.com/alevinval/foosbot/parsing"
	"github.com/alevinval/slack"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

func handleMatch(parser *parsing.Parser, client *slack.Client, message slack.Message) string {
	matches, err := parser.ParseMatch()
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	for i := range matches {
		foosbot.AddMatchWithHistory(matches[i])
	}
	matchIds := []string{}
	for i := range matches {
		matchIds = append(matchIds, matches[i].ShortID())
	}
	ids := strings.Join(matchIds, ",")
	return fmt.Sprintf("%d matches %q registered to history.", len(matches), ids)

}

func handleStats(parser *parsing.Parser, client *slack.Client, message slack.Message) string {
	team, err := parser.ParseStats()
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	stats := foosbot.TeamStats(team)
	if stats.PlayedGames == 0 {
		response := fmt.Sprintf("%s hasn't played any match yet.", team)
		return response
	}
	response := fmt.Sprintf("%s has played %d matches, with a stunning record of %d wins and "+
		"%d defeats.\nRecent match history:\n",
		team, stats.PlayedGames, stats.Wins, stats.Defeats)
	for i := range stats.Matches {
		if i > 3 {
			break
		}
		match := stats.Matches[len(stats.Matches)-1-i]
		history := stats.History[len(stats.History)-1-i]
		if match.WinnerID == team.ID {
			response += fmt.Sprintf("- Won against %s (%s)\n", match.Looser(), humanize.Time(history.PlayedAt))
		} else {
			response += fmt.Sprintf("- Lost against %s (%s)\n", match.Winner(), humanize.Time(history.PlayedAt))
		}
	}
	return response
}

func loadAccessToken() (string, error) {
	tBytes, err := ioutil.ReadFile(".access_token")
	return string(tBytes), err
}

func run(client *slack.Client) {
	log.Printf("connected to slack")
	gNames := []string{}
	for _, group := range client.Groups() {
		gNames = append(gNames, group.Name)
	}
	log.Printf("foosbot listening on: %s", gNames)

	for message := range client.Receiver() {
		in := []byte(message.Text)
		r := bytes.NewReader(in)
		p := parsing.NewParser(r)

		token, err := p.ParseCommand()
		if err == parsing.ErrNotFoosbotCommand {
			continue
		} else if err != nil {
			response := fmt.Sprintf("%s", err)
			client.Respond(message.Channel, response)
			continue
		}
		switch token.Type {
		case parsing.TokenCommandMatch:
			response := handleMatch(p, client, message)
			client.Respond(message.Channel, response)
		case parsing.TokenCommandStats:
			response := handleStats(p, client, message)
			client.Respond(message.Channel, response)
		}
	}
}

func periodicBackup() {
	for {
		time.Sleep(1 * time.Hour)
		foosbot.Context.Store()
	}
}

func main() {
	token, err := loadAccessToken()
	if err != nil {
		log.Printf("cannot open slack access token: %s", err)
		return
	}

	client, err := slack.OpenRTM(token)
	if err != nil {
		log.Printf("cannot connect to slack: %s", err)
		return
	}
	foosbot.Context.Load()
	go run(client)
	go periodicBackup()
	<-sysExit()
	foosbot.Context.Store()
}

func sysExit() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	return ch
}
