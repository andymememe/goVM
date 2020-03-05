package main

import (
	"fmt"
	"os"
	
	"goVM/image"
	"goVM/vm"
)

// LC-3 Architecture
// Reference: https://justinmeiners.github.io/lc3-vm/
func main() {
	argNum := len(os.Args)
	if argNum < 2 {
		/* show usage string */
		fmt.Println("lc3-exe-file [image-file] ...")
		os.Exit(2)
	}

	fmt.Println("Read image...")
	mem := vm.NewMemory()
	for j := 1; j < argNum; j++ {
		image.ReadImage(os.Args[j], mem)
	}
	cpu := vm.NewCPU(mem)

	fmt.Println("CPU start...")
	for cpu.IsRunning {
		instr := cpu.NextInstruction()
		cpu.Execute(instr)
	}
}
