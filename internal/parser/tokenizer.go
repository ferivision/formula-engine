package parser

// Tokenize turns formula text into a token stream, always terminated
// by a single TokenEOF.
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

		case r == '"':
			advance()
			var value []rune
			closed := false
			for i < n {
				c := runes[i]
				if c == '"' {
					advance()
					closed = true
					break
				}
				if c == '\\' && i+1 < n {
					advance()
					value = append(value, runes[i])
					advance()
					continue
				}
				value = append(value, c)
				advance()
			}
			if !closed {
				return nil, newSyntaxError("unterminated string literal", startLine, startCol)
			}
			tokens = append(tokens, Token{Type: TokenString, Value: string(value), Line: startLine, Column: startCol})

		default:
			tt, ok := singleCharTokens[r]
			if !ok {
				return nil, newSyntaxError("unrecognized character "+string(r), startLine, startCol)
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
