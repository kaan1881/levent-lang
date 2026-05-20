package main

type TokenType string

const (
	ILLEGAL  = "ILLEGAL"
	EOF      = "EOF"
	IDENT    = "IDENT"
	INT      = "INT"
	STRING   = "STRING"
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	GT       = ">"
	LT       = "<"
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	COMMA    = ","

	// Mevcutlar
	OLSUM   = "OLSUM"
	YAZDIR  = "YAZDIR"
	EKRAN   = "EKRAN"
	RENK    = "RENK"
	KARE    = "KARE"
	DAIRE   = "DAIRE"
	EGER    = "EGER"
	DONGU   = "DONGU"
	DAHILET = "DAHILET"

	// Yeni Eklenen Gelişmiş Özellikler
	SABIT     = "SABIT"
	FUNC      = "FUNC"
	GOPARALEL = "GOPARALEL"
	DEGILSE   = "DEGILSE"
	DONDUR    = "DONDUR"
	EQ        = "=="
	NE        = "!="
	AND       = "&&"
	OR        = "||"
)

type Token struct {
	Type    TokenType
	Literal string
}

var anahtarKelimeler = map[string]TokenType{
	"olsum":      OLSUM,
	"yazdir":     YAZDIR,
	"ekran":      EKRAN,
	"renk":       RENK,
	"kare":       KARE,
	"daire":      DAIRE,
	"eger":       EGER,
	"dongu":      DONGU,
	"dahilet":    DAHILET,
	"sabit":      SABIT,
	"func":       FUNC,
	"go_paralel": GOPARALEL,
	"degilse":    DEGILSE,
	"dondur":     DONDUR,
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = Token{ASSIGN, string(l.ch)}
	case '+':
		tok = Token{PLUS, string(l.ch)}
	case '-':
		tok = Token{MINUS, string(l.ch)}
	case '/':
		tok = Token{SLASH, string(l.ch)}
	case '*':
		tok = Token{ASTERISK, string(l.ch)}
	case '(':
		tok = Token{LPAREN, string(l.ch)}
	case ')':
		tok = Token{RPAREN, string(l.ch)}
	case '{':
		tok = Token{LBRACE, string(l.ch)}
	case '}':
		tok = Token{RBRACE, string(l.ch)}
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			if tokType, ok := anahtarKelimeler[tok.Literal]; ok {
				tok.Type = tokType
			} else {
				tok.Type = IDENT
			}
			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}