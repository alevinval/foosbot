package foosbot

import "time"

type HistoryEntry struct {
	ID       string    `json:"history_id"`
	MatchID  string    `json:"match_id"`
	PlayedAt time.Time `json:"played_at"`
}

func buildHistoryEntryId(m *Match, at time.Time) string {
	return hash(m.ID, at.String())
}

func NewHistoryEntry(m *Match) *HistoryEntry {
	now := time.Now()
	entry := new(HistoryEntry)
	entry.ID = buildHistoryEntryId(m, now)
	entry.MatchID = m.ID
	entry.PlayedAt = now
	return entry
}
