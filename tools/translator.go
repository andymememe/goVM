package main

import "fmt"

// OpsCode ID
const (
	OPBR   uint16 = iota /* branch */
	OPADD                /* add  */
	OPLD                 /* load */
	OPST                 /* store */
	OPJSR                /* jump register */
	OPAND                /* bitwise and */
	OPLDR                /* load register */
	OPSTR                /* store register */
	OPRTI                /* RTI (unused) */
	OPNOT                /* bitwise not */
	OPLDI                /* load indirect */
	OPSTI                /* store indirect */
	OPJMP                /* jump */
	OPRES                /* reserved (unused) */
	OPLEA                /* load effective address */
	OPTRAP               /* execute trap */
)

// Trap ID
const (
	TRAPGETC  uint16 = 0x20 /* Get character from keyboard, not echoed onto the terminal */
	TRAPOUT   uint16 = 0x21 /* Output a character */
	TRAPPUTS  uint16 = 0x22 /* Output a word string */
	TRAPIN    uint16 = 0x23 /* Get character from keyboard, echoed onto the terminal */
	TRAPPUTSP uint16 = 0x24 /* Output a byte string */
	TRAPHALT  uint16 = 0x25 /* Halt the program */
)

var opToOpcode = map[string]uint16{
	"OPBR":   OPBR,   /* branch */
	"OPADD":  OPADD,  /* add  */
	"OPLD":   OPLD,   /* load */
	"OPST":   OPST,   /* store */
	"OPJSR":  OPJSR,  /* jump register */
	"OPAND":  OPAND,  /* bitwise and */
	"OPLDR":  OPLDR,  /* load register */
	"OPSTR":  OPSTR,  /* store register */
	"OPRTI":  OPRTI,  /* RTI (unused) */
	"OPNOT":  OPNOT,  /* bitwise not */
	"OPLDI":  OPLDI,  /* load indirect */
	"OPSTI":  OPSTI,  /* store indirect */
	"OPJMP":  OPJMP,  /* jump */
	"OPRES":  OPRES,  /* reserved (unused) */
	"OPLEA":  OPLEA,  /* load effective address */
	"OPTRAP": OPTRAP, /* execute trap */
}

var trapToTrapcode = map[string]uint16{
	"TRAPGETC":  TRAPGETC,  /* Get character from keyboard, not echoed onto the terminal */
	"TRAPOUT":   TRAPOUT,   /* Output a character */
	"TRAPPUTS":  TRAPPUTS,  /* Output a word string */
	"TRAPIN":    TRAPIN,    /* Get character from keyboard, echoed onto the terminal */
	"TRAPPUTSP": TRAPPUTSP, /* Output a byte string */
	"TRAPHALT":  TRAPHALT,  /* Halt the program */
}

func opcodeToByte(opcode uint16, params string) (uint16, error) {
	var err error

	switch opcode {
	case OPADD: // Add
		return 0, err
	case OPAND: // And
		return 0, err
	case OPNOT: // Not
		return 0, err
	case OPBR: // Branch
		return 0, err
	case OPJMP: // Jump & Return
		return 0, err
	case OPJSR: // Jump register
		return 0, err
	case OPLD: // Load
		return 0, err
	case OPLDI: // Load Indirect
		return 0, err
	case OPLDR: // Load register
		return 0, err
	case OPLEA: // Load Effective Address
		return 0, err
	case OPST: // Store
		return 0, err
	case OPSTI:
		return 0, err
	case OPSTR:
		return 0, err
	case OPTRAP:
	case OPRES:
	case OPRTI:
	default:
		fmt.Printf("[ERROR] Unsupport opcode: x%X", opcode)
		return 0, err
	}
	fmt.Printf("[ERROR] Unknown opcode: x%X", opcode)
	return 0, err
}

func opcodeToOp(instruction uint16) (string, uint16, bool) {
	op := instruction >> 12

	switch op {
	case OPADD: // Add
		return "", 0, false
	case OPAND: // And
		return "", 0, false
	case OPNOT: // Not
		return "", 0, false
	case OPBR: // Branch
		return "", 0, true
	case OPJMP: // Jump & Return
		return "", 0, false
	case OPJSR: // Jump register
		return "", 0, false
	case OPLD: // Load
		return "", 0, true
	case OPLDI: // Load Indirect
		return "", 0, true
	case OPLDR: // Load register
		return "", 0, false
	case OPLEA: // Load Effective Address
		return "", 0, true
	case OPST: // Store
		return "", 0, true
	case OPSTI:
		return "", 0, true
	case OPSTR:
		return "", 0, false
	case OPTRAP:
		return trapcodeToTrap(instruction), 0, false
	case OPRES:
	case OPRTI:
	default:
		fmt.Printf("[ERROR] Unsupport opcode: x%X", instruction)
		return fmt.Sprintf("; [ERROR] Unsupport opcode: x%X", instruction), 0, false
	}
	fmt.Printf("[ERROR] Unknown opcode: x%X", instruction)
	return fmt.Sprintf("; [ERROR] Unknown opcode: x%X", instruction), 0, false
}

func trapcodeToByte(trapOp string) (uint16, error) {
	var err error

	switch trapOp {
	case "GETC": // Get char from keyboard, not echoed onto the terminal
		return 0, err
	case "OUT": // Output char
		return 0, err
	case "PUTS": // Output a word string
		return 0, err
	case "IN": // Get char from keyboard, echoed onto the terminal
		return 0, err
	case "PUTSP": // Output a byte string
		return 0, err
	case "HALT": // Halt
		return 0, err
	}
	fmt.Printf("[ERROR] Unknown trap: %s", trapOp)
	return 0, err
}

func trapcodeToTrap(instruction uint16) string {
	switch instruction & 0xFF {
	case TRAPGETC: // Get char from keyboard, not echoed onto the terminal
		return "GETC"
	case TRAPOUT: // Output char
		return "OUT"
	case TRAPPUTS: // Output a word string
		return "PUTS"
	case TRAPIN: // Get char from keyboard, echoed onto the terminal
		return "IN"
	case TRAPPUTSP: // Output a byte string
		return "PUTSP"
	case TRAPHALT: // Halt
		return "HALT"
	}
	fmt.Printf("[ERROR] Unknown trap code: x%X", instruction&0xFF)
	return fmt.Sprintf("; [ERROR] Unknown trap code: x%X", instruction&0xFF)
}
