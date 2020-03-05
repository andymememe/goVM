package vm

// CPU define the CPU
type CPU struct {
	IsRunning bool
	register  []uint16
	memory    *Memory
}

// NewCPU return a CPU
func NewCPU(mem *Memory) *CPU {
	reg := make([]uint16, RCOUNT)
	reg[RPC] = PCSTART

	return &CPU{
		IsRunning: true,
		register:  reg,
		memory:    mem,
	}
}

// Load data into memory
func (c *CPU) Load(addr uint16, val uint16) {
	c.memory.MemoryWrite(addr, val)
}

// Execute the opcode
func (c *CPU) Execute(instruction uint16) {
	op := instruction >> 12

	switch op {
	case OPADD: // Add
		c.add(instruction)
		break
	case OPAND: // And
		c.and(instruction)
		break
	case OPNOT: // Not
		c.not(instruction)
		break
	case OPBR: // Branch
		c.br(instruction)
		break
	case OPJMP: // Jump & Return
		c.jmp(instruction)
		break
	case OPJSR: // Jump register
		c.jsr(instruction)
		break
	case OPLD: // Load
		c.ld(instruction)
		break
	case OPLDI: // Load Indirect
		c.ldi(instruction)
		break
	case OPLDR: // Load register
		c.ldr(instruction)
		break
	case OPLEA: // Load Effective Address
		c.lea(instruction)
		break
	case OPST: // Store
		c.st(instruction)
		break
	case OPSTI:
		c.sti(instruction)
		break
	case OPSTR:
		c.str(instruction)
		break
	case OPTRAP:
		c.trap(instruction)
		break
	case OPRES:
	case OPRTI:
	default:
		c.badOpCode(instruction)
		break
	}
}

// NextInstruction return next instruction
func (c *CPU) NextInstruction() uint16 {
	addr := c.register[RPC]
	c.register[RPC]++
	return c.memory.MemoryRead(addr)
}

// Update condition flag
func (c *CPU) updateFlags(regID uint16) {
	if c.register[regID] == 0 {
		c.register[RCOND] = FLZRO
	} else if (c.register[regID] >> 15) != 0 {
		c.register[RCOND] = FLNEG
	} else {
		c.register[RCOND] = FLPOS
	}
}

// The immediate mode value needs to be added to a 16-bit number.
// For positive numbers, we can fill in 0's for the additional bits and the value is the same.
// For negative numbers, we have to fill 1.
func signExtend(x uint16, bitCount int) uint16 {
	if ((x >> (bitCount - 1)) & 1) != 0 {
		x |= (0xFFFF << bitCount)
	}
	return x
}
