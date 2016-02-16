package foosbot

type Queries struct {
	ctx *Context
}

func (q *Queries) MatchByID(id string) (*Match, bool) {
	m, ok := q.ctx.MatchesMap[id]
	return m, ok
}

func (q *Queries) MatchByTeams(a, b *Team) (match *Match, ok bool) {
	id := buildMatchId(a, b)
	return q.MatchByID(id)
}

func (q *Queries) MatchesWithTeam(t *Team) (matches []*Match, history []*HistoryEntry) {
	outcomes := []string{}
	for _, match := range q.ctx.Matches {
		for _, team := range match.Teams {
			if t.ID == team.ID {
				outcomes = append(outcomes, match.ID)
				break
			}
		}
	}
	for _, entry := range q.ctx.History {
		if in(outcomes, entry.MatchID) {
			m, _ := q.MatchByID(entry.MatchID)
			matches = append(matches, m)
			history = append(history, entry)
		}
	}
	return matches, history
}

func (q *Queries) TeamByID(id string) (team *Team, ok bool) {
	team, ok = q.ctx.TeamsMap[id]
	return
}

func (q *Queries) TeamByPlayers(players ...*Player) (team *Team, ok bool) {
	id := buildTeamId(players...)
	return q.TeamByID(id)
}

func (q *Queries) PlayerByID(playerID string) (player *Player, ok bool) {
	player, ok = q.ctx.PlayersMap[playerID]
	return
}

func (q *Queries) PlayerByName(name string) (player *Player, ok bool) {
	player, ok = q.ctx.PlayersNameMap[name]
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
