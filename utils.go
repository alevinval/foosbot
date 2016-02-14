package foosbot

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
)

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

func readGzFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	return ioutil.ReadAll(gzr)
}

func writeGzFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	gzw, err := gzip.NewWriterLevel(f, gzip.BestSpeed)
	if err != nil {
		return err
	}
	defer gzw.Close()

	_, err = gzw.Write(data)
	return err
}
