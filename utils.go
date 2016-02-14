package foosbot

import (
	"crypto/sha256"
	"encoding/hex"
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
