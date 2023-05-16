package htg

type TokenType struct {
  Type string
  Lexeme string
}

var (
  LEFT_ANGLE = TokenType{Type: "LEFT_ANGLE", Lexeme: "<"}
  RIGHT_ANGLE = TokenType{Type: "RIGHT_ANGLE", Lexeme: ">"}
  SLASH = TokenType{Type: "SLASH", Lexeme: "/"}
  EQUAL = TokenType{Type: "EQUAL", Lexeme: "="}
  STRING = TokenType{Type: "STRING", Lexeme: "STRING"}
  NAME = TokenType{Type: "NAME", Lexeme: "NAME"}
  EOF = TokenType{Type: "EOF", Lexeme: ""}
)
