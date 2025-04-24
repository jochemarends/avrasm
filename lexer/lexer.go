package lexer

import (
	"avrasm/token"
	"fmt"
	"io"
	"text/scanner"
)

type Error struct {
    Msg string
    Pos scanner.Position
}

func (e Error) Error() string {
	return fmt.Sprintf("@%d:%d: %s", e.Pos.Line, e.Pos.Column, e.Msg)
}

func Tokenize(r io.Reader) (tokens []token.Token, errors []error) {
	l := New(r)
	for {
		tok, err := l.Scan()
		if err != nil {
			errors = append(errors, *err)
		} else {
			tokens = append(tokens, *tok)
			if tok.Kind == token.EOF {
				break
			}
		}
	}
	return
}

type Lexer struct {
    scanner scanner.Scanner
}

// Returns a new [Lexer] to read from r.
func New(r io.Reader) *Lexer {
    var l Lexer
    l.scanner.Init(r)
	l.scanner.Mode = scanner.ScanChars
	l.scanner.Whitespace ^= 1<<'\n'
    return &l
}

// Scan reads and returns the next token.
func (l *Lexer) Scan() (*token.Token, *Error) {
    r := l.scanner.Scan()

    switch r {
    case scanner.EOF:
        return l.newToken(token.EOF, ""), nil
    case ',':
        return l.newToken(token.Comma, string(r)), nil
    case '.':
        return l.newToken(token.Dot, string(r)), nil
    case '\n':
        return l.newToken(token.Newline, string(r)), nil
    case ';':
		for l.scanner.Peek() != '\n' && l.scanner.Peek() != scanner.EOF {
			l.scanner.Scan()
		}
		return l.Scan()
    }

    if isDigit(r) {
        return l.scanInt(string(r)), nil
    }

	if isLetter(r) || r == '_' {
		return l.scanIdent(string(r)), nil
	}

    return nil, l.newError(fmt.Sprintf("illegal character `%c`", r))
}

func (l *Lexer) scanInt(lexeme string) *token.Token {
    for {
        r := l.scanner.Peek()

        if !isDigit(r) {
            return l.newToken(token.Integer, lexeme)
        }

        lexeme += string(r)
        l.scanner.Scan()
    }
}

func (l *Lexer) scanIdent(lexeme string) *token.Token {
    for {
        r := l.scanner.Peek()

        if !isDigit(r) && !isLetter(r) {
            return l.newToken(token.Ident, lexeme)
        }

        lexeme += string(r)
        l.scanner.Scan()
    }
}

func (l *Lexer) newToken(kind token.Kind, lexeme string) *token.Token {
    return &token.Token{
        Kind:   kind,
        Lexeme: lexeme,
		Pos: 	l.scanner.Position,
    }
}

func (l *Lexer) newError(msg string) *Error {
    return &Error{
        Msg: msg,
        Pos: l.scanner.Position,
    }
}

func isDigit(r rune) bool {
    return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
    return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

