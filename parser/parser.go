package parser

import (
	"fmt"
	"myRVCC/ast"
	"myRVCC/lexer"
	"myRVCC/logger"
	"myRVCC/token"
	"strconv"
)

const (
	//运算符优先级等级
	_ = iota
	LOWEST
	EQUALS       // ==
	LESS_GREATER // > or <
	SUM          // + -
	PRODUCT      // * /
	PREFIX       // -X or !X
	CALL         // 函数调用,括号
)

// tokenKind与运算符优先级等级的对应关系
var precedences = map[token.TokenKind]int{
	token.ADD:    SUM,
	token.SUB:    SUM,
	token.MUL:    PRODUCT,
	token.DIV:    PRODUCT,
	token.LPAREN: CALL,
	token.EQ:     EQUALS,
	token.NEQ:    EQUALS,
	token.LT:     LESS_GREATER,
	token.LEQ:    LESS_GREATER,
	token.GT:     LESS_GREATER,
	token.GEQ:    LESS_GREATER,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenKind]prefixParseFn
	infixParseFns  map[token.TokenKind]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.prefixParseFns = make(map[token.TokenKind]prefixParseFn)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupExpression)
	p.registerPrefix(token.SUB, p.parsePrefixExpression)
	p.registerPrefix(token.ADD, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenKind]infixParseFn)
	p.registerInfix(token.ADD, p.parseInfixExpression)
	p.registerInfix(token.SUB, p.parseInfixExpression)
	p.registerInfix(token.MUL, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LEQ, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GEQ, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerPrefix(kind token.TokenKind, fn prefixParseFn) {
	p.prefixParseFns[kind] = fn
}
func (p *Parser) registerInfix(kind token.TokenKind, fn infixParseFn) {
	p.infixParseFns[kind] = fn
}
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
func (p *Parser) curTokenIs(t token.TokenKind) bool {
	return p.curToken.Kind == t
}

func (p *Parser) peekTokenIs(t token.TokenKind) bool {
	return p.peekToken.Kind == t
}

func (p *Parser) expectPeek(t token.TokenKind) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		//todo
		logger.Panic("[%s] expectPeek error", p.curToken.Literal)
		return false
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Kind != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Kind {
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Kind]
	if prefix == nil {
		//todo
		logger.Panic("no prefix parse function for %s found", p.curToken.Kind)
		return nil
	}
	leftExp := prefix()
	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Kind]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if pre, ok := precedences[p.peekToken.Kind]; ok {
		return pre
	}
	return LOWEST
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		//todo
		logger.Panic(msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) curPrecedence() int {
	if pre, ok := precedences[p.curToken.Kind]; ok {
		return pre
	}
	return LOWEST
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}
