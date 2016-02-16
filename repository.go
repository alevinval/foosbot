package foosbot

import (
	"compress/gzip"
	"encoding/json"
	"os"
)

const (
	dbVersion = "0.0.1"
)

type repository struct {
	Version  string     `json:"version"`
	Matches  []*Match   `json:"history"`
	Outcomes []*Outcome `json:"matches"`
	Teams    []*Team    `json:"teams"`
}

func storeRepository(c *Context, compress bool) error {
	repo := &repository{Version: dbVersion, Matches: c.Matches, Outcomes: c.Outcomes, Teams: c.Teams}
	f, err := os.OpenFile(c.RepositoryName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	if compress {
		gzw, err := gzip.NewWriterLevel(f, gzip.BestSpeed)
		if err != nil {
			return err
		}
		defer gzw.Close()
		return json.NewEncoder(gzw).Encode(repo)
	} else {
		return json.NewEncoder(f).Encode(repo)
	}
}

func loadRepository(path string, decompress bool) (repo *repository, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	repo = new(repository)
	if decompress {
		gzr, zerr := gzip.NewReader(f)
		if zerr != nil {
			err = zerr
			return
		}
		defer gzr.Close()
		err = json.NewDecoder(gzr).Decode(repo)
	} else {
		err = json.NewDecoder(f).Decode(repo)
	}
	return
}
