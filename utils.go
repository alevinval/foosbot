package foosbot

import (
	"crypto/sha256"
	"encoding/hex"
)

func hash(input ...string) string {
	digest := sha256.New()
	for _, in := range input {
		digest.Write([]byte(in))
	}
	h := digest.Sum(nil)
	return hex.EncodeToString(h)
}

func repeated(a, b []string) bool {
	hA, hB := map[string]int{}, map[string]int{}
	for i := range a {
		hA[a[i]] += 1
	}
	for i := range b {
		hA[b[i]] += 1
	}
	for k := range hA {
		if hA[k] > 1 {
			return true
		}
		_, ok := hB[k]
		if ok {
			return true
		}
	}
	return false
}
