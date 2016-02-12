package foosbot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type State struct {
	Players []Player
	Teams   []*Team
	Matches []*Match
}

var (
	digest = sha256.New()
)

func hash(input ...string) string {
	digest.Reset()
	for _, in := range input {
		digest.Write([]byte(in))
	}
	h := digest.Sum(nil)
	return hex.EncodeToString(h)
}

func Reset() {
	players = []Player{}
	teams = []*Team{}
	matches = []*Match{}
}

func Store() {
	s := State{Matches: matches}
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
	for _, match := range s.Matches {
		log.Printf("loading match %q", match.ShortID())
		AddMatch(match)
	}
}
