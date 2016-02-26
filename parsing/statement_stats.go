package parsing

import "github.com/alevinval/foosbot"

type (
	StatStatement struct {
		Names []string
	}
)

func (s *StatStatement) Execute(ctx *foosbot.Context) (*foosbot.Stats, error) {
	if len(s.Names) == 1 {
		player := foosbot.NewPlayer(s.Names[0])
		stats := ctx.PlayerStats(player)
		return &stats.Stats, nil

	} else {
		players := []*foosbot.Player{}
		for _, name := range s.Names {
			player := foosbot.NewPlayer(name)
			players = append(players, player)
		}
		team, err := foosbot.NewTeam(players...)
		if err != nil {
			return nil, newCommandError(err.Error())
		}
		stats := ctx.TeamStats(team)
		return &stats.Stats, nil
	}
}
