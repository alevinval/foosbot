package foosbot

import "errors"

var (
	ErrOutcomeImpossible = errors.New("an outcome cannot have the same winner and looser")
)

type Outcome struct {
	ID          string `json:"outcome_id"`
	WinnerID    string `json:"winner_id"`
	LooserID    string `json:"looser_id"`
	Occurrences int    `json:"-"`
}

func NewOutcome(winner, looser *Team) (outcome *Outcome, err error) {
	if winner.ID == looser.ID {
		err = ErrOutcomeImpossible
		return
	}
	outcome = &Outcome{
		ID:          BuildOutcomeID(winner, looser),
		WinnerID:    winner.ID,
		LooserID:    looser.ID,
		Occurrences: 0,
	}
	return
}

func (o *Outcome) ShortID() string {
	return o.ID[:8]
}

func (o *Outcome) IsWinner(t *Team) bool {
	return o.WinnerID == t.ID
}

func (o *Outcome) IsLooser(t *Team) bool {
	return o.LooserID == t.ID
}

func BuildOutcomeID(winner, looser *Team) string {
	ids := []string{winner.ID, looser.ID}
	return hash(ids...)
}
