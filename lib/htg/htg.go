package htg

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/thejezzi/html-to-gapp/lib/logger"
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
	var prevLine string
	if htg.lexer.line-1 <= len(allLinesAsList) {
		prevLine = ""
	} else {
		prevLine = prevLineNumber + ": " + allLinesAsList[htg.lexer.line-1]
	}

	var nextline string
	if htg.lexer.line+1 >= len(allLinesAsList) {
		nextline = ""
	} else {
		nextline = nextLineNumber + ": " + allLinesAsList[htg.lexer.line+1]
	}
	arrowTip := strings.Repeat(" ", htg.lexer.pos+2) + "^"
	errMsgAtPipe := strings.Repeat(" ", htg.lexer.pos+2) + ""
	logger.Error(prevLine)
	logger.Error(line)
	logger.Error(logger.Colorize(logger.FG_RED, arrowTip))
	logger.Error(logger.Colorize(logger.FG_CYAN, errMsgAtPipe+hint))
	logger.Error(nextline)
}

func (htg *HTG) PrettyError(line, pos int, message, hint string) {
	errorMsg := fmt.Sprintf("[line %d:%d] Error: %s", line, pos, message)
	logger.Error(errorMsg)
	htg.printErrPosition(hint)
	os.Exit(1)
}

func (htg *HTG) Error(line, pos int, message string) {
	errorMsg := fmt.Sprintf("[line %d:%d] Error: %s", line, pos, message)
	logger.Error(errorMsg)
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
		<p>This is a test</p>
	</body>
</html>

	`

	// scan tokens
	tokens := htg.lexer.ScanTokens()
	parser := htg.parser.Parse(tokens)

	for _, tag := range parser {
		printTag(tag, 0)
	}
}

func printTag(tag *Tag, indent int) {
	tagname := fmt.Sprintf(strings.Repeat("  ", indent)+"%v", tag.name)

	fmt.Println(tagname) //, tag.value)
	//for _, attr := range tag.attributes {
	//	fmt.Println(strings.Repeat("  ", indent) + "~" + attr.name + "[" + attr.value + "]")
	//}

	for _, child := range tag.children {
		printTag(child, indent+2)
	}
}
