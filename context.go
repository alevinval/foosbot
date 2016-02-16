package foosbot

import (
	"github.com/cheggaaa/pb"
	"log"
	"os"
)

const (
	DefaultRepositoryName = "foosbot.db"
)

type Context struct {
	RepositoryName string
	Query          Queries

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
	c.Query = Queries{c}
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

func (c *Context) Store() error {
	log.Println("storing repository")
	err := storeRepository(c)
	if err != nil {
		log.Printf("error storing repository: %s", err)
		return err
	}
	return nil
}

func (c *Context) Load() error {
	log.Println("loading repository")
	repo, err := loadRepository(c.RepositoryName)
	if os.IsNotExist(err) {
		log.Printf("repository not found")
		return nil
	} else if err != nil {
		log.Printf("error loading repository: %s", err)
		return err
	}
	log.Println("building match history")
	loadedMatches := map[string]*Match{}
	bar := pb.StartNew(len(repo.Matches))
	for _, match := range repo.Matches {
		loadedMatches[match.ID] = match
		bar.Increment()
	}
	bar.Finish()

	bar = pb.StartNew(len(repo.History))
	for _, historyEntry := range repo.History {
		match, ok := loadedMatches[historyEntry.MatchID]
		if !ok {
			log.Panicf("corrupted history %q", historyEntry.MatchID)
		}
		c.AddMatch(match, historyEntry)
		bar.Increment()
	}
	bar.Finish()
	return nil
}
