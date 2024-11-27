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
	ASSIGN       //赋值 =
	EQUALS       // ==
	LESS_GREATER // > or <
	SUM          // + -
	PRODUCT      // * /
	PREFIX       // -X or !X
	CALL         // 函数调用,括号
)

// tokenKind与运算符优先级等级的对应关系
var precedences = map[token.TokenKind]int{
	token.ASSIGN: ASSIGN,
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
	prefixParseFn func(env *ast.Env) ast.Expression
	infixParseFn  func(ast.Expression, *ast.Env) ast.Expression
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
	p.registerPrefix(token.IDENT, p.parseIdentifierExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	// * 解引用
	p.registerPrefix(token.MUL, p.parsePrefixExpression)
	// & 取地址
	p.registerPrefix(token.AMPERSAND, p.parsePrefixExpression)

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
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)

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
	env := ast.CreateEnv(nil)
	program := &ast.Program{
		Env: env,
	}
	program.Statements = []ast.Statement{}
	for p.curToken.Kind != token.EOF {
		if p.curTokenIs(token.COMMENT) || p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement(env)
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement(env *ast.Env) ast.Statement {
	switch p.curToken.Kind {
	case token.RETURN:
		return p.parseReturnStatement(env)
	case token.LBRACE:
		return p.parseBlockStatement(env)
	case token.FOR:
		return p.parseForStatement(env)
	case token.WHILE:
		return p.parseWhileStatement(env)
	default:
		return p.parseExpressionStatement(env)
	}
}

func (p *Parser) parseExpressionStatement(env *ast.Env) *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST, env)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int, env *ast.Env) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Kind]
	if prefix == nil {
		//todo
		logger.Panic("no prefix parse function for %s found", p.curToken.Kind)
		return nil
	}
	leftExp := prefix(env)
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Kind]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp, env)
	}
	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if pre, ok := precedences[p.peekToken.Kind]; ok {
		return pre
	}
	return LOWEST
}

func (p *Parser) parseIntegerLiteral(env *ast.Env) ast.Expression {
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

func (p *Parser) parseInfixExpression(left ast.Expression, env *ast.Env) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	//等号为右结合，所以减1
	if p.curToken.Kind == token.ASSIGN {
		precedence--
	}

	p.nextToken()
	expression.Right = p.parseExpression(precedence, env)
	return expression
}

func (p *Parser) curPrecedence() int {
	if pre, ok := precedences[p.curToken.Kind]; ok {
		return pre
	}
	return LOWEST
}

func (p *Parser) parseGroupExpression(env *ast.Env) ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST, env)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parsePrefixExpression(env *ast.Env) ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX, env)
	return expression
}

func (p *Parser) parseIdentifierExpression(env *ast.Env) ast.Expression {
	ident := env.FindOrCreateIdentifier(&p.curToken)
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
		Obj:   ident,
	}
}

func (p *Parser) parseReturnStatement(env *ast.Env) ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST, env)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseBlockStatement(env *ast.Env) *ast.BlockStatement {
	stmt := &ast.BlockStatement{Token: p.curToken}
	//todo 暂时还不切换env
	stmt.Env = env //ast.CreateEnv(env)
	stmt.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.COMMENT) || p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
			continue
		}
		stmt.Statements = append(stmt.Statements, p.parseStatement(stmt.Env))
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIfExpression(env *ast.Env) ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		//todo
		logger.Panic("[%s] parseIfExpression error", p.curToken.Literal)
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST, env)
	if !p.expectPeek(token.RPAREN) {
		logger.Panic("[%s] parseIfExpression error", p.curToken.Literal)
		return nil
	}
	p.nextToken()
	expression.Consequence = p.parseStatement(env)
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		p.nextToken()
		expression.Alternative = p.parseStatement(env)
	}
	return expression
}

func (p *Parser) parseForStatement(env *ast.Env) ast.Statement {
	stmt := &ast.ForStatement{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		//todo
		logger.Panic("[%s] parseForStatement error", p.curToken.Literal)
	}
	p.nextToken()
	if !p.curTokenIs(token.SEMICOLON) {
		stmt.InitExpr = p.parseExpression(LOWEST, env)
		if !p.expectPeek(token.SEMICOLON) {
			logger.Panic("[%s] parseForStatement error", p.curToken.Literal)
			return nil
		}
	}
	p.nextToken()
	if !p.curTokenIs(token.SEMICOLON) {
		stmt.Condition = p.parseExpression(LOWEST, env)
		if !p.expectPeek(token.SEMICOLON) {
			logger.Panic("[%s] parseForStatement error", p.curToken.Literal)
			return nil
		}
	}
	p.nextToken()
	if !p.curTokenIs(token.RPAREN) {
		stmt.Inc = p.parseExpression(LOWEST, env)

		if !p.expectPeek(token.RPAREN) {
			logger.Panic("[%s] parseForStatement error", p.curToken.Literal)
			return nil
		}
	}
	p.nextToken()
	stmt.Consequence = p.parseStatement(env)
	return stmt
}

func (p *Parser) parseWhileStatement(env *ast.Env) ast.Statement {
	stmt := &ast.WhileStatement{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		logger.Panic("[%s] parseWhileStatement error", p.curToken.Literal)
	}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST, env)
	if !p.expectPeek(token.RPAREN) {
		logger.Panic("[%s] parseWhileStatement error", p.curToken.Literal)
	}
	p.nextToken()
	stmt.Consequence = p.parseStatement(env)
	return stmt
}
