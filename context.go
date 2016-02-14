package foosbot

import (
	"encoding/json"
	"github.com/cheggaaa/pb"
	"log"
	"os"
)

const (
	DefaultDatabaseName = "foosbot.db"
)

type repository struct {
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

func (c *context) Store() error {
	log.Printf("storing repository")
	repo := repository{History: Context.History, Matches: Context.Matches}
	data, err := json.Marshal(repo)
	if err != nil {
		log.Printf("error serializing repository: %s", err)
		return err
	}
	writeGzFile(c.DatabaseName, data)
	return nil
}

func (c *context) Load() error {
	log.Println("loading repository")
	data, err := readGzFile(c.DatabaseName)
	if os.IsNotExist(err) {
		log.Printf("repository not found")
		return nil
	} else if err != nil {
		log.Printf("error reading repository: %s", err)
		return err
	}

	repository := new(repository)
	err = json.Unmarshal(data, repository)
	if err != nil {
		log.Printf("error loading repository: %s", err)
		return err
	}

	log.Println("building match history")
	loadedMatches := map[string]*Match{}
	bar := pb.StartNew(len(repository.Matches))
	for _, match := range repository.Matches {
		loadedMatches[match.ID] = match
		bar.Increment()
	}
	bar.Finish()
	bar = pb.StartNew(len(repository.History))
	for _, entry := range repository.History {
		match, ok := loadedMatches[entry.MatchID]
		if !ok {
			log.Panicf("corrupted history %q", entry.MatchID)
		}
		AddMatch(match, entry)
		bar.Increment()
	}
	bar.Finish()

	return nil
}
