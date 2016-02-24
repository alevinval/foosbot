package foosbot

type (
	MatchStatus bool
	MatchResult struct {
		Match    *Match
		Outcome  *Outcome
		Team     *Team
		Opponent *Team
		Status   MatchStatus
	}
)

const (
	StatusWon  = MatchStatus(true)
	StatusLost = MatchStatus(false)
)

func (s MatchStatus) String() string {
	if s {
		return "won"
	} else {
		return "lost"
	}
}
