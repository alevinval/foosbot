package parsing

type tokenType int

const (
	TokenEOF tokenType = iota
	TokenIllegal
	TokenWhitespace
	TokenIdentifier
	TokenDigit
	TokenCommandMatch
	TokenCommandStats
	TokenKeywordFoosbot
	TokenKeywordVS
)

func (t tokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
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
	case TokenIllegal:
		return "illegal token"
	default:
		return "unkown token"
	}
}
