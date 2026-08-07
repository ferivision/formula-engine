package parser

import "testing"

func TestTokenize_SkipsTrailingEOF(t *testing.T) {
	tokens, err := Tokenize("")
	if err != nil {
		t.Fatalf("Tokenize() error = %v", err)
	}
	if len(tokens) != 1 || tokens[0].Type != TokenEOF {
		t.Fatalf("Tokenize(\"\") = %+v, want single EOF token", tokens)
	}
}

func TestTokenize_TokenTypes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "integer",
			input: "42",
			want:  []Token{{Type: TokenNumber, Value: "42"}},
		},
		{
			name:  "decimal",
			input: "3.14",
			want:  []Token{{Type: TokenNumber, Value: "3.14"}},
		},
		{
			name:  "identifier",
			input: "subtotal",
			want:  []Token{{Type: TokenIdent, Value: "subtotal"}},
		},
		{
			name:  "identifier with digits and underscore",
			input: "field_2",
			want:  []Token{{Type: TokenIdent, Value: "field_2"}},
		},
		{
			name:  "operators",
			input: "+-*/",
			want: []Token{
				{Type: TokenPlus, Value: "+"},
				{Type: TokenMinus, Value: "-"},
				{Type: TokenStar, Value: "*"},
				{Type: TokenSlash, Value: "/"},
			},
		},
		{
			name:  "punctuation",
			input: "(,)",
			want: []Token{
				{Type: TokenLParen, Value: "("},
				{Type: TokenComma, Value: ","},
				{Type: TokenRParen, Value: ")"},
			},
		},
		{
			name:  "string literal",
			input: `"hello world"`,
			want:  []Token{{Type: TokenString, Value: "hello world"}},
		},
		{
			name:  "string literal with escaped quote",
			input: `"he said \"hi\""`,
			want:  []Token{{Type: TokenString, Value: `he said "hi"`}},
		},
		{
			name:  "nested function call expression",
			input: "MAX(1, MIN(2, 3.5) * x)",
			want: []Token{
				{Type: TokenIdent, Value: "MAX"},
				{Type: TokenLParen, Value: "("},
				{Type: TokenNumber, Value: "1"},
				{Type: TokenComma, Value: ","},
				{Type: TokenIdent, Value: "MIN"},
				{Type: TokenLParen, Value: "("},
				{Type: TokenNumber, Value: "2"},
				{Type: TokenComma, Value: ","},
				{Type: TokenNumber, Value: "3.5"},
				{Type: TokenRParen, Value: ")"},
				{Type: TokenStar, Value: "*"},
				{Type: TokenIdent, Value: "x"},
				{Type: TokenRParen, Value: ")"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.input)
			if err != nil {
				t.Fatalf("Tokenize(%q) error = %v", tt.input, err)
			}
			// last token is always EOF; strip it for comparison against `want`
			got = got[:len(got)-1]
			if len(got) != len(tt.want) {
				t.Fatalf("Tokenize(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
			for i := range tt.want {
				if got[i].Type != tt.want[i].Type || got[i].Value != tt.want[i].Value {
					t.Errorf("token[%d] = %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestTokenize_TracksLineAndColumn(t *testing.T) {
	got, err := Tokenize("1 +\n  22")
	if err != nil {
		t.Fatalf("Tokenize() error = %v", err)
	}

	want := []Token{
		{Type: TokenNumber, Value: "1", Line: 1, Column: 1},
		{Type: TokenPlus, Value: "+", Line: 1, Column: 3},
		{Type: TokenNumber, Value: "22", Line: 2, Column: 3},
	}

	for i, w := range want {
		if got[i].Line != w.Line || got[i].Column != w.Column {
			t.Errorf("token[%d] position = (line %d, col %d), want (line %d, col %d)",
				i, got[i].Line, got[i].Column, w.Line, w.Column)
		}
	}
}
