package foosbot

const (
	DefaultRepositoryName = "foosbot.db"
)

type Context struct {
	RepositoryName string

	// Actual context
	History        []*HistoryEntry
	Matches        []*Match
	MatchesMap     map[string]*Match
	Teams          []*Team
	TeamsMap       map[string]*Team
	Players        []*Player
	PlayersMap     map[string]*Player
	PlayersNameMap map[string]*Player
}

func (c *Context) Reset() {
	c.History = []*HistoryEntry{}
	c.Matches = []*Match{}
	c.MatchesMap = map[string]*Match{}
	c.Teams = []*Team{}
	c.TeamsMap = map[string]*Team{}
	c.Players = []*Player{}
	c.PlayersMap = map[string]*Player{}
	c.PlayersNameMap = map[string]*Player{}
}

func NewContext() *Context {
	return newContext()
}

func newContext() *Context {
	c := new(Context)
	c.RepositoryName = DefaultRepositoryName
	c.Reset()
	return c
}

func (c *Context) AddMatchWithHistory(match *Match) {
	entry := NewHistoryEntry(match)
	c.AddMatch(match, entry)
}

func (c *Context) AddMatch(match *Match, entry *HistoryEntry) {
	c.History = append(c.History, entry)
	m, ok := c.MatchesMap[match.ID]
	if ok {
		match = m
		m.N++
		return
	}
	match.N++
	c.Matches = append(c.Matches, match)
	c.MatchesMap[match.ID] = match
	for _, team := range match.Teams {
		c.AddTeam(team)
	}
	return
}

func (c *Context) AddTeam(team *Team) {
	_, ok := c.TeamsMap[team.ID]
	if ok {
		return
	}
	c.Teams = append(c.Teams, team)
	c.TeamsMap[team.ID] = team

	for _, player := range team.Players {
		c.AddPlayer(player)
	}
	return
}

func (c *Context) AddPlayer(player *Player) {
	_, ok := c.PlayersMap[player.ID]
	if ok {
		return
	}
	c.Players = append(c.Players, player)
	c.PlayersMap[player.ID] = player
	c.PlayersNameMap[player.Name] = player
	return
}
