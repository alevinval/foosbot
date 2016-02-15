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

func buildMatchId(a, b *Team) string {
	ids := []string{a.ID, b.ID}
	return hash(ids...)
}

func NewMatch(winner, looser *Team) *Match {
	match := new(Match)
	match.ID = buildMatchId(winner, looser)
	match.WinnerID = winner.ID
	match.Teams = []*Team{winner, looser}
	return match
}

func (m *Match) ShortID() string {
	return strings.ToUpper(m.ID[:8])
}

func (m *Match) String() string {
	return fmt.Sprintf("Match %q", m.ShortID())
}

func (m *Match) IsWinner(t *Team) bool {
	return m.WinnerID == t.ID
}

func (m *Match) Winner() *Team {
	for _, team := range m.Teams {
		if team.ID == m.WinnerID {
			return team
		}
	}
	return nil
}

func (m *Match) Loosers() []*Team {
	loosers := []*Team{}
	for _, team := range m.Teams {
		if team.ID != m.WinnerID {
			loosers = append(loosers, team)
		}
	}
	return loosers
}

func (m *Match) IsLooser(t *Team) bool {
	if m.WinnerID == t.ID {
		return false
	}
	for i := range m.Teams {
		if t.ID == m.Teams[i].ID {
			return true
		}
	}
	return false
}
