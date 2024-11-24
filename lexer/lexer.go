package lexer

import (
	"myRVCC/token"
	"myRVCC/utils"
	"os"
	"strings"
	"text/scanner"
)

type Lexer struct {
	scanner.Scanner
	code string
}

func New(input string) *Lexer {
	//转换input到io.Reader
	reader := strings.NewReader(input)
	s := scanner.Scanner{}
	s.Init(reader)
	return &Lexer{s, input}
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
		l.Error(position, "unexpected character: %s", string(ch))
		os.Exit(1)
	}
	l.Next()
	return tok
}

func (l *Lexer) readNumber() string {
	l.Scanner.Scan()
	return l.TokenText()
}

func (l *Lexer) Error(pos scanner.Position, format string, arg ...interface{}) {
	utils.VErrorAt(l.code, pos, format, arg...)
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