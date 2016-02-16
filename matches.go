package foosbot

import (
	"time"
)

type Match struct {
	ID        string    `json:"history_id"`
	OutcomeID string    `json:"match_id"`
	PlayedAt  time.Time `json:"played_at"`
}

func NewMatch(outcome *Outcome) (entry *Match) {
	now := time.Now()
	entry = &Match{
		ID:        BuildHistoryID(outcome, now),
		OutcomeID: outcome.ID,
		PlayedAt:  now,
	}
	return
}

func (h *Match) ShortID() string {
	return h.ID[:8]
}

func BuildHistoryID(outcome *Outcome, playedAt time.Time) string {
	return hash(outcome.ID, playedAt.String())
}
