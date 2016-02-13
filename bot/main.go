package main

import (
	"bytes"
	"fmt"
	"github.com/alevinval/bobbie/auth"
	"github.com/alevinval/bobbie/slack"
	"github.com/alevinval/foosbot"
	"github.com/alevinval/foosbot/parsing"
	"github.com/dustin/go-humanize"
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

func run() {
	accessToken, err := auth.AccessToken()
	if err != nil {
		log.Fatalf("error loading access token: %s", err)
	}
	client := slack.NewClient(accessToken.AccessToken)
	client.StartRTM()
	log.Printf("foosbot: slack rtm started")
	groupNames := []string{}
	for _, group := range client.Groups {
		groupNames = append(groupNames, group.Name)
	}
	log.Printf("foosbot running on: %s", groupNames)

	for message := range client.Messages {
		in := []byte(message.Text)
		r := bytes.NewReader(in)
		p := parsing.NewParser(r)

		token, err := p.ParseCommand()
		if err == parsing.ErrNotFoosbotCommand {
			continue
		} else if err != nil {
			response := fmt.Sprintf("%s", err)
			client.Say(message.Channel, response)
			continue
		}
		switch token.Type {
		case parsing.TokenCommandMatch:
			response := handleMatch(p, client, message)
			client.Say(message.Channel, response)
		case parsing.TokenCommandStats:
			response := handleStats(p, client, message)
			client.Say(message.Channel, response)
		}
	}
}

func periodicBackup() {
	for {
		time.Sleep(1 * time.Hour)
		foosbot.Store()
	}
}

func main() {
	foosbot.Load()
	go run()
	go periodicBackup()
	<-sysExit()
	foosbot.Store()
}

func sysExit() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	return ch
}
