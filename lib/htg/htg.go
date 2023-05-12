package htg

import "github.com/thejezzi/html-to-gapp/lib/logger"

func Convert(htmlFile string, gappFile string) {

}

func TestRun() {
	// new lexer instance
	lexer := Lexer{source: `
  	<html>
    	<head>
      		<title>Test</title>
    	</head>
    	<body>

      		<!-- Test 

			-->
      		<h1 class="Testclass">Test</h1>
      		<p>Test</p>
            <p data-bind="a001"></p>
			<input type="text" name="test" value="test" />
    	</body>
  	</html>
  `}

	logger.Info(lexer.source)

	// scan tokens
	tokens := lexer.ScanTokens()

	for _, token := range tokens {
		logger.Info(token.String())
	}

}
