package vm

import (
	"fmt"

	"github.com/zetamatta/go-getch"
)

/*
Trap Routine

Trap routine is predefined routine for performing common tasks and interacting with I/O devices.

When a trap code is called, the REG[PC] is moved to that code's address.
The CPU executes instructions of the procedure, and when it is complete,
the PC is reset to the location the trap was called from.

Note: This is why programs start at address 0x3000 instead of 0x0.
The lower addresses are left empty to leave space for the trap routine code.

 * 1111 0000 TrapVec(8)
*/
func (c *CPU) trap(instruction uint16) {
	switch instruction & 0xFF {
	case TRAPGETC: // Get char from keyboard, not echoed onto the terminal
		c.trapGetc()
		break
	case TRAPOUT: // Output char
		c.trapOut()
		break
	case TRAPPUTS: // Output a word string
		c.trapPuts()
		break
	case TRAPIN: // Get char from keyboard, echoed onto the terminal
		c.trapIn()
		break
	case TRAPPUTSP: // Output a byte string
		c.trapPutsp()
		break
	case TRAPHALT: // Halt
		c.trapHalt()
		break
	}
}

/*
Input Char
*/
func (c *CPU) trapGetc() {
	c.register[RR0] = uint16(getch.Rune()) // Blocking get rune
}

/*
Output Char
*/
func (c *CPU) trapOut() {
	fmt.Printf("%c", rune(c.register[RR0]))
	err := getch.Flush()
	if err != nil {
		panic(fmt.Sprintf("trapPuts: %s", err.Error()))
	}
}

/*
Output a word string
*/
func (c *CPU) trapPuts() {
	addr := c.register[RR0]
	r := c.memory.MemoryRead(addr)
	for r != 0x0000 {
		fmt.Printf("%c", rune(r))
		addr++
		r = c.memory.MemoryRead(addr)
	}
	err := getch.Flush()
	if err != nil {
		panic(fmt.Sprintf("trapPuts: %s", err.Error()))
	}
}

/*
Prompt for Input Character
*/
func (c *CPU) trapIn() {
	fmt.Print("Enter a character: ")
	aRune := getch.Rune()
	fmt.Printf("%c", aRune)
	c.register[RR0] = uint16(aRune)
}

/*
Output String
*/
func (c *CPU) trapPutsp() {
	addr := c.register[RR0]
	r := c.memory.MemoryRead(addr)
	for r != 0x0000 {
		rune1 := rune(r & 0xFF)
		fmt.Printf("%c", rune1)
		rune2 := rune(r >> 8)
		if rune2 != rune(0x0000) {
			fmt.Printf("%c", rune2)
		}
		addr++
		r = c.memory.MemoryRead(addr)
	}
	err := getch.Flush()
	if err != nil {
		panic(fmt.Sprintf("trapPuts: %s", err.Error()))
	}
}

/*
Halt Program
*/
func (c *CPU) trapHalt() {
	err := getch.Flush()
	if err != nil {
		panic(fmt.Sprintf("trapPuts: %s", err.Error()))
	}
	c.IsRunning = false
}
