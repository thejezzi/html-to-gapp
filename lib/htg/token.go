package htg

import "fmt"

type Token struct {
	tokenType string
	value     string
  start     int
  end       int
  line      int
}

func (t *Token) String() string {
  return fmt.Sprintf("%s %s [%d:%d-%d]", t.tokenType, t.value, t.line, t.start, t.end)
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
