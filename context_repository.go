package foosbot

import (
	"compress/gzip"
	"encoding/json"
	"github.com/cheggaaa/pb"
	"log"
	"os"
)

type repository struct {
	History []*HistoryEntry `json:"history"`
	Matches []*Match        `json:"matches"`
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

func storeRepository(c *Context) error {
	repo := &repository{History: c.History, Matches: c.Matches}
	f, err := os.OpenFile(c.RepositoryName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	gzw, err := gzip.NewWriterLevel(f, gzip.BestSpeed)
	if err != nil {
		return err
	}
	defer gzw.Close()

	return json.NewEncoder(gzw).Encode(repo)
}

func loadRepository(path string) (*repository, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	repo := new(repository)
	err = json.NewDecoder(gzr).Decode(repo)
	return repo, err
}
