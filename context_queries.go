package foosbot

func (c *Context) MatchByID(id string) (*Match, bool) {
	m, ok := c.MatchesMap[id]
	return m, ok
}

func (c *Context) MatchByTeams(a, b *Team) (match *Match, ok bool) {
	teamID := buildMatchId(a, b)
	match, ok = c.MatchesMap[teamID]
	return
}

func (c *Context) MatchesWithTeam(t *Team) (matches []*Match, history []*HistoryEntry) {
	outcomes := []string{}
	for _, match := range c.Matches {
		for _, team := range match.Teams {
			if t.ID == team.ID {
				outcomes = append(outcomes, match.ID)
				break
			}
		}
	}
	for _, entry := range c.History {
		if in(outcomes, entry.MatchID) {
			m, _ := c.MatchByID(entry.MatchID)
			matches = append(matches, m)
			history = append(history, entry)
		}
	}
	return matches, history
}

func (c *Context) TeamByID(id string) (team *Team, ok bool) {
	team, ok = c.TeamsMap[id]
	return
}

func (c *Context) TeamByPlayers(players ...*Player) (team *Team, ok bool) {
	teamID := buildTeamId(players...)
	team, ok = c.TeamsMap[teamID]
	return
}

func (c *Context) PlayerByID(playerID string) (player *Player, ok bool) {
	player, ok = c.PlayersMap[playerID]
	return
}

func (c *Context) PlayerByName(name string) (player *Player, ok bool) {
	player, ok = c.PlayersNameMap[name]
	return
}

func in(arr []string, m string) bool {
	for i := range arr {
		if arr[i] == m {
			return true
		}
	}
	return false
}
