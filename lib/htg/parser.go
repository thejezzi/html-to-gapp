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

	for !p.check(RIGHT_ANGLE.Type) && !p.check(SLASH.Type) && !p.isAtEnd() {
		thisTag.attributes = append(thisTag.attributes, p.attribute())
	}

	if p.check(SLASH.Type) {

		// Skip the slash
		p.advance()

		if !p.check(RIGHT_ANGLE.Type) {
			htg.PrettyError(p.peek().line, p.peek().start, "Expected self closing tag.", "Help: Did you forgot a right angle bracket?")
		}
		// Skip right angle bracket to prevent error
		p.advance()
		return thisTag
	}

	if !p.match(RIGHT_ANGLE) {
		htg.PrettyError(p.peek().line, p.peek().start, "Unexpected character.", "Help: Remove space after '<' character.")
	}

	switch thisTag.name {
	case "DOCTYPE":
		return thisTag
	case "meta":
		return thisTag
	case "link":
		return thisTag
	case "img":
		return thisTag
	case "input":
		return thisTag
	case "br":
		return thisTag
	case "hr":
		return thisTag
	case "area":
		return thisTag
	case "base":
		return thisTag
	case "col":
		return thisTag
	case "command":
		return thisTag
	case "embed":
		return thisTag
	case "keygen":
		return thisTag
	case "param":
		return thisTag
	case "source":
		return thisTag
	case "track":
		return thisTag
	case "wbr":
		return thisTag
	}

	// while not end of file and not a closing tag
	for !p.isAtEnd() && !(p.previous().tokenType == SLASH.Type) && !(p.peek().tokenType == NAME.Type) && !(p.peek().value == thisTag.name) {

		// If the current and next one are `</` skip them to eventually fullfil the condition to end the loop
		if p.peek().tokenType == LEFT_ANGLE.Type && p.peekNext().tokenType == SLASH.Type {
			p.advance()
			p.advance()
			break
		}
		thisTag.children = append(thisTag.children, p.declaration())
	}

	if p.match(NAME) {
		if p.previous().value != thisTag.name {
			htg.PrettyError(p.peek().line, p.peek().start, "No end tag provided.", "Help: You must close the tag with the same name.")
		}
	}

	if !p.match(RIGHT_ANGLE) {
		if len(p.tokens) == p.current+1 {
			lastToken := p.tokens[len(p.tokens)-1]
			htg.PrettyError(lastToken.line, lastToken.start, "Unexpected end of file.", "Help: Consider closing all tags")
		}
		htg.PrettyError(p.peek().line, p.peek().start, "Unexpected character.", "Help: Remove space after '<' character.")
	}

	return thisTag
}

func (p *Parser) attribute() *Attribute {
	var attr_name string
	var attr_value string
	if p.match(ATTRIBUTE) {
		attr_name = p.previous().value
	} else {
		htg.PrettyError(p.peek().line, p.peek().start, "Expected attribute", "Help: ... I ... I dont know")
	}
	if !p.match(EQUAL) {
		if p.check(RIGHT_ANGLE.Type) || p.check(SLASH.Type) {
			return p.addAttr(attr_name, "")
		}
		if p.check(ATTRIBUTE.Type) {
			return p.addAttr(attr_name, "")
		}
		htg.PrettyError(p.peek().line, p.peek().start, "Unexpected character.", "Help: If you try to define a value for an attribute, you must use '=' character.")
	}
	if p.match(STRING) {
		attr_value = p.previous().value
	}
	return p.addAttr(attr_name, attr_value)
}

func (p *Parser) text() *Tag {
	if p.match(TEXT) {
		token := p.addText()
		// advance twice to skip the next token so that the loop can check on
		// previous as slash and peek as name
		return token
	}
	htg.PrettyError(p.peek().line, p.peek().start, "Unexpected string value.", "Help: Something went terribly wrong.")
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
