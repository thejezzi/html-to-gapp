package htg

import (
	"strings"
)

type Lexer struct {
	source       string
	tokens       []Token
	currentToken int
	line         int
	pos          int
	start        int
	current      int
}

func (lexer *Lexer) ScanTokens() []Token {
	for !lexer.isAtEnd() {

		lexer.start = lexer.current
		lexer.scanToken()
	}
	lexer.tokens = append(lexer.tokens, Token{"EOF", ""})
	return lexer.tokens
}

func (lexer *Lexer) isAtEnd() bool {
	return lexer.current >= len(lexer.source)
}

func (lexer *Lexer) scanToken() {
	switch c := lexer.advance(); c {
	case " ":
	case "\r":
	case "\t":
	case "\n":
		lexer.line++
		lexer.pos = 0
	// html comment
	case "<":
		if lexer.peek() == "!" && lexer.peekNext() == "-" && lexer.peekNextTwo() == "-" {
			lexer.Comment()
			break
		}
		if lexer.peek() == "/" {
			lexer.advance()
			lexer.addToken("TAG_END_CLOSE")
			break
		}
		lexer.addToken("TAG_START")
	case ">":
		lexer.addToken("TAG_END")
	case "/":
		if lexer.match(">") {
			lexer.addToken("TAG_END_SELF_CLOSE")
		}
		htg.Error(lexer.line, lexer.pos, "Unexpected character.", "Unexpected character '/'")
	case "=":
		lexer.addToken("EQUALS")
	case "\"":
		lexer.string()
	case "'":
		lexer.string()
	default:
		if lexer.isAlphaNumeric(c) {
			lexer.identifier()
			break
		}
		htg.Error(lexer.line, lexer.pos, "Unexpected character.", "Unexpected character '"+c+"'")
	}
}

func (lexer *Lexer) advance() string {
	lexer.current++
	lexer.pos++
	return string(lexer.source[lexer.current-1])
}

func (lexer *Lexer) advanceToken() {
	lexer.start = lexer.current
}

func (lexer *Lexer) addToken(tokenType string) {
	text := lexer.source[lexer.start:lexer.current]
	lexer.addTokenWithLiteral(tokenType, text)
}

func (lexer *Lexer) addTokenWithLiteral(tokenType string, literal string) {
	lexer.tokens = append(lexer.tokens, Token{tokenType, literal})
}

func (lexer *Lexer) peek() string {
	if lexer.isAtEnd() {
		return "\000"
	}
	return string(lexer.source[lexer.current])
}

func (lexer *Lexer) peekNext() string {
	if lexer.current+1 >= len(lexer.source) {
		return "\000"
	}
	return string(lexer.source[lexer.current+1])
}

func (lexer *Lexer) peekNextTwo() string {
	if lexer.current+2 >= len(lexer.source) {
		return "\000"
	}
	return string(lexer.source[lexer.current+2])
}

func (lexer *Lexer) match(expected string) bool {
	if lexer.isAtEnd() {
		return false
	}
	if lexer.source[lexer.current] != expected[0] {
		return false
	}
	lexer.current++
	return true
}

func (lexer *Lexer) Comment() {

	// skip "<!--"
	for i := 0; i < 3; i++ {
		lexer.advance()
	}
	for lexer.peek() != "-" && lexer.peekNext() != "-" && lexer.peekNextTwo() != ">" && !lexer.isAtEnd() {
		lexer.advance()
	}

	for i := 0; i < 4; i++ {
		lexer.advance()
	}
}

func (lexer *Lexer) string() {
	for lexer.isAlphaNumeric(lexer.peek()) {
		lexer.advance()
	}
	// skip closing quote
	lexer.advance()

	value := lexer.source[lexer.start+1 : lexer.current-1]
	lexer.addTokenWithLiteral("STRING", value)
}

func (lexer *Lexer) isAlphaNumeric(c string) bool {
	return strings.Contains("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", c)
}

func (lexer *Lexer) isHyphen(c string) bool {
	return strings.Contains("-", c)
}

func (lexer *Lexer) isUnderscore(c string) bool {
	return strings.Contains("_", c)
}

func (lexer *Lexer) identifier() {
	for lexer.isAlphaNumeric(lexer.peek()) || lexer.isHyphen(lexer.peek()) || lexer.isUnderscore(lexer.peek()) {
		lexer.advance()
	}
	lexer.addToken("IDENTIFIER")
}
