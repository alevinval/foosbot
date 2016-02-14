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

func AddMatches(c *foosbot.Context, matches []*foosbot.Match) string {
	for i := range matches {
		c.AddMatchWithHistory(matches[i])
	}
	ids := []string{}
	for i := range matches {
		ids = append(ids, matches[i].ShortID())
	}
	idsStr := strings.Join(ids, ",")
	return fmt.Sprintf("%d matches %q registered to history.", len(matches), idsStr)

}

func GetStats(c *foosbot.Context, team *foosbot.Team) string {
	stats := c.TeamStats(team)
	if stats.PlayedGames == 0 {
		return fmt.Sprintf("%s hasn't played any match yet.", team)
	}
	response := fmt.Sprintf("%s has played %d matches, with a stunning record of %d wins and "+
		"%d defeats.\n", team, stats.PlayedGames, stats.Wins, stats.Defeats)
	response += fmt.Sprintf("Recent match history:\n")
	for i := range stats.Matches {
		idx := len(stats.Matches) - 1 - i
		match := stats.Matches[idx]
		history := stats.History[idx]
		if match.WinnerID == team.ID {
			response += fmt.Sprintf("- Won against %s (%s)\n",
				match.Looser(), humanize.Time(history.PlayedAt))
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

func bot(c *foosbot.Context, client *slack.Client) {
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
			matches, err := p.ParseMatch()
			if err != nil {
				client.Respond(message.Channel, err.Error())
				continue
			}
			response := AddMatches(c, matches)
			client.Respond(message.Channel, response)
		case parsing.TokenCommandStats:
			team, err := p.ParseStats()
			if err != nil {
				client.Respond(message.Channel, err.Error())
				continue
			}
			response := GetStats(c, team)
			client.Respond(message.Channel, response)
		}
	}
}

func main() {
	c := foosbot.NewContext()
	err := c.Load()
	if err != nil {
		log.Printf("cannot load repository")
		return
	}
	go backup(c)

	token, err := loadToken()
	if err != nil {
		log.Printf("cannot open slack access token: %s", err)
		return
	}

	client, err := slack.OpenRTM(token)
	if err != nil {
		log.Printf("cannot connect to slack: %s", err)
		return
	}
	go bot(c, client)

	<-exit()
	err = c.Store()
	if err != nil {
		log.Printf("cannot store repository")
		return
	}
}
