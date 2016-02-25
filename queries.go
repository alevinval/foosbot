package foosbot

type (
	QueryBuilder struct {
		ctx     *Context
		indexes *Indexing
	}
)

func Query(ctx *Context) *QueryBuilder {
	return &QueryBuilder{
		ctx:     ctx,
		indexes: ctx.indexes,
	}
}

func (qb *QueryBuilder) OutcomeByID(id string) (*Outcome, bool) {
	outcome, ok := qb.indexes.outcomesMap[id]
	return outcome, ok
}

func (qb *QueryBuilder) OutcomeByTeams(a, b *Team) (outcome *Outcome, ok bool) {
	id := BuildOutcomeID(a, b)
	return qb.OutcomeByID(id)
}

func (qb *QueryBuilder) TeamByID(id string) (team *Team, ok bool) {
	team, ok = qb.indexes.teamsMap[id]
	return
}

func (qb *QueryBuilder) TeamByPlayers(players ...*Player) (team *Team, ok bool) {
	id := BuildTeamID(players...)
	team, ok = qb.TeamByID(id)
	return
}

func (qb *QueryBuilder) PlayerByID(playerID string) (player *Player, ok bool) {
	player, ok = qb.indexes.playersMap[playerID]
	return
}

func (qb *QueryBuilder) PlayerByName(name string) (player *Player, ok bool) {
	player, ok = qb.indexes.playersNameMap[name]
	return
}

func (qb *QueryBuilder) TeamsWithPlayer(p *Player) (teams []*Team) {
	return qb.indexes.playersTeamMap[p.ID]
}

func (qb *QueryBuilder) TeamWithPlayer(teams []*Team, player *Player) *Team {
	for _, team := range teams {
		_, ok := qb.indexes.teamsPlayerMap[team.ID][player.ID]
		if ok {
			return team
		}
	}
	return nil
}

func (qb *QueryBuilder) Matches() *MatchesQueryBuilder {
	return QueryMatches(qb.ctx)
}
