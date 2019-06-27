package main

import "fmt"

/* メモリ 1MB */
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

var registersName = [8]string{
	"EAX",
	"ECX",
	"EDX",
	"EBX",
	"ESP",
	"EBP",
	"ESI",
	"EDI",
}

type emulator struct {
	// 汎用レジスタ
	registers [registersCount]uint32

	// EFLAGSレジスタ
	eflgas uint32

	// メモリ
	memory *uint8

	// プログラムカウンタ
	eip uint32
}

// エミュレータを作成する
// C言語と違ってメモリの確保に関しては実装しなくてよいから楽といえば楽なのかな？？
func createEmu(eip, esp uint32) *emulator {

	emu := &emulator{}
	emu.eip = eip
	emu.registers[esp] = esp

	return emu

}

// 汎用レジスタとプログラムカウンタの値を出力
func (emu *emulator) dumpRegisters() {

	for i, s := range registersName {
		fmt.Printf("%s = %08x\n", s, emu.registers[i])
	}
	fmt.Printf("EIP = %08x\n", emu.eip)
}

// エミュレータのメモリアドレス + index の領域にアクセスする
// が、現状エミュレータのメモリ空間を実現できていない。
func (emu *emulator) getCode8(index int) uint32 {
}

func main() {

	emu := createEmu(1, 1)
	emu.dumpRegisters()
}
