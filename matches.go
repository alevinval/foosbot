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

func (c *Context) AddMatch(match *Match, entry *HistoryEntry) {
	c.addHistoryEntry(entry)
	m, ok := c.MatchesMap[match.ID]
	if ok {
		match = m
		m.N++
		return
	}
	match.N++
	c.Matches = append(c.Matches, match)
	c.MatchesMap[match.ID] = match
	for _, team := range match.Teams {
		c.AddTeam(team)
	}
	return
}

func (c *Context) AddMatchWithHistory(match *Match) {
	entry := NewHistoryEntry(match)
	c.AddMatch(match, entry)
}

func (c *Context) MatchByID(id string) (*Match, bool) {
	m, ok := c.MatchesMap[id]
	return m, ok
}

func (c *Context) MatchByTeams(a, b *Team) (match *Match, ok bool) {
	teamID := buildMatchId(a, b)
	match, ok = c.MatchesMap[teamID]
	return
}

func (c *Context) MatchesWithTeam(t *Team) (matches []*Match, history []*HistoryEntry) {
	outcomes := []string{}
	for _, match := range c.Matches {
		for _, team := range match.Teams {
			if t.ID == team.ID {
				outcomes = append(outcomes, match.ID)
				break
			}
		}
	}
	for _, entry := range c.History {
		if in(outcomes, entry.MatchID) {
			m, _ := c.MatchByID(entry.MatchID)
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
