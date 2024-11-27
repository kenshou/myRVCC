package token

import (
	"text/scanner"
)

type TokenKind int

const (
	ILLEGAL TokenKind = iota //
	EOF                      //文件终止符，文件的最后
	COMMENT
	//字面量
	literal_begin
	IDENT //标识符/变量名/函数名
	INT   //int整形字面量
	literal_end

	//操作符号
	operator_begin
	ADD //+
	SUB //-
	MUL //*
	DIV // /

	EQ  //==
	NEQ //!=
	LT  // <
	LEQ // <=
	GT  // >
	GEQ // >=

	LPAREN //(
	LBRACE //{

	RPAREN //)
	RBRACE // }

	ASSIGN //=

	SEMICOLON //;
	operator_end

	keyword_begin
	RETURN //return
	IF
	ELSE
	FOR
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT: "IDENT",
	INT:   "INT",

	ADD:    "+",
	SUB:    "-",
	MUL:    "*",
	DIV:    "/",
	EQ:     "==",
	NEQ:    "!=",
	LT:     "<",
	LEQ:    "<=",
	GT:     ">",
	GEQ:    ">=",
	ASSIGN: "=",

	LPAREN: "(",
	LBRACE: "{",

	RPAREN:    ")",
	RBRACE:    "}",
	SEMICOLON: ";",

	RETURN: "return",
	IF:     "if",
	ELSE:   "else",
	FOR:    "for",
}
var keywords = map[string]TokenKind{}

func init() {
	for i, keyword := range tokens[keyword_begin+1 : keyword_end] {
		keywords[keyword] = TokenKind(i + int(keyword_begin) + 1)
	}
}
func LookUpIdent(ident string) TokenKind {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
func (tk TokenKind) String() string {
	return tokens[tk]
}

type Token struct {
	Kind     TokenKind
	Literal  string
	Position scanner.Position
}
