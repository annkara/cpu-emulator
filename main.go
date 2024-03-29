package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

/* メモリ 1MB */
const memorySize = 1024 * 1024

const (
	EAX = iota
	ECX
	EDX
	EBX
	ESP
	EBP
	ESI
	EDI
	REGISTERSCOUNT
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

/* エミュレータ構造体 */
type emulator struct {
	// 汎用レジスタ
	registers [REGISTERSCOUNT]uint32

	// EFLAGSレジスタ
	eflgas uint32

	// 8bit * memorySize のメモリ空間
	memory [memorySize]byte

	// プログラムカウンタ 32bit
	eip uint32
}

// 命令テーブル
var instructions [256]func(*emulator)

// エミュレータを作成する
func createEmu(eip, esp uint32) *emulator {

	emu := &emulator{}
	emu.eip = eip
	emu.registers[ESP] = esp

	return emu

}

// 汎用レジスタとプログラムカウンタの値を出力
func (emu *emulator) dumpRegisters() {

	for i, s := range registersName {
		fmt.Printf("%s = %08x\n", s, emu.registers[i])
	}
	fmt.Printf("EIP = %08x\n", emu.eip)
}

func code8(emu *emulator, index uint32) uint8 {
	return emu.memory[emu.eip+index]
}

func signCode8(emu *emulator, index uint32) int8 {
	return int8(emu.memory[emu.eip+index])
}

func code32(emu *emulator, index uint32) uint32 {

	ret := emu.memory[emu.eip+index:]
	return binary.LittleEndian.Uint32(ret)
}

func movR32Imm32(emu *emulator) {

	reg := code8(emu, 0) - 0xB8
	value := code32(emu, 1)
	emu.registers[reg] = value
	emu.eip += 5

}

func shortJump(emu *emulator) {

	diff := signCode8(emu, 1)
	emu.eip += uint32(diff + 2)
}

func nearJump(emu *emulator) {

	diff := code32(emu, 1)
	emu.eip += (diff + 5)
}

// 関数テーブルの初期化
func initInstructions() {
	for i := 0; i < 8; i++ {
		instructions[0xB8+i] = movR32Imm32
	}
	instructions[0xE9] = nearJump
	instructions[0xEB] = shortJump
}

func main() {

	if len(os.Args) != 2 {
		println("usage: px86 filename\n")
		return
	}

	// EIP、ESPが0x7c00の状態のエミュレータを作る
	emu := createEmu(0x7c00, 0x7c00)

	// 機械語ファイルを読み込んで、エミュレータのメモリ上に格納する
	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// 0x200 = 512バイトのスライスを作成し
	// バイナリファイルを読み込み、copy関数でエミュレータのメモリに格納
	data := make([]byte, 0x200)
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	copy(emu.memory[0x7c00:0x7c00+count], data[:count])

	initInstructions()

	for emu.eip < memorySize {
		code := code8(emu, 0)

		fmt.Printf("EIP = %X, Code = %02X\n", emu.eip, code)

		if instructions[code] == nil {
			fmt.Printf("\n\nNot Implemented: %x\n", code)
			break
		}

		instructions[code](emu)

		if emu.eip == 0x00 {
			println("\n\nend of program.\n\n")
			break
		}
	}
	emu.dumpRegisters()
	return
}
