package foosbot

import (
	"encoding/json"
	"github.com/cheggaaa/pb"
	"io/ioutil"
	"log"
	"os"
)

const (
	DefaultDatabaseName = "foosbot.db"
)

type Repository struct {
	History []*HistoryEntry `json:"history"`
	Matches []*Match        `json:"matches"`
}

type context struct {
	DatabaseName string

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

func (c *context) Reset() {
	c.History = []*HistoryEntry{}
	c.Matches = []*Match{}
	c.MatchesMap = map[string]*Match{}
	c.Teams = []*Team{}
	c.TeamsMap = map[string]*Team{}
	c.Players = []*Player{}
	c.PlayersMap = map[string]*Player{}
	c.PlayersNameMap = map[string]*Player{}
}

var (
	Context = newContext()
)

func newContext() *context {
	c := new(context)
	c.DatabaseName = DefaultDatabaseName
	c.Reset()
	return c
}

func (c *context) Store() {
	s := Repository{History: Context.History, Matches: Context.Matches}

	log.Println("serialising database")
	data, err := json.Marshal(s)
	if err != nil {
		log.Printf("cannot store state: %s", err)
		return
	}

	log.Println("dumping database")
	err = ioutil.WriteFile(c.DatabaseName, data, os.ModePerm)
	if err != nil {
		log.Fatalf("cannot store state: %s", err)
		return
	}
}

func (c *context) Load() {
	log.Printf("reading database")
	data, err := ioutil.ReadFile(c.DatabaseName)
	if err != nil {
		log.Printf("cannot load state: %s", err)
		return
	}

	log.Printf("loading database")
	repository := new(Repository)
	err = json.Unmarshal(data, repository)
	if err != nil {
		log.Printf("cannot load state: %s", err)
		return
	}
	matchIndex := map[string]*Match{}

	log.Println("loading matches")
	bar := pb.StartNew(len(repository.Matches))
	for _, match := range repository.Matches {
		matchIndex[match.ID] = match
		bar.Increment()
	}
	bar.Finish()

	log.Println("rebuilding match history")
	bar = pb.StartNew(len(repository.History))
	for _, entry := range repository.History {
		match, ok := matchIndex[entry.MatchID]
		if !ok {
			log.Panicf("corrupted history %q", entry.MatchID)
		}
		AddMatch(match, entry)
		bar.Increment()
	}
	bar.Finish()
}
