package main

import (
	"bytes"
	"fmt"
	"github.com/alevinval/bobbie/auth"
	"github.com/alevinval/bobbie/slack"
	"github.com/alevinval/foosbot"
	"github.com/alevinval/foosbot/parsing"
	"log"
	"os"
	"os/signal"
	"strings"
)

func handleMatch(parser *parsing.Parser, client *slack.Client, message slack.Message) string {
	matches, err := parser.ParseMatch()
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	for _, match := range matches {
		foosbot.AddMatch(match)
	}

	matchIds := []string{}
	for _, match := range matches {
		matchIds = append(matchIds, match.ShortID())
	}
	ids := strings.Join(matchIds, ",")
	return fmt.Sprintf("%d matches %q registered to history.", len(matches), ids)

}

func handleStats(parser *parsing.Parser, client *slack.Client, message slack.Message) string {
	playerNames, err := parser.ParseStats()
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	players := []foosbot.Player{}
	for _, name := range playerNames {
		player, ok := foosbot.PlayerByName(name)
		if !ok {
			return fmt.Sprintf("Player %q not found", name)
		}
		players = append(players, player)
	}
	team, ok := foosbot.TeamByPlayers(players...)
	if !ok {
		return fmt.Sprintf("Those players are not in a team")
	}
	wins, defeats := 0, 0
	teamMatches := foosbot.MatchesWithTeam(team)
	for _, match := range teamMatches {
		if match.Winner == team.ID {
			wins += match.N
		} else {
			defeats += match.N
		}
	}
	if len(teamMatches) == 0 {
		response := fmt.Sprintf("Team \"%s\" - (%s, %s) hasn't played any match yet.",
			team.ShortID(), team.Players[0].Name, team.Players[1].Name)
		return response
	}
	lastMatch := teamMatches[len(teamMatches)-1]
	opponent := lastMatch.Teams[0]
	if opponent.ID == team.ID {
		opponent = lastMatch.Teams[1]
	}
	outcome := "lost"
	if team.ID == lastMatch.Winner {
		outcome = "won"
	}
	response := fmt.Sprintf("Team \"%s\" - (%s, %s) has played %d matches, with a stunning record of %d wins and "+
		"%d defeats.\nTheir last match was against \"%s\" - (%s, %s) and they %s",
		team.ShortID(), team.Players[0].Name, team.Players[1].Name, wins+defeats, wins, defeats,
		lastMatch.ShortID(), opponent.Players[0].Name, opponent.Players[1].Name, outcome)
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

		command, err := p.ParseCommand()
		if err == parsing.ErrNotFoosbotCommand {
			continue
		} else if err != nil {
			response := fmt.Sprintf("%s", err)
			client.Say(message.Channel, response)
			continue
		}
		switch command {
		case parsing.TokenCommandMatch:
			response := handleMatch(p, client, message)
			client.Say(message.Channel, response)
		case parsing.TokenCommandStats:
			response := handleStats(p, client, message)
			client.Say(message.Channel, response)
		}
	}
}

func main() {
	foosbot.Load()
	go run()
	<-sysExit()
	foosbot.Store()
}

func sysExit() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	return ch
}
