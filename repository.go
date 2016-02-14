package foosbot

import (
	"encoding/json"
	"log"
	"os"
)

type repository struct {
	History []*HistoryEntry `json:"history"`
	Matches []*Match        `json:"matches"`
}

func storeRepository(c *Context) error {
	r := repository{History: c.History, Matches: c.Matches}
	data, err := json.Marshal(r)
	if err != nil {
		log.Printf("error serializing repository: %s", err)
		return err
	}
	return writeGzFile(c.RepositoryName, data)
}

func loadRepository(path string) (*repository, error) {
	data, err := readGzFile(path)
	if os.IsNotExist(err) {
		log.Printf("repository not found")
		return nil, err
	} else if err != nil {
		log.Printf("error reading repository: %s", err)
		return nil, err
	}
	r := new(repository)
	err = json.Unmarshal(data, r)
	return r, err
}
