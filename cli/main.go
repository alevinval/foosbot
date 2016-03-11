package main

import (
	"bytes"
	"fmt"
	"github.com/alevinval/foosbot"
	"github.com/alevinval/foosbot/parsing"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

func addMatchCommand(ctx *foosbot.Context, statement *parsing.MatchStatement) string {
	err := ctx.ExecuteMatch(*statement)
	if err != nil {
		panic(err)
	}
	total := statement.TeamOneScore + statement.TeamTwoScore
	return fmt.Sprintf("%d matches registered to history.", total)
}

func getLeaderboard(ctx *foosbot.Context) string {
	stats := ctx.PlayersStatsFromMatches(10, 10)
	response := ctx.ReportLeaderBoard(stats)
	return response
}

func statsCommand(ctx *foosbot.Context, statement *parsing.StatStatement) string {
	stats, err := ctx.ExecuteStats(*statement)
	if err != nil {
		return err.Error()
	}
	return ctx.ReportStats(stats)
}

func process(ctx *foosbot.Context, p *parsing.Parser) string {
	token, err := p.ParseCommand()
	if err == parsing.ErrNotFoosbotCommand {
		return "no fb command"
	} else if err != nil {
		return fmt.Sprintf("%s", err)
	}
	switch token.Type {
	case parsing.TokenCommandMatch:
		matchStatement, err := p.ParseMatch()
		if err != nil {
			return err.Error()
		}
		return addMatchCommand(ctx, matchStatement)
	case parsing.TokenCommandLeaderboard:
		return getLeaderboard(ctx)
	case parsing.TokenCommandStats:
		statStatement, err := p.ParseStats()
		if err != nil {
			return err.Error()
		}
		return statsCommand(ctx, statStatement)
	default:
		return "Not yet supported"
	}
	return "What"
}

func main() {
	foosbot.Verbose = false
	ctx := foosbot.NewContext()
	err := ctx.Load()
	if err != nil {
		panic(err)
	}
	q := foosbot.Query(ctx)

	var limit int

	app := cli.NewApp()
	app.Name = "foosbot"
	app.Usage = "The foosball match tracker"
	app.Commands = []cli.Command{
		{
			Name:    "match",
			Aliases: []string{"m"},
			Usage:   "Match history listing",
			Action: func(c *cli.Context) {
				matches := q.Matches().Limit(limit).Get()
				for _, match := range matches {
					fmt.Println(foosbot.Print(match))
				}
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "limit",
					Value:       10,
					Usage:       "Limits to the specified amount of results",
					Destination: &limit,
				},
			},
		},
		{
			Name: "run",
			Action: func(c *cli.Context) {
				args := c.Args()
				s := args.First()
				s += " "
				s += strings.Join(args.Tail(), " ")
				r := bytes.NewReader([]byte(s))
				parser := parsing.NewParser(r)
				out := process(ctx, parser)
				fmt.Println(out)
				ctx.Store()
			},
		},
	}
	app.Run(os.Args)
}
