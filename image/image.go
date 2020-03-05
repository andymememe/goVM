package image

import (
	"fmt"
	"os"

	"goVM/vm"
)

// ReadImage get filename and read image
func ReadImage(imagePath string, mem *vm.Memory) {
	f, err := os.Open(imagePath)
	if err != nil {
		panic(fmt.Sprintf("ReadImage/os.Open: %s", err.Error()))
	}
	defer f.Close()

	err = readImageFile(f, mem)
	if err != nil {
		panic(fmt.Sprintf("ReadImage/readImageFile: %s", err.Error()))
	}
}

// ReadImageFile reads file and load into memory
func readImageFile(file *os.File, mem *vm.Memory) error {
	var origin uint16

	fInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fSize := fInfo.Size()
	if fSize%2 != 0 {
		fSize++
	}
	dataByte := make([]byte, fSize)

	// Read file
	readN, err := file.Read(dataByte)
	if err != nil {
		return err
	}
	fmt.Printf("Read %d byte\n", readN)

	// The origin tells us where in memory to place the image
	origin = getUint16(dataByte[0:2])
	addr := origin
	curByte := int64(2)

	// Read byte into memory
	for curByte < fSize {
		mem.MemoryWrite(addr, getUint16(dataByte[curByte:curByte+2]))
		curByte += 2
		addr++
	}

	return nil
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
