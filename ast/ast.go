package ast

import "avrasm/token"

type File struct {
	Filename   string
	Statements []Statement
}

type Statement interface {}

type Instruction struct {
	Mnemonic *token.Token
	Operands []*Operand
}

type Storage struct {
	Directive *token.Token
	Operands  []*Operand
}

type Operand struct {
	Token *token.Token
}

