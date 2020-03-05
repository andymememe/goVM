package vm

import "fmt"

/*
Add two number and save to register

 * 0001 DR(3) SR1(3) 0 00 SR2(3) => DR = SR1 + SR2
 * 0001 DR(3) SR1(3) 1 IMM(5)    => DR = SR1 + IMM
*/
func (c *CPU) add(instruction uint16) {
	r0 := (instruction >> 9) & 0x7
	r1 := (instruction >> 6) & 0x7
	immFlag := (instruction >> 5) & 0x1

	if immFlag != 0 { // R0 = R1 + to-16bit(imm-5bit)
		imm := signExtend(instruction&0x1F, 5)
		c.register[r0] = c.register[r1] + imm
	} else { // R0 = R1 + R2
		r2 := instruction & 0x7
		c.register[r0] = c.register[r1] + c.register[r2]
	}

	// Update Flags by result (R0)
	c.updateFlags(r0)
}

/*
Bitwise and two number and save to register

 * 0101 DR(3) SR1(3) 0 00 SR2(3) => DR = SR1 & SR2
 * 0101 DR(3) SR1(3) 1 IMM(5)    => DR = SR1 & IMM
*/
func (c *CPU) and(instruction uint16) {
	r0 := (instruction >> 9) & 0x7
	r1 := (instruction >> 6) & 0x7
	immFlag := (instruction >> 5) & 0x1

	if immFlag != 0 { // R0 = R1 & to-16bit(imm-5bit)
		imm := signExtend(instruction&0x1F, 5)
		c.register[r0] = c.register[r1] & imm
	} else { // R0 = R1 & R2
		r2 := instruction & 0x7
		c.register[r0] = c.register[r1] & c.register[r2]
	}

	// Update Flags by result (R0)
	c.updateFlags(r0)
}

/*
Bitwise not the number in a register and save to register

 * 1001 DR(3) SR(3) 1 11111 => DR = ^SR
*/
func (c *CPU) not(instruction uint16) {
	r0 := (instruction >> 9) & 0x7
	r1 := (instruction >> 6) & 0x7

	// R0 = ^R1
	c.register[r0] = ^c.register[r1]

	// Update Flags by result (R0)
	c.updateFlags(r0)
}

/*
Branch

If ((n And N) Or (z And Z) Or (p And P)) then Reg[PC] = Reg[PC] + PCOffset

 * 0000 N(1) Z(1) P(1) PCOffset(9) => If Condition(N, Z, R) match REG[COND] Then Reg[PC] + PCOffset
*/
func (c *CPU) br(instruction uint16) {
	pcOffset := signExtend(instruction&0x1ff, 9) // to-16bit(pcOffset-9bit)
	condFlag := (instruction >> 9) & 0x7         // condition

	if (condFlag & c.register[RCOND]) != 0 { // condition & Reg[RCOND] != 0
		c.register[RPC] += pcOffset
	}
}

/*
Jump, also Return

Return is a special case of Jump whenever BaseR is R7

 * 1100 000 BaseR(3) 000000 => REG[PC] = REG[BaseR] (JMP)
 * 1100 000 111      000000 => REG[PC] = REG[R7]    (RET)
*/
func (c *CPU) jmp(instruction uint16) {
	r1 := (instruction >> 6) & 0x7
	c.register[RPC] = c.register[r1] // Reg[RPC] = Reg[r1]
}

/*
Jump register

Save Reg[PC] to R7, and jump to somewhere

 * 0100 1 PCOffset(11)    => R7 = REG[PC], REG[PC] = REG[PC] + PCOffset
 * 0100 0 00 BaseR 000000 => R7 = REG[PC], REG[PC] = BaseR
*/
func (c *CPU) jsr(instruction uint16) {
	r1 := (instruction >> 6) & 0x7
	longPCOffset := signExtend(instruction&0x7ff, 11) // to-16bit(PCOffset-11bit)
	longFlag := (instruction >> 11) & 1

	c.register[RR7] = c.register[RPC]
	if longFlag != 0 {
		c.register[RPC] += longPCOffset /* JSR */
	} else {
		c.register[RPC] = c.register[r1] /* JSRR */
	}
}

