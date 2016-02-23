package foosbot

import (
	"compress/gzip"
	"encoding/json"
	"os"
)

const (
	DatabaseVersion     = "0.0.1"
	DefaultDatabaseName = "foosbot.db"
)

type Database struct {
	Version  string     `json:"version"`
	Matches  []*Match   `json:"matches"`
	Outcomes []*Outcome `json:"outcomes"`
	Teams    []*Team    `json:"teams"`
}

func StoreDB(ctx *Context, compress bool) error {
	db := &Database{
		Version:  DatabaseVersion,
		Matches:  ctx.Matches,
		Outcomes: ctx.Outcomes,
		Teams:    ctx.Teams,
	}
	f, err := os.OpenFile(ctx.DatabaseName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
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
		return json.NewEncoder(gzw).Encode(db)
	} else {
		return json.NewEncoder(f).Encode(db)
	}
}

func LoadDB(path string, decompress bool) (db *Database, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	db = new(Database)
	if decompress {
		gzr, zerr := gzip.NewReader(f)
		if zerr != nil {
			err = zerr
			return
		}
		defer gzr.Close()
		err = json.NewDecoder(gzr).Decode(db)
	} else {
		err = json.NewDecoder(f).Decode(db)
	}
	return
}
