package ast

import "avrasm/token"

type File struct {
	Filename     string
	Instructions []Instruction
}

type Instruction struct {
	Mnemonic *token.Token
	Operands []*Operand
}

type Operand struct {
	Token *token.Token
}

