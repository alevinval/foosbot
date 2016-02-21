package foosbot

import (
	"time"
)

type Match struct {
	ID        string    `json:"match_id"`
	OutcomeID string    `json:"outcome_id"`
	PlayedAt  time.Time `json:"played_at"`
}

func NewMatch(outcome *Outcome) (match *Match) {
	now := time.Now()
	match = &Match{
		ID:        BuildMatchID(outcome, now),
		OutcomeID: outcome.ID,
		PlayedAt:  now,
	}
	return
}

func (m *Match) ShortID() string {
	return m.ID[:8]
}

func BuildMatchID(outcome *Outcome, playedAt time.Time) string {
	return hash(outcome.ID, playedAt.String())
}
