package htg

// import (
//   "fmt"
//   "github.com/thejezzi/html-to-gapp/lib/logger"
// )

type Parser struct {
	tokens  []Token
	current int

	tags []*Tag
}

/*
* Grammar
*
* declaration -> tag | text
* tag -> '<' NAME attributes? '>' declaration* '</' NAME '>'
* attributes -> attribute*
* attribute -> NAME '=' STRING
* text -> STRING
 */
func (p *Parser) Parse(tokens []Token) []*Tag {
	p.tokens = tokens
	p.current = 0

	p.tags = make([]*Tag, 0)

	for !p.isAtEnd() {
		p.tags = append(p.tags, p.declaration())
	}
	return p.tags
}

func (p *Parser) declaration() *Tag {
	if p.match(LEFT_ANGLE) {
		return p.tag()
	} else {
		return p.text()
	}
}

func (p *Parser) tag() *Tag {
	thisTag := &Tag{}
	currentToken := p.peek()
	if p.match(NAME) {
		thisTag.name = currentToken.value
	}

	for !p.check(RIGHT_ANGLE.Type) && !p.isAtEnd() {
		p.attribute()
	}

	if !p.match(RIGHT_ANGLE) {
		htg.Error(p.peek().line, p.peek().start, "Unexpected character.", "Help: Remove space after '<' character.")
	}

	if thisTag.name == "DOCTYPE" {
		return thisTag
	}
	// while not end of file and not a closing tag
	for !p.isAtEnd() && !(p.previous().tokenType == SLASH.Type) && !(p.peek().tokenType == NAME.Type) && !(p.peek().value == thisTag.name) {
		thisTag.children = append(thisTag.children, p.declaration())
		if p.peek().tokenType == LEFT_ANGLE.Type && p.peekNext().tokenType == SLASH.Type {
			p.advance()
			p.advance()
			break
		}
	}

	if p.match(NAME) {
		if p.previous().value != thisTag.name {
			htg.Error(p.peek().line, p.peek().start, "No end tag provided.", "Help: You must close the tag with the same name.")
		}
	}

	if !p.match(RIGHT_ANGLE) {
		htg.Error(p.peek().line, p.peek().start, "Unexpected character.", "Help: Remove space after '<' character.")
	}

	return thisTag
}

func (p *Parser) attribute() {
	var attr_name string
	var attr_value string
	if p.match(ATTRIBUTE) {
		attr_name = p.previous().value
	} else {
		htg.Error(p.peek().line, p.peek().start, "Unexpected character.", "Help: Remove space after '<' character.")
	}
	if p.match(EQUAL) {
		p.advance()
	} else {
		if p.peek().tokenType == RIGHT_ANGLE.Type {
			p.addAttr(attr_name, "")
			return
		}
		htg.Error(p.peek().line, p.peek().start, "Unexpected character.", "Help: If you try to define a value for an attribute, you must use '=' character.")
	}
	if p.match(STRING) {
		attr_value = p.previous().value
	}
	p.addAttr(attr_name, attr_value)
}

func (p *Parser) text() *Tag {
	if p.match(TEXT) {
		token := p.addText()
		// advance twice to skip the next token so that the loop can check on
		// previous as slash and peek as name
		return token
	}
	htg.Error(p.peek().line, p.peek().start, "Unexpected string value.", "Help: Something went terribly wrong.")
	p.synchronize()
	return nil
}

func (p *Parser) addText() *Tag {
	return &Tag{name: "TEXT", value: p.previous().value}
}

func (p *Parser) addAttr(name string, value string) *Attribute {
	return &Attribute{name: name, value: value}
}

func (p *Parser) match(tokenType TokenType) bool {
	if p.check(tokenType.Type) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) check(tokenType string) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF.Type
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) peekNext() Token {
	return p.tokens[p.current+1]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType string, message string) Token {
	if p.check(tokenType) {
		return p.advance()
	}
	panic(message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		switch p.peek().tokenType {
		case LEFT_ANGLE.Type:
			return
		case STRING.Type:
			return
		}
		p.advance()
	}
}
