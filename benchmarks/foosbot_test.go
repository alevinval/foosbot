package benchmarks_test

import (
	"github.com/alevinval/foosbot"
	"io/ioutil"
	"log"
	"math/rand"
	"testing"
)

var (
	letters = [][]rune{
		[]rune("abcd"),
		[]rune("efgh"),
		[]rune("ijkl"),
		[]rune("mnop"),
	}
	_ = setup()
)

func setup() bool {
	log.SetOutput(ioutil.Discard)
	return true
}

func randomPlayers(n int) []*foosbot.Player {
	players := []*foosbot.Player{}
	i := 0
	for i < 4 {
		n := rand.Intn(4)
		name := string(letters[i][n])
		players = append(players, foosbot.NewPlayer(name))
		i++
	}
	return players
}

func randomMatches(n int) []*foosbot.Match {
	matches := []*foosbot.Match{}
	for n > 0 {
		players := randomPlayers(4)
		rin := rand.Perm(4)
		t1 := foosbot.NewTeam(players[rin[0]], players[rin[1]])
		t2 := foosbot.NewTeam(players[rin[2]], players[rin[3]])
		m := foosbot.NewMatch(t1, t2)
		matches = append(matches, m)
		n--
	}
	return matches
}

func addMatches(m []*foosbot.Match) {
	for k := 0; k < len(m); k++ {
		foosbot.AddMatch(m[k])

	}
}
func BenchmarkCreate100KHistory(b *testing.B) {
	m := randomMatches(100000)
	benchmarkBuildHistory(b, m)
}

func BenchmarkStore100KMatches(b *testing.B) {
	m := randomMatches(100000)
	addMatches(m)
	for i := 0; i < b.N; i++ {
		benchmarkStoreState(b)
	}
}

func BenchmarkLoad100KMatches(b *testing.B) {
	m := randomMatches(100000)
	addMatches(m)
	foosbot.Store()
	foosbot.Reset()
	benchmarkLoadState(b)
}

func benchmarkBuildHistory(b *testing.B, m []*foosbot.Match) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		addMatches(m)
	}
}

func benchmarkStoreState(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		foosbot.Store()
	}
}

func benchmarkLoadState(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		foosbot.Load()
	}
}
