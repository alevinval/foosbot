package benchmarks_test

import (
	"github.com/alevinval/foosbot"
	"io/ioutil"
	"log"
	"math/rand"
	"testing"
)

var (
	letters = []string{
		"abcd",
		"efgh",
		"ijkl",
		"mnop",
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

func randomOutcomes(n int) []*foosbot.Outcome {
	outcomes := []*foosbot.Outcome{}
	for n > 0 {
		players := randomPlayers(4)
		rin := rand.Perm(4)
		t1, _ := foosbot.NewTeam(players[rin[0]], players[rin[1]])
		t2, _ := foosbot.NewTeam(players[rin[2]], players[rin[3]])
		outcome, _ := foosbot.NewOutcome(t1, t2)
		outcomes = append(outcomes, outcome)
		n--
	}
	return outcomes
}

func addMatches(c *foosbot.Context, outcomes []*foosbot.Outcome) {
	for k := 0; k < len(outcomes); k++ {
		c.AddMatchWithOutcome(outcomes[k])
	}
}
func BenchmarkCreateBigHistory(b *testing.B) {
	c := foosbot.NewContext()
	m := randomOutcomes(100000)
	benchmarkBuildHistory(b, c, m)
}

func BenchmarkStoreBigHistory(b *testing.B) {
	c := foosbot.NewContext()
	m := randomOutcomes(100000)
	addMatches(c, m)
	for i := 0; i < b.N; i++ {
		benchmarkStoreState(b, c)
	}
}

func BenchmarkLoadBigHistory(b *testing.B) {
	c := foosbot.NewContext()
	m := randomOutcomes(100000)
	addMatches(c, m)
	c.Store()
	benchmarkLoadState(b, c)
}

func benchmarkBuildHistory(b *testing.B, c *foosbot.Context, outcomes []*foosbot.Outcome) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		addMatches(c, outcomes)
	}
}

func benchmarkStoreState(b *testing.B, c *foosbot.Context) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.Store()
	}
}

func benchmarkLoadState(b *testing.B, c *foosbot.Context) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.Reset()
		c.Load()
	}
}
