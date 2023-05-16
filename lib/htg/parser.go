package htg

// import (
//   "fmt"
//   "github.com/thejezzi/html-to-gapp/lib/logger"
// )

type Parser struct {
  tokens []Token
  current int
}

type Tag struct {
  name string
  attributes []Attribute
}

type Attribute struct {
  name string
  value string
}

type Tags struct {
  tag Tag
  children []Tags
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
// func (p *Parser) Parse(tokens []Token) []Stmt {
//   p.tokens = tokens
//   p.current = 0
//
//   var tags []Tags
//   for !p.isAtEnd() {
//     tags = append(statements, p.declaration())
//   }
//
//   return tags
// }
//
// func (p *Parser) declaration() Tags {
//   if p.match(LEFT_ANGLE) {
//     return p.tag()
//   }
//
//   return p.text()
// }
