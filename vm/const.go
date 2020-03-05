package vm

// Register ID
const (
	// RR0 ~ RR7 are general purpose reg
	RR0 uint16 = iota
	RR1
	RR2
	RR3
	RR4
	RR5
	RR6
	RR7
	RPC   /* Program Counter */
	RCOND /* Condition Flags, program can check logical conditions. */
	RCOUNT
)

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

// Flag ID for RCOND
const (
	FLPOS uint16 = 1 << 0 /* P */
	FLZRO uint16 = 1 << 1 /* Z */
	FLNEG uint16 = 1 << 2 /* N */
)

// set the PC to starting position
// 0x3000 is the default
const (
	PCSTART uint16 = 0x3000
)

// Special Memory Addr ID
const (
	MRKBSR uint16 = 0xFE00 // Keyboard Status
	MRKBDR uint16 = 0xFE02 // Keyboard Data
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
