package main

import (
	"fmt"
	"github.com/alevinval/foosbot"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	foosbot.Verbose = false
	ctx := foosbot.NewContext()
	ctx.Load()
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
					fmt.Println(ctx.Print(match))
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
	}
	app.Run(os.Args)
}
