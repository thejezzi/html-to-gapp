package htg

import (
	"fmt"
	"github.com/thejezzi/html-to-gapp/lib/logger"
	"os"
	"strings"
)

type Lexer struct {
	source       string
	tokens       []Token
	currentToken int
	line         int
	start        int
	current      int
}

type Token struct {
	tokenType string
	value     string
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
	// html comment
	case "<":
		if lexer.peek() == "!" && lexer.peekNext() == "-" && lexer.peekNextTwo() == "-" {
			lexer.Comment()
			break
		}
		lexer.Tag()
	default:
		logger.Error(fmt.Sprintf("Unexpected character: %s", c))
	}
}

func (lexer *Lexer) advance() string {
	lexer.current++
	return string(lexer.source[lexer.current-1])
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
	for lexer.peek() != "-" && lexer.peekNext() != "-" && lexer.peekNextTwo() != ">" && !lexer.isAtEnd() {
		lexer.advance()
	}
	for i := 0; i < 3; i++ {
		lexer.advance()
	}
}

func (lexer *Lexer) Tag() {
	for (lexer.peek() != ">") && !lexer.isAtEnd() {
		lexer.advance()
	}
	lexer.advance()

	// A Tag has the following format:
	// <tagname attribute="value" attribute="value">
	// We need to extract the tagname and the attributes

	// Get the text
	text := lexer.source[lexer.start+1 : lexer.current-1]

	// Get the tagname
	tagname := strings.Split(text, " ")[0]

	if string(tagname[0]) == "/" {
		lexer.addTokenWithLiteral("TAG_CLOSE", tagname[1:])
		return
	}

	// Get the attributes text
	attributes := text[len(tagname):]
	// Get the attributes as a map
	attributesMap := make(map[string]string)

	for _, attribute := range strings.Split(attributes, " ") {
		if strings.Contains(attribute, "=") {
			attributeSplit := strings.Split(attribute, "=")
			attributeName := attributeSplit[0]
			attributeValue := attributeSplit[1]
			if attributeValue[0] == '"' && attributeValue[len(attributeValue)-1] == '"' {
				attributeValue = attributeValue[1 : len(attributeValue)-1]
			} else {
				logger.Error(fmt.Sprintf("Attribute value is not surrounded by quotes: %s", attributeValue))
				os.Exit(1)
			}
			attributesMap[attributeName] = attributeValue
		}
	}

	// Create the tokens
	lexer.tokens = append(lexer.tokens, Token{"TAG", tagname})
	for attributeName, attributeValue := range attributesMap {
		lexer.tokens = append(lexer.tokens, Token{"ATTRIBUTE", attributeName})
		lexer.tokens = append(lexer.tokens, Token{"ATTRIBUTE_VALUE", attributeValue})
	}

	// Get everything after the tag opening
	for lexer.peekNext() != "<" && !lexer.isAtEnd() {
		lexer.advance()
	}
}
