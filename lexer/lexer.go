package lexer

import (
	"myRVCC/logger"
	"myRVCC/token"
	"myRVCC/utils"
	"os"
	"strings"
	"text/scanner"
	"unicode"
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
	case '*':
		tok = newToken(token.MUL, ch, position)
	case '/':
		tok = newToken(token.DIV, ch, position)
	case '(':
		tok = newToken(token.LPAREN, ch, position)
	case ')':
		tok = newToken(token.RPAREN, ch, position)
	case '{':
		tok = newToken(token.LBRACE, ch, position)
	case '}':
		tok = newToken(token.RBRACE, ch, position)
	case '=':
		l.Next()
		if l.Peek() == '=' {
			tok.Literal = "=="
			tok.Kind = token.EQ
		} else {
			tok = newToken(token.ASSIGN, ch, position)
			return tok
		}
	case '!':
		l.Next()
		if l.Peek() == '=' {
			tok.Literal = "!="
			tok.Kind = token.NEQ
		} else {
			//todo
			logger.Panic("[%s] unexpected character", ch)
		}
	case '>':
		l.Next()
		if l.Peek() == '=' {
			tok.Literal = ">="
			tok.Kind = token.GEQ
		} else {
			tok = newToken(token.GT, ch, position)
			return tok
		}
	case '<':
		l.Next()
		if l.Peek() == '=' {
			tok.Literal = "<="
			tok.Kind = token.LEQ
		} else {
			tok = newToken(token.LT, ch, position)
			return tok
		}
	case ';':
		tok = newToken(token.SEMICOLON, ch, position)
	case scanner.EOF:
		tok.Kind = token.EOF
		tok.Literal = ""
	default:
		if isIdentRune(ch, 0) {
			tok.Literal = l.readIdentifier()
			tok.Kind = token.LookUpIdent(tok.Literal)
			return tok
		} else if isDigit(ch) {
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

func (l *Lexer) readIdentifier() string {
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

// 判断是否是标识符.i 为0时，表示第一个必须不为数字
func isIdentRune(ch rune, i int) bool {
	return ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
}
