package foosbot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type State struct {
	Matches      []*Match
	MatchHistory []string
}

var (
	Matches      = []*Match{}
	MatchHistory = []string{}
	MatchesMap   = map[string]*Match{}

	Teams    = []*Team{}
	TeamsMap = map[string]*Team{}

	Players        = []*Player{}
	PlayersMap     = map[string]*Player{}
	PlayersNameMap = map[string]*Player{}
)

func Reset() {
	Players = []*Player{}
	Teams = []*Team{}
	Matches = []*Match{}
}

func Store() {
	s := State{Matches: Matches, MatchHistory: MatchHistory}
	data, err := json.Marshal(s)
	if err != nil {
		log.Printf("cannot store state: %s", err)
		return
	}
	err = ioutil.WriteFile("foosbot.db", data, os.ModePerm)
	if err != nil {
		log.Fatalf("cannot store state: %s", err)
		return
	}
}

func Load() {
	s := new(State)
	data, err := ioutil.ReadFile("foosbot.db")
	if err != nil {
		log.Printf("cannot load state: %s", err)
		return
	}
	err = json.Unmarshal(data, s)
	if err != nil {
		log.Printf("cannot load state: %s", err)
		return
	}
	log.Printf("loading state from db")
	matchIndex := map[string]*Match{}
	for _, match := range s.Matches {
		matchIndex[match.ID] = match
	}
	for _, matchID := range s.MatchHistory {
		match, ok := matchIndex[matchID]
		if !ok {
			log.Panicf("corrupted history %q", matchID)
		}
		log.Printf("loading match %q", match.ShortID())
		AddMatch(match)
	}
}
