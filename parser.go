package main

import (
	"fmt"
	"strconv"
)

// İfade öncelikleri
const (
	_ = iota
	LOWEST
	SUM      // + veya -
	PRODUCT  // * veya /
)

var precedences = map[TokenType]int{
	PLUS:     SUM,
	MINUS:    SUM,
	ASTERISK: PRODUCT,
	SLASH:    PRODUCT,
}

// AST Düğümleri
type Node interface {
	TokenLiteral() string
}

type Expression interface {
	Node
	expressionNode()
}

type IntegerLiteral struct {
	Token Token
	Value int
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

type InfixExpression struct {
	Token    Token // Operatör (+, -, *, /)
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// Parser Yapısı
type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Pratt Parsing: İfade önceliğini hesaplar
func (p *Parser) ParseExpression(precedence int) Expression {
	left := &IntegerLiteral{Token: p.curToken}
	val, _ := strconv.Atoi(p.curToken.Literal)
	left.Value = val

	for p.peekToken.Type != EOF && precedence < p.peekPrecedence() {
		p.nextToken()
		left = &InfixExpression{
			Token:    p.curToken,
			Left:     left,
			Operator: p.curToken.Literal,
			Right:    p.parseInfix(p.curPrecedence()),
		}
	}
	return left
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok { return p }
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok { return p }
	return LOWEST
}

func (p *Parser) parseInfix(precedence int) Expression {
	p.nextToken()
	return p.ParseExpression(precedence)
}