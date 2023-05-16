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
  insideTag    bool
}

func (lexer *Lexer) ScanTokens() []Token {
	for !lexer.isAtEnd() {

		lexer.start = lexer.current
		lexer.scanToken()
	}
	lexer.tokens = append(lexer.tokens, Token{"EOF", "", lexer.start, lexer.current, lexer.line})
	return lexer.tokens
}

func (lexer *Lexer) isAtEnd() bool {
	return lexer.current >= len(lexer.source)
}

func (lexer *Lexer) scanToken() {
	switch c := lexer.advance(); c {
	case " ":
    if lexer.insideTag && lexer.peekPrevious() == "<" {
      htg.Error(lexer.line, lexer.pos, "Unexpected character.", "Help: Remove space after '<' character.")
    }
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
    lexer.insideTag = true
  
    if lexer.peek() == "!" && lexer.peekNext() == "D" {
      lexer.addToken(LEFT_ANGLE.Type)
      lexer.advance()
      lexer.advanceToken()
      lexer.tag()
      break
    }

		lexer.addToken(LEFT_ANGLE.Type)
	case ">":
		lexer.addToken(RIGHT_ANGLE.Type)
    lexer.insideTag = false
	case "/":
    lexer.addToken(SLASH.Type)
	case "=":
		lexer.addToken(EQUAL.Type)
	case "\"":
    if !lexer.insideTag {
      lexer.literal()
      break;
    }
		lexer.string()
	default:
    if !lexer.insideTag {
      lexer.literal()
      break;
    }

    if ( lexer.peekPrevious() == "<" || lexer.peekPrevious() == "/" ) && lexer.isAlpha(c) {
      lexer.tag()
      break
    }

    if lexer.isAlpha(c) {
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
  // Replace all leading and trailing spaces as well as newlines
  // with regex
  text = strings.TrimSpace(text)
  text = strings.Replace(text, "\n", "", -1)
	lexer.addTokenWithLiteral(tokenType, text)
}

func (lexer *Lexer) addTokenWithLiteral(tokenType string, literal string) {
	lexer.tokens = append(lexer.tokens, Token{tokenType, literal, lexer.start, lexer.current, lexer.line})
}

func (lexer *Lexer) peek() string {
	if lexer.isAtEnd() {
		return "\000"
	}
	return string(lexer.source[lexer.current])
}

func (lexer *Lexer) peekPrevious() string {
  if lexer.current-2 < 0 {
    return "\000"
  }
  return string(lexer.source[lexer.current-2])
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

func (lexer *Lexer) isAlpha(c string) bool {
  return strings.Contains("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", c)
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
	for lexer.peek() != "\"" && !lexer.isAtEnd() {
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
	for lexer.isAlpha(lexer.peek()) || lexer.isHyphen(lexer.peek()) || lexer.isUnderscore(lexer.peek()) {
		lexer.advance()
	}
	lexer.addToken("ATTRIBUTE")
}

func (lexer *Lexer) literal() {
  for !lexer.isAtEnd() && lexer.peek() != "<" {
    lexer.advance()
  }
  lexer.addToken("LITERAL")
}

func (lexer *Lexer) tag() {
  for lexer.isAlphaNumeric(lexer.peek()) {
    lexer.advance()
  }
  lexer.addToken("TAG")
}
