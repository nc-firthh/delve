// This script generates pkg/dwarf/op/opcodes.go from pkg/dwarf/op/opcodes.txt

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
)

type Opcode struct {
	Name string
	Code string
	Args string
	Func string
}

func main() {
	fh, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	outfh := os.Stdout
	if os.Args[2] != "-" {
		outfh, err = os.Create(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		defer outfh.Close()
	}

	opcodes := []Opcode{}
	s := bufio.NewScanner(fh)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		fields := strings.Split(line, "\t")
		opcode := Opcode{Name: fields[0], Code: fields[1], Args: fields[2]}
		if len(fields) > 3 {
			opcode.Func = fields[3]
		}
		opcodes = append(opcodes, opcode)
	}

	var buf bytes.Buffer

	fmt.Fprintf(&buf, `// THIS FILE IS AUTOGENERATED, EDIT opcodes.table INSTEAD
	
package op
`)

	// constants
	fmt.Fprintf(&buf, "const (\n")
	for _, opcode := range opcodes {
		fmt.Fprintf(&buf, "%s Opcode = %s\n", opcode.Name, opcode.Code)
	}
	fmt.Fprintf(&buf, ")\n\n")

	// name map
	fmt.Fprintf(&buf, "var opcodeName = map[Opcode]string{\n")
	for _, opcode := range opcodes {
		fmt.Fprintf(&buf, "%s: %q,\n", opcode.Name, opcode.Name)
	}
	fmt.Fprintf(&buf, "}\n")

	// arguments map
	fmt.Fprintf(&buf, "var opcodeArgs = map[Opcode]string{\n")
	for _, opcode := range opcodes {
		fmt.Fprintf(&buf, "%s: %s,\n", opcode.Name, opcode.Args)
	}
	fmt.Fprintf(&buf, "}\n")

	// function map
	fmt.Fprintf(&buf, "var oplut = map[Opcode]stackfn{\n")
	for _, opcode := range opcodes {
		if opcode.Func != "" {
			fmt.Fprintf(&buf, "%s: %s,\n", opcode.Name, opcode.Func)
		}
	}
	fmt.Fprintf(&buf, "}\n")

	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	outfh.Write(src)
}
