package main

import (
	"github.com/sirupsen/logrus"
)

type Cpu struct {
	pc uint16
	sp uint16

	a uint8
	x uint8
	y uint8

	c uint8
	n uint8
	z uint8
	v uint8

	mem [1 << 16]uint8

	ops [256]OpTableEntry

	cycles uint32
}

func (c *Cpu) readMem(addr uint16) uint8 {
	return c.mem[addr]
}

func (c *Cpu) writeMem(addr uint16, val uint8) {
	c.mem[addr] = val
}

func (c *Cpu) Reset() {
	c.a = 0
	c.x = 0
	c.y = 0

	// TODO(Ben): Implement proper reset vector handling
	c.pc = 0

	c.initOpTable()

	// Test program - calculates 7th Fibonacci number
	fib := []uint8{
		0xA2, 0x01, 0x86, 0x00, 0x38, 0xA0, 0x07, 0x98, 0xE9, 0x03, 0xA8, 0x18, 0xA9, 0x02, 0x85, 0x01,
		0xA6, 0x01, 0x65, 0x00, 0x85, 0x01, 0x86, 0x00, 0x88, 0xD0, 0xF5,
	}
	for i, s := range fib {
		c.mem[i] = s
	}
}

func (c *Cpu) Run() {
	for {
		opcode := c.readMem(c.pc)

		// TODO(Ben): Implement proper behaviour
		if opcode == 0x00 {
			break
		}

		op := c.ops[opcode]
		addr := op.mode(c.pc)
		op.op(addr)
		c.cycles += uint32(op.cycles)
		logrus.Infof("CPU A=0x%02X X=0x%002X Y=0x%02X PC=0x%04X SP=0x%04X", c.a, c.x, c.y, c.pc, c.sp)
		c.pc += uint16(op.bytes)
	}

	logrus.Println("CPU A:", c.a)
}

func main() {
	cpu := Cpu{}

	cpu.Reset()
	cpu.Run()
}
