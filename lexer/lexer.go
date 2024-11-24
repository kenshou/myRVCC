package lexer

import (
	"myRVCC/token"
	"strings"
	"text/scanner"
)

type Lexer struct {
	scanner.Scanner
}

func New(input string) *Lexer {
	//转换input到io.Reader
	reader := strings.NewReader(input)
	s := scanner.Scanner{}
	s.Init(reader)
	return &Lexer{s}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	ch := l.Peek()
	//跳过空白符号
	for l.Whitespace&(1<<uint(ch)) != 0 {
		l.Next()
		ch = l.Peek()
	}
	//保存当前位置
	position := l.Pos()
	tok.Position = position
	switch ch {
	case '+':
		tok = newToken(token.ADD, ch, position)
	case '-':
		tok = newToken(token.SUB, ch, position)
	case scanner.EOF:
		tok.Kind = token.EOF
		tok.Literal = ""
	default:
		if isDigit(ch) {
			tok.Literal = l.readNumber()
			tok.Kind = token.INT
			return tok
		}
	}
	l.Next()
	return tok
}

func (l *Lexer) readNumber() string {
	l.Scanner.Scan()
	return l.TokenText()
}

func newToken(tokenKind token.TokenKind, ch rune, position scanner.Position) token.Token {
	return token.Token{
		Kind:     tokenKind,
		Literal:  string(ch),
		Position: position,
	}
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
