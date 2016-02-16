package foosbot

import (
	"flag"
	"github.com/cheggaaa/pb"
	"log"
	"os"
)

const (
	DefaultRepositoryName = "foosbot.db"
)

var (
	compress   = flag.Bool("compress", false, "compresses the database")
	decompress = flag.Bool("decompress", false, "decompresses the database")
)

type Context struct {
	RepositoryName string
	Query          Queries
	Matches        []*Match
	Outcomes       []*Outcome
	Teams          []*Team
	Players        []*Player
	outcomesMap    map[string]*Outcome
	teamsMap       map[string]*Team
	playersMap     map[string]*Player
	playersNameMap map[string]*Player
}

func (ctx *Context) Reset() {
	ctx.Matches = []*Match{}
	ctx.Outcomes = []*Outcome{}
	ctx.Teams = []*Team{}
	ctx.Players = []*Player{}
	ctx.outcomesMap = map[string]*Outcome{}
	ctx.teamsMap = map[string]*Team{}
	ctx.playersMap = map[string]*Player{}
	ctx.playersNameMap = map[string]*Player{}
}

func NewContext() *Context {
	return newContext()
}

func newContext() *Context {
	ctx := new(Context)
	ctx.Query = Queries{ctx}
	ctx.RepositoryName = DefaultRepositoryName
	ctx.Reset()
	return ctx
}

func (ctx *Context) AddMatchWithOutcome(outcome *Outcome) {
	match := NewMatch(outcome)
	ctx.AddMatch(match, outcome)
}

func (ctx *Context) AddMatch(match *Match, outcome *Outcome) {
	ctx.Matches = append(ctx.Matches, match)
	m, ok := ctx.outcomesMap[outcome.ID]
	if ok {
		outcome = m
		m.Occurrences++
		return
	}
	outcome.Occurrences++
	ctx.Outcomes = append(ctx.Outcomes, outcome)
	ctx.outcomesMap[outcome.ID] = outcome
	return
}

func (ctx *Context) AddTeam(team *Team) {
	_, ok := ctx.teamsMap[team.ID]
	if ok {
		return
	}
	ctx.Teams = append(ctx.Teams, team)
	ctx.teamsMap[team.ID] = team

	for _, player := range team.Players {
		ctx.AddPlayer(player)
	}
	return
}

func (ctx *Context) AddPlayer(player *Player) {
	_, ok := ctx.playersMap[player.ID]
	if ok {
		return
	}
	ctx.Players = append(ctx.Players, player)
	ctx.playersMap[player.ID] = player
	ctx.playersNameMap[player.Name] = player
	return
}

func (ctx *Context) Store() error {
	log.Println("storing repository")
	err := storeRepository(ctx, *compress)
	if err != nil {
		log.Printf("error storing repository: %s", err)
		return err
	}
	return nil
}

func (ctx *Context) Load() error {
	log.Println("loading repository")
	repo, err := loadRepository(ctx.RepositoryName, *decompress)
	if os.IsNotExist(err) {
		log.Printf("repository not found")
		return nil
	} else if err != nil {
		log.Printf("error loading repository: %s", err)
		return err
	}

	log.Println("loading teams")
	bar := pb.StartNew(len(repo.Teams))
	for _, team := range repo.Teams {
		ctx.AddTeam(team)
		bar.Increment()
	}
	bar.Finish()

	log.Println("loading match history")
	loadedOutcomes := map[string]*Outcome{}
	bar = pb.StartNew(len(repo.Outcomes))
	for _, outcome := range repo.Outcomes {
		loadedOutcomes[outcome.ID] = outcome
		bar.Increment()
	}
	bar.Finish()

	bar = pb.StartNew(len(repo.Matches))
	for _, match := range repo.Matches {
		outcome, ok := loadedOutcomes[match.OutcomeID]
		if !ok {
			log.Panicf("corrupted history %q", match.OutcomeID)
		}
		ctx.AddMatch(match, outcome)
		bar.Increment()
	}
	bar.Finish()
	return nil
}
