package ast

import (
	"bytes"
	"myRVCC/logger"
	"myRVCC/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}
type IdentifierObj struct {
	Value  string //变量名
	Offset int64  //fp的偏移量
}
type Env struct {
	IdentObjArr []*IdentifierObj
	StackSize   int64
	Parent      *Env
}

func CreateEnv(parent *Env) *Env {
	return &Env{
		Parent: parent,
	}
}
func (e *Env) FindOrCreateIdentifier(ident *token.Token) *IdentifierObj {
	if ident.Kind != token.IDENT {
		logger.Panic("[%s] FindOrCreateIdentifier: token is not int", ident.Literal)
	}
	for _, obj := range e.IdentObjArr {
		if obj.Value == ident.Literal {
			return obj
		}
	}
	obj := &IdentifierObj{
		Value: ident.Literal,
	}
	e.IdentObjArr = append(e.IdentObjArr, obj)
	return obj
}

// Statement 语句
type Statement interface {
	Node
	statementNode()
}

// Expression 表达式
type Expression interface {
	Node
	expressionNode()
}

// Program 程序
type Program struct {
	Env        *Env
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}

	return out.String()
}

// ExpressionStatement 表达式语句
type ExpressionStatement struct {
	Token      token.Token // 第一个token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral 整形字面量
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) expressionNode() {}

// Identifier 标识符
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
	Obj   *IdentifierObj
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// InfixExpression / 中缀表达式
type InfixExpression struct {
	Token    token.Token // operator当前运算符号，比如 * +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type BlockStatement struct {
	Token      token.Token // {
	Statements []Statement
	Env        *Env
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
