package foosbot

type Match struct {
	ID     string  `json:"match_id"`
	Winner string  `json:"winner_id"`
	Teams  []*Team `json:"teams"`
	N      int     `json:"n"`
}

var (
	matches    = []*Match{}
	matchesMap = map[string]*Match{}
)

func (m *Match) ShortID() string {
	return m.ID[:8]
}

func buildMatchId(a, b *Team) string {
	ids := []string{a.ID, b.ID}
	return hash(ids...)
}

func NewMatch(winner, looser *Team) *Match {
	matchID := buildMatchId(winner, looser)
	match := new(Match)
	match.ID = matchID
	match.Winner = winner.ID
	match.Teams = []*Team{winner, looser}
	match.N = 1
	return match
}

func AddMatch(match *Match) {
	m, ok := matchesMap[match.ID]
	if ok {
		m.N++
		return
	}
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
	for _, match := range matches {
		for _, team := range match.Teams {
			if t.ID == team.ID {
				foundMatches = append(foundMatches, match)
				break
			}
		}
	}
	return
}
