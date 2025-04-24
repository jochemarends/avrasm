package encoding

import (
	"avrasm/arch"
	"avrasm/ast"
	"avrasm/token"
	"encoding/binary"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

type typeMapper func(o *ast.Operand) (arch.Word, error)

type constraintMapper func(w arch.Word) (arch.Word, error)

func isByte(w arch.Word) (arch.Word, error) {
	if w < 0 || w > 0xff {
        return w, fmt.Errorf("value %d cannot be represented as a byte", w)
    }
	return w, nil
}

func integer(o *ast.Operand) (arch.Word, error) {
	if o.Token.Kind != token.Integer {
		return 0, fmt.Errorf("@%d:%d: expected an integer operand", o.Token.Pos.Line, o.Token.Pos.Column)
	}
	data, err := strconv.Atoi(o.Token.Lexeme)
	return arch.Word(data), err
}

func register(o *ast.Operand) (arch.Word, error) {
	lexeme := o.Token.Lexeme

	if len(lexeme) < 2 || lexeme[0] != 'r' {
		return 0, fmt.Errorf("@%d:%d: expected an register operand", o.Token.Pos.Line, o.Token.Pos.Column)
	}

	num, err := strconv.Atoi(lexeme[1:])
	if err != nil || num < 0 || num > 31 {
		return 0, fmt.Errorf("@%d:%d: invalid register format `%s`", o.Token.Pos.Line, o.Token.Pos.Column, lexeme)
	}
	return arch.Word(num), nil
}

type operand struct {
	symbol      rune
	kind        typeMapper
	constraints []constraintMapper
}

func (o *operand) encode(node *ast.Operand) (word arch.Word, err error) {
	word, err = o.kind(node)

	for _, constraint := range o.constraints {
		word, err = constraint(word)
		if err != nil {
			break
		}
	}

	return
}

type instruction struct {
	encoding string
	operands []operand
}

var encodings = map[string]instruction{
    "nop": {
        encoding: "0000 0000 0000 0000",
    },
	"mov": {
		encoding: "0010 11rd dddd rrrr",
		operands: []operand{
			{symbol: 'd', kind: register},
			{symbol: 'r', kind: register},
		},
	},
}

func Encode(stmt *ast.Statement, w io.Writer) error {
	switch v := (*stmt).(type) {
	case *ast.Storage:
        err := EncodeStorage(v, w)
        if err != nil {
            return err
        }
	case *ast.Instruction:
        word, err := EncodeInstruction(v)
        if err != nil {
            return err
        }
        if err := binary.Write(w, binary.LittleEndian, word); err != nil {
            return err
        }
	}
	return nil
}

func EncodeStorage(storage *ast.Storage, w io.Writer) error {
	if storage.Directive.Lexeme != "byte" {
		return fmt.Errorf("@%d:%d: invalid storage directive `%s`", storage.Directive.Pos.Line, storage.Directive.Pos.Column, storage.Directive.Lexeme)
	}

	encoding := operand{
		kind: integer,
		constraints: []constraintMapper{isByte},
	}

	for _, op := range storage.Operands {
		data, err := encoding.encode(op)
		if err != nil {
			return err
		}
        if err := binary.Write(w, binary.LittleEndian, byte(data)); err != nil {
            return err
        }
	}
	return nil
}

func EncodeInstruction(instr *ast.Instruction) (encoded arch.Word, err error) {
	encoding, ok := encodings[strings.ToLower(instr.Mnemonic.Lexeme)]

	if len(instr.Operands) != len(encoding.operands) {
		return 0, fmt.Errorf("@%d:%d: expected %d operands, received %d", instr.Mnemonic.Pos.Line, instr.Mnemonic.Pos.Column, len(encoding.operands), len(instr.Operands))
	}

	if !ok {
		return 0, fmt.Errorf("@%d:%d: invalid mnemonic `%s`", instr.Mnemonic.Pos.Line, instr.Mnemonic.Pos.Column, instr.Mnemonic.Lexeme)
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

