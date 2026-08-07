package parser

// Tokenize turns formula text into a token stream, always terminated
// by a single TokenEOF. String literals and unrecognized-character
// errors are handled in a later ticket (FASE-2.2).
func Tokenize(input string) ([]Token, error) {
	runes := []rune(input)
	n := len(runes)

	var tokens []Token
	line, col := 1, 1
	i := 0

	advance := func() rune {
		r := runes[i]
		i++
		if r == '\n' {
			line++
			col = 1
		} else {
			col++
		}
		return r
	}

	for i < n {
		startLine, startCol := line, col
		r := runes[i]

		switch {
		case r == ' ' || r == '\t' || r == '\n' || r == '\r':
			advance()

		case isDigit(r):
			start := i
			for i < n && (isDigit(runes[i]) || runes[i] == '.') {
				advance()
			}
			tokens = append(tokens, Token{Type: TokenNumber, Value: string(runes[start:i]), Line: startLine, Column: startCol})

		case isIdentStart(r):
			start := i
			for i < n && isIdentPart(runes[i]) {
				advance()
			}
			tokens = append(tokens, Token{Type: TokenIdent, Value: string(runes[start:i]), Line: startLine, Column: startCol})

		default:
			tt, ok := singleCharTokens[r]
			if !ok {
				// Unrecognized character: silently skipped for now.
				// FASE-2.2 replaces this with a proper SyntaxError.
				advance()
				continue
			}
			advance()
			tokens = append(tokens, Token{Type: tt, Value: string(r), Line: startLine, Column: startCol})
		}
	}

	tokens = append(tokens, Token{Type: TokenEOF, Line: line, Column: col})
	return tokens, nil
}

var singleCharTokens = map[rune]TokenType{
	'+': TokenPlus,
	'-': TokenMinus,
	'*': TokenStar,
	'/': TokenSlash,
	'(': TokenLParen,
	')': TokenRParen,
	',': TokenComma,
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isIdentStart(r rune) bool {
	return r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isIdentPart(r rune) bool {
	return isIdentStart(r) || isDigit(r)
}
