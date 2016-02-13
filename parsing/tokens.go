package parsing

type Token int

const (
	TokenEOF Token = iota
	TokenWhitespace
	TokenIdentifier
	TokenDigit
	TokenCommandMatch
	TokenCommandStats
	TokenKeywordFoosbot
	TokenKeywordVS
)

func (t Token) String() string {
	switch t {
	case TokenEOF:
		return "eof"
	case TokenWhitespace:
		return "whitespace"
	case TokenIdentifier:
		return "identifier"
	case TokenDigit:
		return "digit"
	case TokenCommandMatch:
		return "match"
	case TokenCommandStats:
		return "stats"
	case TokenKeywordFoosbot:
		return "foosbot"
	case TokenKeywordVS:
		return "vs"
	default:
		return "-unkown-"
	}
}
