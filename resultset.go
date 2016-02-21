package foosbot

const (
	StatusWon  = MatchStatus(true)
	StatusLost = MatchStatus(false)
)

type MatchStatus bool

func (s MatchStatus) String() string {
	if s {
		return "won"
	} else {
		return "lost"
	}
}

type MatchResult struct {
	Match    *Match
	Status   MatchStatus
	Outcome  *Outcome
	Team     *Team
	Opponent *Team
}
