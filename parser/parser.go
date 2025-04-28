package parser

import (
	"avrasm/ast"
	"avrasm/lexer"
	"avrasm/token"
	"fmt"
	"io"
	"iter"
	"slices"
)

func ParseFile(r io.Reader) (file *ast.File, err error) {
	l := lexer.New(r)

	var toks []token.Token
	for {
		tok, err := l.Scan()
		if err != nil {
			return nil, err
		}
		toks = append(toks, *tok)
		if tok.Kind == token.EOF {
			break
		}
	}

	p := New(slices.Values(toks))
	file = &ast.File{}
	for !p.Done() {
		instr, err := p.ParseInstruction()
		if err != nil {
			return nil, err
		}
		file.Statements = append(file.Statements, instr)
	}
	return
}

type Parser struct {
	next      func() (token.Token, bool)
	stop      func()
	currToken *token.Token
	nextToken *token.Token
	isDone	  bool
}

func New(tokens iter.Seq[token.Token]) *Parser {
	var p Parser
	p.next, p.stop = iter.Pull(tokens)	
	p.advance()
	return &p
}

func (p *Parser) Close() {
	p.stop()
}

func (p *Parser) advance() {
	p.currToken = p.nextToken
	tok, ok := p.next()

	if !ok || tok.Kind == token.EOF {
		p.isDone = true
	} else {
		p.nextToken = &tok
	}
}

func (p *Parser) Done() bool {
	return p.isDone
}

func (p *Parser) test(kinds ...token.Kind) bool {
	if p.nextToken != nil {
		return slices.Contains(kinds, p.nextToken.Kind)
	}
	return false
}

func (p *Parser) match(kinds ...token.Kind) bool {
	if p.test(kinds...) {
        p.advance()
        return true
    }
    return false
}

func (p *Parser) expect(kinds ...token.Kind) (tok *token.Token, err error) {
	if p.match(kinds...) {
		tok = p.currToken
    } else {
		if p.nextToken == nil {
			err = fmt.Errorf("@%d:%d: expected token of type `%v`", p.nextToken.Pos.Line, p.nextToken.Pos.Column, kinds)
		} else {
			err = fmt.Errorf("@%d:%d: expected token of in `%v`, received `%v`", p.nextToken.Pos.Line, p.nextToken.Pos.Column, kinds, p.nextToken.Kind)
		}
	}
    return
}

func (p *Parser) ParseInstruction() (stmt ast.Statement, err error) {
	if p.test(token.Dot) {
		stmt, err = p.parseStorage()
	} else {
		stmt, err = p.parseInstruction()
	}

	if err == nil {
		_, err = p.expect(token.Newline)
	}

	return
}

func (p *Parser) parseStorage() (*ast.Storage, error) {
	_, err := p.expect(token.Dot)
	if err != nil {
		return nil, err
	}

	tok, err := p.expect(token.Ident)
	if err != nil {
		return nil, err
	}

	ops, err := p.parseOperands()
	if err != nil {
		return nil, err
	}

	return &ast.Storage{Directive: tok, Operands: ops}, nil
}

func (p *Parser) parseInstruction() (*ast.Instruction, error) {
	tok, err := p.expect(token.Ident)
	if err != nil {
		return nil, err
	} else {
		ops, err := p.parseOperands()
		if err != nil {
			return nil, err
		}
		return &ast.Instruction{Mnemonic: tok, Operands: ops}, nil
	}
}

func (p *Parser) parseOperands() ([]*ast.Operand, error) {
	var ops []*ast.Operand

	if !p.match(token.Ident, token.Integer) {
		return ops, nil
	} else {
		ops = append(ops, &ast.Operand{Token: p.currToken})
	}
	
	for p.match(token.Comma) {
		tok, err := p.expect(token.Ident, token.Integer)
		if err != nil {
			return nil, err
		}
		ops = append(ops, &ast.Operand{Token: tok})
	}

	return ops, nil
}

func (p *Parser) parseOperand() (*ast.Operand, error) {
	tok, err := p.expect(token.Ident, token.Integer)
	if err != nil {
		return nil, err
	}
	return &ast.Operand{Token: tok}, nil
}

