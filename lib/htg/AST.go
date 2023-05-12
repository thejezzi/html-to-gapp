package htg

import (
	"github.com/thejezzi/html-to-gapp/lib/logger"
)

type Node struct {
	nodeType string
	value    string
	children []Node
}

type AST struct {
	tree    []Node
	start   int
	current int
}

func (ast *AST) Parse(tokens []Token) {
	for _, token := range tokens {
		switch token.tokenType {
		case "TAG":
			// ast.Tag()
		case "EOF":
			logger.Info("EOF")
		default:
			logger.Error("Unexpected token type")
		}
	}
}
