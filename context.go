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

func newContext() *Context {
	c := new(Context)
	c.RepositoryName = DefaultRepositoryName
	c.Reset()
	return c
}

func NewContext() *Context {
	c := newContext()
	return c
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
