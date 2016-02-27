package foosbot

import (
	"flag"
	"github.com/cheggaaa/pb"
	"log"
	"os"
)

var (
	Verbose    = true
	compress   = flag.Bool("compress", false, "compresses the database")
	decompress = flag.Bool("decompress", false, "decompresses the database")
)

type Indexing struct {
	outcomesMap    map[string]*Outcome
	teamsMap       map[string]*Team
	teamsPlayerMap map[string]map[string]*Player
	playersMap     map[string]*Player
	playersNameMap map[string]*Player
	playersTeamMap map[string][]*Team
}

func (idx *Indexing) Reset() {
	idx.outcomesMap = map[string]*Outcome{}
	idx.teamsMap = map[string]*Team{}
	idx.teamsPlayerMap = map[string]map[string]*Player{}
	idx.playersMap = map[string]*Player{}
	idx.playersNameMap = map[string]*Player{}
	idx.playersTeamMap = map[string][]*Team{}
}

type Context struct {
	DatabaseName string
	Matches      []*Match
	Outcomes     []*Outcome
	Teams        []*Team
	Players      []*Player
	indexes      *Indexing
}

func (ctx *Context) Reset() {
	ctx.Matches = []*Match{}
	ctx.Outcomes = []*Outcome{}
	ctx.Teams = []*Team{}
	ctx.Players = []*Player{}
	ctx.indexes.Reset()
}

func NewContext() *Context {
	return newContext()
}

func newContext() *Context {
	indexes := new(Indexing)
	indexes.Reset()
	ctx := &Context{
		DatabaseName: DefaultDatabaseName,
		indexes:      indexes,
	}
	return ctx
}

func (ctx *Context) AddMatchWithOutcome(outcome *Outcome) {
	match := NewMatch(outcome)
	ctx.AddMatch(match, outcome)
}

func (ctx *Context) AddMatch(match *Match, outcome *Outcome) {
	ctx.Matches = append(ctx.Matches, match)
	m, ok := ctx.indexes.outcomesMap[outcome.ID]
	if ok {
		outcome = m
		m.Occurrences++
		return
	}
	outcome.Occurrences++
	ctx.Outcomes = append(ctx.Outcomes, outcome)
	ctx.indexes.outcomesMap[outcome.ID] = outcome
	return
}

func (ctx *Context) AddTeam(team *Team) {
	_, ok := ctx.indexes.teamsMap[team.ID]
	if ok {
		return
	}
	ctx.Teams = append(ctx.Teams, team)
	ctx.indexes.teamsMap[team.ID] = team
	ctx.indexes.teamsPlayerMap[team.ID] = map[string]*Player{}

	for _, player := range team.Players {
		ctx.AddPlayer(player)
		ctx.indexes.teamsPlayerMap[team.ID][player.ID] = player
		ctx.indexes.playersTeamMap[player.ID] = append(ctx.indexes.playersTeamMap[player.ID], team)
	}
	return
}

func (ctx *Context) AddPlayer(player *Player) {
	_, ok := ctx.indexes.playersMap[player.ID]
	if ok {
		return
	}
	ctx.Players = append(ctx.Players, player)
	ctx.indexes.playersMap[player.ID] = player
	ctx.indexes.playersNameMap[player.Name] = player
	return
}

func (ctx *Context) Store() error {
	if Verbose {
		log.Println("storing database")
	}
	err := StoreDB(ctx, *compress)
	if err != nil {
		log.Printf("error storing database: %s", err)
		return err
	}
	return nil
}

func (ctx *Context) Load() error {
	var bar *pb.ProgressBar
	if Verbose {
		log.Println("loading database")
	}
	db, err := LoadDB(ctx.DatabaseName, *decompress)
	if os.IsNotExist(err) {
		log.Printf("database not found")
		return nil
	} else if err != nil {
		log.Printf("error loading database: %s", err)
		return err
	}

	if Verbose {
		log.Println("loading teams")
		bar = pb.StartNew(len(db.Teams))
	}
	for _, team := range db.Teams {
		ctx.AddTeam(team)
		if Verbose {
			bar.Increment()
		}
	}
	if Verbose {
		bar.Finish()
	}

	if Verbose {
		log.Println("loading match history")
		bar = pb.StartNew(len(db.Outcomes))
	}
	loadedOutcomes := map[string]*Outcome{}
	for _, outcome := range db.Outcomes {
		loadedOutcomes[outcome.ID] = outcome
		if Verbose {
			bar.Increment()
		}
	}
	if Verbose {
		bar.Finish()
	}

	if Verbose {
		bar = pb.StartNew(len(db.Matches))
	}
	for _, match := range db.Matches {
		outcome, ok := loadedOutcomes[match.OutcomeID]
		if !ok {
			log.Panicf("corrupted history %q", match.ID)
		}
		ctx.AddMatch(match, outcome)
		if Verbose {
			bar.Increment()
		}
	}
	if Verbose {
		bar.Finish()
	}
	return nil
}
