package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	tool := flag.String("tool", "asm", "Tool select: asm/dasm")
	input := flag.String("input", "", "Input filename")
	output := flag.String("output", "", "Output filename")

	flag.Parse()

	if *input == "" || *output == "" {
		fmt.Println("Please enter input/output filename")
		fmt.Println("type -h for help")
		os.Exit(1)
	}

	switch *input {
	case "asm":
		assemble(*input, *output)
	case "dasm":
		disassemble(*input, *output)
	default:
		fmt.Printf("Unsupport tool: %s\n", *tool)
		fmt.Println("Support tool are [asm, dasm]")
		fmt.Println("type -h for help")
		os.Exit(1)
	}
}

func openInOut(inp string, opt string) (*os.File, *os.File, error) {
	f, err := os.Open(inp)
	if err != nil {
		return nil, nil, err
	}

	fopt, err := os.Open(opt)
	if err != nil {
		return nil, nil, err
	}

	return f, fopt, nil
}

func checkRuneInRunes(rs []rune, r rune) bool {
	for _, ar := range rs {
		if ar == r {
			return true
		}
	}
	return false
}

func removeEmptyString(strs []string) []string {
	ret := make([]string, 0)
	for _, val := range strs {
		if len(val) > 0 {
			ret = append(ret, val)
		}
	}
	return ret
}

func getNumber(str string) (uint16, error) {
	var num uint64
	var err error
	var base int

	if str[0:1] == "x" {
		base = 16
	} else if str[0:1] == "#" {
		base = 10
	} else {
		err = fmt.Errorf("Unknown format of number: %s", str)
		return 0, err
	}

	num, err = strconv.ParseUint(str[1:], base, 16)
	if err != nil {
		return 0, err
	}
	return uint16(num), nil
}

// Convert 2 bytes to uint16
//
// LC3 is big-endian.
func getUint16(x []byte) uint16 {
	if len(x) != 2 {
		panic("getUint16 get length not equal to 2")
	}
	return uint16(x[0])<<8 | uint16(x[1])
}

func sepOpt(optSet []uint8, offset int, val uint16) {
	hi := uint8(val & 0xF)
	lo := uint8(val >> 8)
	optSet[offset] = hi
	optSet[offset] = lo
}

func assemble(inp string, opt string) {
	var err error
	var addr uint16 = 0
	var firstAddr uint16 = 0
	var offset int = 0
	optSet := make([]uint8, 2048)
	symbolList := make(map[string]uint16)
	codeList := make(map[int]string)

	f, fopt, err := openInOut(inp, opt)
	if err != nil {
		fmt.Printf("Assemble/openInOut: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()
	defer fopt.Close()

	// Read file
	ending := false
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if ending {
			break
		}

		if offset >= 2048 {
			fmt.Printf("Assemble failed: the obj size is larger than 2048, current neede size: %d\n", offset)
		}

		// Get line
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// Remove comment
		splitComment := strings.Split(line, ";")
		splitComment = removeEmptyString(splitComment)
		code := strings.TrimSpace(splitComment[0])

		// Check label
		runeCodes := []rune(code)
		if checkRuneInRunes(runeCodes, ':') {
			splitLabel := strings.Split(code, ":")
			splitLabel = removeEmptyString(splitLabel)
			symbolList[strings.TrimSpace(splitLabel[0])] = addr
			if len(splitLabel[1]) != 0 {
				code = strings.TrimSpace(splitLabel[1])
			} else {
				continue
			}
		}

		// Check code
		runeCodes = []rune(code)
		if runeCodes[0] == '.' { // Directive
			dir := line[1:]
			splitParam := strings.Split(dir, " ")
			splitParam = removeEmptyString(splitParam)
			switch splitParam[0] {
			case "ORIG":
				firstAddr, err = getNumber(splitParam[1])
				if err != nil {
					fmt.Printf("Assemble/getNumber(ORIG): %s\n", err.Error())
					os.Exit(1)
				}
				addr = firstAddr - 1 // Before first memory address

				sepOpt(optSet, 0, firstAddr)
				offset += 2
				break
			case "FILL":
				elem, err := getNumber(splitParam[1])
				if err != nil {
					fmt.Printf("Assemble/getNumber(FILL): %s\n", err.Error())
					os.Exit(1)
				}

				sepOpt(optSet, offset, elem)
				offset += 2
				break
			case "BLKW":
				count, err := getNumber(splitParam[1])
				if err != nil {
					fmt.Printf("Assemble/getNumber(BLKW): %s\n", err.Error())
					os.Exit(1)
				}

				for i := 0; i < int(count); i++ {
					sepOpt(optSet, offset, 0)
					offset += 2
				}
				break
			case "STRINGZ":
				str := splitParam[1][1 : len(splitParam[1])-1]
				bytes := []byte(str)
				for _, b := range bytes {
					optSet[offset] = b
					offset++
				}
				optSet[offset] = 0
				offset++
				break
			case "END":
				ending = true
			default:
				fmt.Printf("Unknown directive: %s\n", splitParam[0])
				os.Exit(1)
				break
			}
		} else { // Opcode
			codeList[offset] = code
			offset += 2
		}
		addr++
	}
	if err = scanner.Err(); err != nil {
		fmt.Printf("Assemble/fileScanner.Scan: %s\n", err.Error())
		os.Exit(1)
	}

	for loc := range codeList {
		code := codeList[loc]
		splitCode := strings.Split(code, " ")
		splitCode = removeEmptyString(splitCode)
		if opcode, ok := opToOpcode[splitCode[0]]; ok {
			joinParam := strings.Join(splitCode[1:], "")
			res, err := opcodeToByte(opcode, joinParam)
			if err != nil {
				fmt.Printf("Assemble/opcodeToByte: %s\n", code)
				os.Exit(1)
			}
			sepOpt(optSet, loc, res)

		} else {
			fmt.Printf("Assemble/opToOpcode[code]: %s\n", code)
			os.Exit(1)
		}
	}
}

func disassemble(inp string, opt string) {
	symCount := 0
	symbolList := make(map[uint16]string)
	memMapping := make(map[uint16]string)

	f, fopt, err := openInOut(inp, opt)
	if err != nil {
		fmt.Printf("Assemble/openInOut: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()
	defer fopt.Close()

	fInfo, err := f.Stat()
	if err != nil {
		fmt.Printf("Disassemble/f.Stat: %s\n", err.Error())
	}
	fSize := fInfo.Size()
	if fSize%2 != 0 {
		fSize++
	}
	dataByte := make([]byte, fSize)

	// Read file
	readN, err := f.Read(dataByte)
	if err != nil {
		fmt.Printf("Disassemble/f.Read: %s\n", err.Error())
	}
	fmt.Printf("Read %d byte\n", readN)

	// The origin tells us where in memory to place the image
	origin := getUint16(dataByte[0:2])
	addr := origin
	curByte := int64(2)

	fopt.WriteString(fmt.Sprintf(".ORIG x%x\n", origin))

	// Read byte into memory
	for curByte < fSize {
		ins := getUint16(dataByte[curByte : curByte+2])
		asm, jmpAddr, hasLabel := opcodeToOp(ins)
		memMapping[addr] = asm
		if hasLabel {
			symbolList[jmpAddr] = fmt.Sprintf("Label%d", symCount)
		}
		curByte += 2
		addr++
	}
}
