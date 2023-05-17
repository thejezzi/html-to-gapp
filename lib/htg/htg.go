package htg

import (
	"fmt"
	"github.com/thejezzi/html-to-gapp/lib/logger"
	"os"
	"strconv"
	"strings"
)

type HTG struct {
	lexer  *Lexer
	parser *Parser
}

func (htg *HTG) printErrPosition(hint string) {

	// Int to string
	lineNumber := strconv.Itoa(htg.lexer.line)
	prevLineNumber := strconv.Itoa(htg.lexer.line - 1)
	nextLineNumber := strconv.Itoa(htg.lexer.line + 1)

	allLinesAsList := strings.Split(htg.lexer.source, "\n")
	line := lineNumber + ": " + allLinesAsList[htg.lexer.line]
	prevLine := prevLineNumber + ": " + allLinesAsList[htg.lexer.line-1]
	nextLine := nextLineNumber + ": " + allLinesAsList[htg.lexer.line+1]
	arrowTip := strings.Repeat(" ", htg.lexer.pos+2) + "^"
	errMsgAtPipe := strings.Repeat(" ", htg.lexer.pos+2) + ""
	logger.Error(prevLine)
	logger.Error(line)
	logger.Error(logger.Colorize(logger.FG_RED, arrowTip))
	logger.Error(logger.Colorize(logger.FG_CYAN, errMsgAtPipe+hint))
	logger.Error(nextLine)
}

func (htg *HTG) Error(line, pos int, message, hint string) {
	errorMsg := fmt.Sprintf("[line %d:%d] Error: %s", line, pos, message)
	logger.Error(errorMsg)
	htg.printErrPosition(hint)
	os.Exit(1)
}

func Convert(htmlFile string, gappFile string) {

}

var htg HTG

func TestRun() {
	// new lexer instance
	htg = HTG{
		lexer:  &Lexer{},
		parser: &Parser{},
	}

	htg.lexer.source = `
<!DOCTYPE html>
<html>
<head>
  <title>Test</title>
</head>
<body>
  <h1>Hello World</h1>
  <p class="Test">This is a test</p>
  <div class="Test">
	<p>This is a test</p>
	<p>This is a test</p>
	<div class="sub">
		<p>This is a test</p>
		<p>This is a test</p>
	</div>
  </div>
  <script>
    var x = 1;
    var y = 2;
    var z = x + y;
    console.log(z);
  </script>
</body>
</html>
`

	logger.Info(htg.lexer.source)

	// scan tokens
	tokens := htg.lexer.ScanTokens()
	parser := htg.parser.Parse(tokens)

	for _, tag := range parser {
		printTag(tag, 0)
	}
}

func printTag(tag *Tag, indent int) {
	if tag.name == "TEXT" {
		fmt.Println(strings.Repeat(" ", indent), tag.name, tag.value)
	} else {
		fmt.Println(strings.Repeat(" ", indent), tag.name)
	}
	for _, attr := range tag.attributes {
		fmt.Println(strings.Repeat(" ", indent*2), attr.name+" = "+attr.value)
	}
	for _, child := range tag.children {
		printTag(child, indent+2)
	}
}
