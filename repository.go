package foosbot

import (
	"compress/gzip"
	"encoding/json"
	"os"
)

type repository struct {
	History []*HistoryEntry `json:"history"`
	Matches []*Match        `json:"matches"`
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
