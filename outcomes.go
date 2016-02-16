package foosbot

import "errors"

var (
	ErrOutcomeImpossible = errors.New("an outcome cannot have the same winner and looser")
)

type Outcome struct {
	ID          string `json:"match_id"`
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

func (m *Outcome) ShortID() string {
	return m.ID[:8]
}

func (m *Outcome) IsWinner(t *Team) bool {
	return m.WinnerID == t.ID
}

func (m *Outcome) IsLooser(t *Team) bool {
	return m.LooserID == t.ID
}

func BuildOutcomeID(winner, looser *Team) string {
	ids := []string{winner.ID, looser.ID}
	return hash(ids...)
}