/*
Load value in memory into register

 * 0010 DR(3) PCOffset(9) => DR = MEM[PCOffset + REG[PC]]
*/
func (c *CPU) ld(instruction uint16) {
	r0 := (instruction >> 9) & 0x7
	pcOffset := signExtend(instruction&0x1ff, 9) // to-16bit(pcOffset-9bit)

	// R0 = Mem[pcOffset + Reg[RPC]]
	c.register[r0] = c.memory.MemoryRead(c.register[RPC] + pcOffset)
	c.updateFlags(r0)
}

/*
Indirectly load value in memory into register

The memory requires 16 bits to address.
LDI is useful for loading values that are stored in locations far away from the current PC.

To use it, the address of the final location needs to be stored in a neighborhood nearby.

 * 1010 DR(3) PCOffset(9) => DR = MEM[MEM[PCOffset + REG[PC]]]
*/
func (c *CPU) ldi(instruction uint16) {
	r0 := (instruction >> 9) & 0x7
	pcOffset := signExtend(instruction&0x1ff, 9) // to-16bit(pcOffset-9bit)

	// R0 = Mem[Mem[pcOffset + Reg[RPC]]]
	c.register[r0] = c.memory.MemoryRead(c.memory.MemoryRead(c.register[RPC] + pcOffset))
	c.updateFlags(r0)
}

/*
Load value in memory of addr in a register into register

 * 0110 DR(3) BaseR(3) Offset(6) => DR = MEM[Offset + REG[BaseR]]
*/
func (c *CPU) ldr(instruction uint16) {
	r0 := (instruction >> 9) & 0x7
	r1 := (instruction >> 6) & 0x7
	offset := signExtend(instruction&0x3f, 6) // to-16bit(offset-6bit)

	// R0 = Mem[offset + Reg[r1]]
	c.register[r0] = c.memory.MemoryRead(c.register[r1] + offset)
	c.updateFlags(r0)
}

/*
Load effective address

 * 1110 DR(3) PCOffset (9) => DR = REG[PC] + PCOffset
*/
func (c *CPU) lea(instruction uint16) {
	r0 := (instruction >> 9) & 0x7
	pcOffset := signExtend(instruction&0x1ff, 9) // to-16bit(pcOffset-9bit)

	// R0 = pcOffset + Reg[RPC]
	c.register[r0] = c.register[RPC] + pcOffset
	c.updateFlags(r0)
}

/*
Store value in register to memory

 * 0011 SR(3) PCOffset(9) => MEM[REG[PC] + PCOffset] = REG[SR]
*/
func (c *CPU) st(instruction uint16) {
	r0 := (instruction >> 9) & 0x7
	pcOffset := signExtend(instruction&0x1ff, 9) // to-16bit(pcOffset-9bit)

	// Mem[pcOffset + Reg[RPC]] = Reg[r0]
	c.memory.MemoryWrite(c.register[RPC]+pcOffset, c.register[r0])
}

/*

Indirectly store value in register to memory

 * 1011 SR(3) PCOffset(9) => MEM[MEM[REG[PC] + PCOffset]] = REG[SR]
*/
func (c *CPU) sti(instruction uint16) {
	r0 := (instruction >> 9) & 0x7
	pcOffset := signExtend(instruction&0x1ff, 9) // to-16bit(pcOffset-9bit)

	// Mem[Mem[pcOffset + Reg[RPC]]] = Reg[r0]
	c.memory.MemoryWrite(c.memory.MemoryRead(c.register[RPC]+pcOffset), c.register[r0])
}

/*
Store value in register into memory whose address in a register

 * 0111 SR(3) BaseR(3) Offset(6) => MEM[BaseR + Offset] = SR
*/
func (c *CPU) str(instruction uint16) {
	r0 := (instruction >> 9) & 0x7
	r1 := (instruction >> 6) & 0x7
	offset := signExtend(instruction&0x3F, 6) // to-16bit(offset-6bit)

	// Mem[Reg[r1] + offset] = Reg[r0]
	c.memory.MemoryWrite(c.register[r1]+offset, c.register[r0])
}

/*
Bad OpCode and panic.

RES, RTI are unused, so panic too.
*/
func (c *CPU) badOpCode(instruction uint16) {
	panic(fmt.Sprintf("Bad OpCode: %b", instruction))
}
