package vm

import (
	"fmt"
	"math"

	"github.com/zetamatta/go-getch"
)

// Memory is define the memory
type Memory struct {
	storage []uint16
}

// NewMemory return Memory
func NewMemory() *Memory {
	return &Memory{
		storage: make([]uint16, math.MaxUint16),
	}
}

// MemoryRead return the value of storage[addr]
func (m *Memory) MemoryRead(addr uint16) uint16 {
	if addr == MRKBSR {
		evt, err := getch.Within(1000) // Non blocking wait event
		if err != nil {
			panic(fmt.Sprintf("MemoryRead: %s", err.Error()))
		}

		if k := evt.Key; k != nil {
			m.storage[MRKBSR] = (1 << 15)
			m.storage[MRKBDR] = uint16(k.Rune) | k.Scan
		} else {
			m.storage[MRKBSR] = 0
		}
	}
	return m.storage[addr]
}

// MemoryWrite set the value of storage[addr] to val
func (m *Memory) MemoryWrite(addr uint16, val uint16) {
	m.storage[addr] = val
}
