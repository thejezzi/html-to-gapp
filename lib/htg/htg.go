package htg

func Convert(htmlFile string, gappFile string) {
	// new lexer instance
	lexer := Lexer{source: `
  <html>
    <head>
      <title>Test</title>
    </head>
    <body>

      <!-- Test -->
      <h1 class="Testclass">Test</h1>
      <p>Test</p>
    </body>
  </html>
  `}

	// scan tokens
	tokens := lexer.ScanTokens()

	ast := AST{}
	ast.Parse(tokens)

}
