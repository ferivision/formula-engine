package parser

// Parse turns a token stream into an AST using precedence climbing:
// expression := term (('+' | '-') term)*
// term       := factor (('*' | '/') factor)*
// factor     := NUMBER | STRING | IDENT | '(' expression ')'
//
// Function-call syntax and detailed syntax-error messages are handled
// in later tickets (FASE-3.2, FASE-3.3).
func Parse(tokens []Token) (*Node, error) {
	p := &astParser{tokens: tokens}
	node, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if p.current().Type != TokenEOF {
		return nil, newSyntaxError("unexpected token", p.current().Line, p.current().Column)
	}
	return node, nil
}

type astParser struct {
	tokens []Token
	pos    int
}

func (p *astParser) current() Token {
	return p.tokens[p.pos]
}

func (p *astParser) advance() Token {
	tok := p.tokens[p.pos]
	if p.pos < len(p.tokens)-1 {
		p.pos++
	}
	return tok
}

func (p *astParser) parseExpression() (*Node, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for p.current().Type == TokenPlus || p.current().Type == TokenMinus {
		op := p.advance()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &Node{Type: NodeBinary, Operator: op.Type, Left: left, Right: right}
	}
	return left, nil
}

func (p *astParser) parseTerm() (*Node, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for p.current().Type == TokenStar || p.current().Type == TokenSlash {
		op := p.advance()
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		left = &Node{Type: NodeBinary, Operator: op.Type, Left: left, Right: right}
	}
	return left, nil
}

func (p *astParser) parseFactor() (*Node, error) {
	tok := p.current()
	switch tok.Type {
	case TokenNumber:
		p.advance()
		return &Node{Type: NodeNumber, Value: tok.Value}, nil
	case TokenString:
		p.advance()
		return &Node{Type: NodeString, Value: tok.Value}, nil
	case TokenIdent:
		p.advance()
		return &Node{Type: NodeIdentifier, Value: tok.Value}, nil
	case TokenLParen:
		p.advance()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.current().Type != TokenRParen {
			return nil, newSyntaxError("expected closing parenthesis", p.current().Line, p.current().Column)
		}
		p.advance()
		return expr, nil
	default:
		return nil, newSyntaxError("unexpected token", tok.Line, tok.Column)
	}
}
