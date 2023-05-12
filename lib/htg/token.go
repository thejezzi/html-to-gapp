package htg

import "fmt"

type Token struct {
	tokenType string
	value     string
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %s", t.tokenType, t.value)
}

var Tokentype []string = []string{
	"TAG_START",
	"TAG_END",
	"TAG_NAME",
	"TAG_END_SELF_CLOSE",
	"TAG_END_CLOSE",
	"TAG_ATTRIBUTE",
	"TAG_ATTRIBUTE_VALUE",
	"LITERAL",
}
