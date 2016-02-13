package foosbot

import (
	"encoding/json"
	"github.com/cheggaaa/pb"
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
	Matches = []*Match{}
	MatchHistory = []string{}
	MatchesMap = map[string]*Match{}

	Teams = []*Team{}
	TeamsMap = map[string]*Team{}

	Players = []*Player{}
	PlayersMap = map[string]*Player{}
	PlayersNameMap = map[string]*Player{}
}

func Store() {
	s := State{Matches: Matches, MatchHistory: MatchHistory}

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
	bar = pb.StartNew(len(s.MatchHistory))
	for _, matchID := range s.MatchHistory {
		match, ok := matchIndex[matchID]
		if !ok {
			log.Panicf("corrupted history %q", matchID)
		}
		bar.Increment()
		AddMatch(match)
	}
	bar.Finish()
}
