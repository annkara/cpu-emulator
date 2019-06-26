package main

const memorySize = 1024 * 1024

const (
	eax = iota
	ecx
	edx
	ebx
	esp
	ebp
	esi
	edi
	registersCount
)

type emulator struct {

	registers [registersCount]uint32

	eflgas uint32

	memory *uint8

	eip uint32

}

func main(){
	println(memorySize)
	println(eax + ebx)
}