package foosbot

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"strings"
)

func namesFromTeam(t *Team) string {
	names := []string{}
	for _, player := range t.Players {
		names = append(names, player.Name)
	}
	return strings.Join(names, ",")
}

func (ctx *Context) Print(i interface{}) (out string) {
	switch obj := i.(type) {
	case *Match:
		out = fmt.Sprintf("%s (%s)", obj.ShortID(), humanize.Time(obj.PlayedAt))
	case *Outcome:
		w, _ := ctx.Query.TeamByID(obj.WinnerID)
		l, _ := ctx.Query.TeamByID(obj.LooserID)
		out = fmt.Sprintf("(%s) vs (%s)", namesFromTeam(w), namesFromTeam(l))
	case *Team:
		out = fmt.Sprintf("%s (%s)", obj.ShortID(), namesFromTeam(obj))
	case *Player:
		out = fmt.Sprintf("%s (%s)", obj.ShortID(), obj.Name)
	default:
		b, _ := json.Marshal(obj)
		out = string(b)
	}
	return
}
