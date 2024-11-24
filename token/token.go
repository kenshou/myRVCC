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
	RPAREN //)

	ASSIGN //=

	SEMICOLON //;
	operator_end
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

	LPAREN:    "(",
	RPAREN:    ")",
	SEMICOLON: ";",
}

func (tk TokenKind) String() string {
	return tokens[tk]
}

type Token struct {
	Kind     TokenKind
	Literal  string
	Position scanner.Position
}
