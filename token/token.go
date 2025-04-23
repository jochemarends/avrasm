package token

import "text/scanner"

type Token struct {
    Kind   Kind
    Lexeme string
	Pos    scanner.Position
}

type Kind string

const (
    Comma   Kind = ","
	Dot	    Kind = "."
    Integer Kind = "Integer"
	Ident   Kind = "Ident"
	Newline Kind = "\\n"
    EOF     Kind = "EOF"
)

