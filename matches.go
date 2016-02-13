package foosbot

import (
	"fmt"
	"strings"
)

type Match struct {
	ID       string  `json:"match_id"`
	WinnerID string  `json:"winner_id"`
	Teams    []*Team `json:"teams"`
	N        int     `json:"-"`
}

func (m *Match) ShortID() string {
	return strings.ToUpper(m.ID[:8])
}

func (m *Match) String() string {
	return fmt.Sprintf("Match %q", m.ShortID())
}

func (m *Match) Winner() *Team {
	for _, team := range m.Teams {
		if team.ID == m.WinnerID {
			return team
		}
	}
	return nil
}

func (m *Match) Looser() *Team {
	for _, team := range m.Teams {
		if team.ID != m.WinnerID {
			return team
		}
	}
	return nil
}

func buildMatchId(a, b *Team) string {
	ids := []string{a.ID, b.ID}
	return hash(ids...)
}

func NewMatch(winner, looser *Team) *Match {
	matchID := buildMatchId(winner, looser)
	match := new(Match)
	match.ID = matchID
	match.WinnerID = winner.ID
	match.Teams = []*Team{winner, looser}
	return match
}

func AddMatch(match *Match, entry *HistoryEntry) {
	addHistoryEntry(entry)
	m, ok := MatchesMap[match.ID]
	if ok {
		match = m
		m.N++
		return
	}
	match.N++
	Matches = append(Matches, match)
	MatchesMap[match.ID] = match
	for _, team := range match.Teams {
		AddTeam(team)
	}
	return
}

func AddMatchWithHistory(match *Match) {
	entry := NewHistoryEntry(match)
	AddMatch(match, entry)
}

func MatchByID(id string) (*Match, bool) {
	m, ok := MatchesMap[id]
	return m, ok
}

func MatchByTeams(a, b *Team) (match *Match, ok bool) {
	teamID := buildMatchId(a, b)
	match, ok = MatchesMap[teamID]
	return
}

func MatchesWithTeam(t *Team) (matches []*Match, history []*HistoryEntry) {
	outcomes := []string{}
	for _, match := range Matches {
		for _, team := range match.Teams {
			if t.ID == team.ID {
				outcomes = append(outcomes, match.ID)
				break
			}
		}
	}
	for _, entry := range History {
		if in(outcomes, entry.MatchID) {
			m, _ := MatchByID(entry.MatchID)
			matches = append(matches, m)
			history = append(history, entry)
		}
	}
	return matches, history
}

func in(arr []string, m string) bool {
	for i := range arr {
		if arr[i] == m {
			return true
		}
	}
	return false
}
