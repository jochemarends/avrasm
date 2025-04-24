package main

import (
	"avrasm/encoding"
	"avrasm/parser"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
)

func assemble(r io.Reader, w io.Writer) error {
	file, err := parser.ParseFile(r)
	if err != nil {
		return err
	}

	for _, instr := range file.Statements {
		err := encoding.Encode(&instr, w)
		if err != nil {
			return err
		}
		// if err := binary.Write(w, binary.LittleEndian, word); err != nil {
		// 	return err
		// }
	}

	return nil
}

func dump(r io.Reader) {
	for true {
		var value uint16
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			break
		}
		fmt.Printf("0x%04X\n", value)
	}
}

func main() {
	output := flag.String("o", "out.bin", "output file name")

	flag.Parse()

	for _, s := range flag.Args() {
		ifile, err := os.Open(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
		defer ifile.Close()

		ofile, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
		defer ofile.Close()

		if err := assemble(ifile, ofile); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
	}
}

