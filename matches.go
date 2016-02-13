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

var (
	head         = 0
	matchHistory = []string{}
	matches      = []*Match{}
	matchesMap   = map[string]*Match{}
)

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

func AddMatch(match *Match) {
	defer func() {
		matchHistory = append(matchHistory, match.ID)
	}()
	m, ok := matchesMap[match.ID]
	if ok {
		m.N++
		return
	}
	match.N++
	matches = append(matches, match)
	matchesMap[match.ID] = match
	for _, team := range match.Teams {
		AddTeam(team)
	}
	return
}

func MatchByID(id string) (*Match, bool) {
	m, ok := matchesMap[id]
	return m, ok
}

func MatchByTeams(a, b *Team) (match *Match, ok bool) {
	teamID := buildMatchId(a, b)
	match, ok = matchesMap[teamID]
	return
}

func MatchesWithTeam(t *Team) (foundMatches []*Match) {
	outcomes := []string{}
	for _, match := range matches {
		for _, team := range match.Teams {
			if t.ID == team.ID {
				outcomes = append(outcomes, match.ID)
				break
			}
		}
	}
	for _, matchID := range matchHistory {
		if in(outcomes, matchID) {
			m, _ := MatchByID(matchID)
			foundMatches = append(foundMatches, m)
		}
	}
	return
}

func in(arr []string, m string) bool {
	for i := range arr {
		if arr[i] == m {
			return true
		}
	}
	return false
}
