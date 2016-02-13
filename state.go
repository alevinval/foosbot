package foosbot

import (
	"encoding/json"
	"github.com/cheggaaa/pb"
	"io/ioutil"
	"log"
	"os"
)

type State struct {
	Matches []*Match
	History []*HistoryEntry
}

var (
	Matches    = []*Match{}
	History    = []*HistoryEntry{}
	MatchesMap = map[string]*Match{}

	Teams    = []*Team{}
	TeamsMap = map[string]*Team{}

	Players        = []*Player{}
	PlayersMap     = map[string]*Player{}
	PlayersNameMap = map[string]*Player{}
)

func Reset() {
	Matches = []*Match{}
	History = []*HistoryEntry{}
	MatchesMap = map[string]*Match{}

	Teams = []*Team{}
	TeamsMap = map[string]*Team{}

	Players = []*Player{}
	PlayersMap = map[string]*Player{}
	PlayersNameMap = map[string]*Player{}
}

func Store() {
	s := State{Matches: Matches, History: History}

	log.Println("serialising database")
	data, err := json.Marshal(s)
	if err != nil {
		log.Printf("cannot store state: %s", err)
		return
	}

	log.Println("dumping database")
	err = ioutil.WriteFile("foosbot.db", data, os.ModePerm)
	if err != nil {
		log.Fatalf("cannot store state: %s", err)
		return
	}
}

func Load() {
	s := new(State)
	log.Printf("reading database")
	data, err := ioutil.ReadFile("foosbot.db")
	if err != nil {
		log.Printf("cannot load state: %s", err)
		return
	}

	log.Printf("loading database")
	err = json.Unmarshal(data, s)
	if err != nil {
		log.Printf("cannot load state: %s", err)
		return
	}
	matchIndex := map[string]*Match{}

	log.Println("loading matches")
	bar := pb.StartNew(len(s.Matches))
	for _, match := range s.Matches {
		matchIndex[match.ID] = match
		bar.Increment()
	}
	bar.Finish()

	log.Println("rebuilding match history")
	bar = pb.StartNew(len(s.History))
	for _, entry := range s.History {
		match, ok := matchIndex[entry.MatchID]
		if !ok {
			log.Panicf("corrupted history %q", entry.MatchID)
		}
		bar.Increment()
		AddMatch(match, entry)
	}
	bar.Finish()
}
