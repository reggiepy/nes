package nes

import "log"

type Memory interface {
	Read(address uint16) byte
	Read16(address uint16) uint16
	Write(address uint16, value byte)
}

// CPU Memory Map

type cpuMemory struct {
	nes *NES
}

func NewCPUMemory(nes *NES) Memory {
	return &cpuMemory{nes}
}

func (mem *cpuMemory) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		return mem.nes.RAM[address%0x0800]
	case address < 0x4000:
		return mem.nes.PPU.readRegister(0x2000 + address%8)
	case address == 0x4014:
		return mem.nes.PPU.readRegister(address)
	case address == 0x4016:
		return mem.nes.Controller1.Read()
	case address == 0x4017:
		return mem.nes.Controller2.Read()
	case address >= 0x6000:
		return mem.nes.Cartridge.Read(address)
	default:
		log.Fatalf("unhandled cpu memory read at address: 0x%04X", address)
	}
	return 0
}

func (mem *cpuMemory) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		mem.nes.RAM[address%0x0800] = value
	case address < 0x4000:
		mem.nes.PPU.writeRegister(0x2000+address%8, value)
	case address == 0x4014:
		mem.nes.PPU.writeRegister(address, value)
	case address == 0x4016:
		mem.nes.Controller1.Write(value)
	case address == 0x4017:
		mem.nes.Controller2.Write(value)
	case address < 0x4020:
		// TODO: I/O registers
	case address >= 0x6000:
		mem.nes.Cartridge.Write(address, value)
	default:
		log.Fatalf("unhandled cpu memory write at address: 0x%04X", address)
	}
}

func (mem *cpuMemory) Read16(address uint16) uint16 {
	lo := uint16(mem.Read(address))
	hi := uint16(mem.Read(address + 1))
	return hi<<8 | lo
}

// PPU Memory Map

type ppuMemory struct {
	nes *NES
}

func NewPPUMemory(nes *NES) Memory {
	return &ppuMemory{nes}
}

func (mem *ppuMemory) Read(address uint16) byte {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		return mem.nes.Cartridge.Read(address)
	case address < 0x3F00:
		address = mem.nes.Cartridge.NameTableAddress(address)
		return mem.nes.PPU.nameTableData[address%2048]
	case address < 0x4000:
		return mem.nes.PPU.readPalette(address % 32)
	default:
		log.Fatalf("unhandled ppu memory read at address: 0x%04X", address)
	}
	return 0
}

func (mem *ppuMemory) Write(address uint16, value byte) {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		mem.nes.Cartridge.Write(address, value)
	case address < 0x3F00:
		address = mem.nes.Cartridge.NameTableAddress(address)
		mem.nes.PPU.nameTableData[address%2048] = value
	case address < 0x4000:
		mem.nes.PPU.writePalette(address%32, value)
	default:
		log.Fatalf("unhandled ppu memory write at address: 0x%04X", address)
	}
}

func (mem *ppuMemory) Read16(address uint16) uint16 {
	lo := uint16(mem.Read(address))
	hi := uint16(mem.Read(address + 1))
	return hi<<8 | lo
}
