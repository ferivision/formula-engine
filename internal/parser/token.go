package parser

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNumber
	TokenIdent
	TokenString
	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash
	TokenLParen
	TokenRParen
	TokenComma
)

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenNumber:
		return "NUMBER"
	case TokenIdent:
		return "IDENT"
	case TokenString:
		return "STRING"
	case TokenPlus:
		return "PLUS"
	case TokenMinus:
		return "MINUS"
	case TokenStar:
		return "STAR"
	case TokenSlash:
		return "SLASH"
	case TokenLParen:
		return "LPAREN"
	case TokenRParen:
		return "RPAREN"
	case TokenComma:
		return "COMMA"
	default:
		return "UNKNOWN"
	}
}

// Token is one lexical unit produced by Tokenize, with its position in
// the original formula text (1-indexed) for error reporting.
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}
