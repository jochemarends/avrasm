package encoding

import (
	"avrasm/arch"
	"avrasm/ast"
	"avrasm/token"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type typeMapper func(o *ast.Operand) (arch.Word, error)

func integer(o *ast.Operand) (arch.Word, error) {
	if o.Token.Kind != token.Integer {
		return 0, fmt.Errorf("expected an integer operand")
	}
	data, err := strconv.Atoi(o.Token.Lexeme)
	return arch.Word(data), err
}

type operand struct {
	symbol rune
	kind   typeMapper
}

type instruction struct {
	encoding string
	operands []operand
}

var encodings = map[string]instruction{
    "nop": {
        encoding: "0000 0000 0000 0000",
    },
	// "mov": {
	// 	encoding: "0010 11rd dddd rrrr",
	// 	operands: []operand{
	// 		{symbol: 'd', kind: integer},
	// 		{symbol: 'r', kind: integer},
	// 	},
	// },
}

func Encode(instr *ast.Instruction) (encoded arch.Word, err error) {
	encoding, ok := encodings[instr.Mnemonic.Lexeme]

	if len(instr.Operands) != len(encoding.operands) {
		return 0, fmt.Errorf("expected %d operands, received %d", len(encoding.operands), len(instr.Operands))
	}

	if !ok {
		return 0, fmt.Errorf("unknown mnemonic `%s`", instr.Mnemonic.Lexeme)
	}

	bits := []byte(strings.ReplaceAll(encoding.encoding, " ", ""))

	for i, encoding := range encoding.operands {
		op := instr.Operands[i]
		data, err := encoding.kind(op)
		if err != nil {
			return 0, err
		}

		for i, c := range slices.Backward(bits) {
			if c == byte(encoding.symbol) {
				bits[i] = byte(data & 1) + byte('0')
				data >>= 1
			}
		}
	}

	word, err := strconv.ParseUint(string(bits), 2, 16)

	if err != nil {
		return 0, err
	}

	return arch.Word(word), nil
}

