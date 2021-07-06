package main

type AddressingMode func(pc uint16) uint16

func (c *Cpu) AddrModeAccumulator(pc uint16) uint16 {
	return 0 // Source address is not used
}

func (c *Cpu) AddrModeImmediate(pc uint16) uint16 {
	return c.pc + 1
}

func (c *Cpu) AddrModeAbsolute(pc uint16) uint16 {
	addr := uint16(c.readMem(pc+2))<<8 + uint16(c.readMem(pc+1))
	return addr
}

func (c *Cpu) AddrModeAbsoluteX(pc uint16) uint16 {
	return c.AddrModeAbsolute(pc) + uint16(c.x) + uint16(c.c)
}

func (c *Cpu) AddrModeAbsoluteY(pc uint16) uint16 {
	return c.AddrModeAbsolute(pc) + uint16(c.y) + uint16(c.c)
}

func (c *Cpu) AddrModeZeroPage(pc uint16) uint16 {
	return uint16(c.readMem(pc + 1))
}

func (c *Cpu) AddrModeZeroPageX(pc uint16) uint16 {
	addr := c.AddrModeZeroPage(pc)
	return addr + uint16(c.x)
}

func (c *Cpu) AddrModeZeroPageY(pc uint16) uint16 {
	addr := c.AddrModeZeroPage(pc)
	return addr + uint16(c.y)
}

func (c *Cpu) AddrModeImplied(pc uint16) uint16 {
	return 0
}

func (c *Cpu) AddrModeRelative(pc uint16) uint16 {
	rel := int8(c.readMem(pc + 1))
	addr := uint16(int32(pc) + int32(rel))
	return addr
}

func (c *Cpu) AddrModeIndirect(pc uint16) uint16 {
	ptr := c.AddrModeAbsolute(pc)
	addr := uint16(c.readMem(ptr+1))<<8 + uint16(c.readMem(ptr))
	return addr
}

func (c *Cpu) AddrModeIndirectX(pc uint16) uint16 {
	ptr := c.AddrModeZeroPage(pc) + uint16(c.x)
	addr := uint16(c.readMem(ptr+1))<<8 + uint16(c.readMem(ptr))
	return addr
}

func (c *Cpu) AddrModeIndirectY(pc uint16) uint16 {
	ptr := c.AddrModeZeroPage(pc)
	addr := uint16(c.readMem(ptr+1))<<8 + uint16(c.readMem(ptr))
	addr += uint16(c.c)
	return addr
}

type Op func(addr uint16)

type OpTableEntry struct {
	opcode uint8
	mode   AddressingMode
	op     Op
	bytes  uint8
	cycles uint8
}

func (c *Cpu) initOpTable() {
	// LDA
	c.ops[0xA9] = OpTableEntry{opcode: 0xA9, mode: c.AddrModeImmediate, op: c.Op_LDA, bytes: 2, cycles: 2}
	c.ops[0xA5] = OpTableEntry{opcode: 0xA5, mode: c.AddrModeZeroPage, op: c.Op_LDA, bytes: 2, cycles: 3}
	c.ops[0xB9] = OpTableEntry{opcode: 0xB9, mode: c.AddrModeZeroPageX, op: c.Op_LDA, bytes: 2, cycles: 4}
	c.ops[0xAD] = OpTableEntry{opcode: 0xAD, mode: c.AddrModeAbsolute, op: c.Op_LDA, bytes: 3, cycles: 4}
	c.ops[0xBD] = OpTableEntry{opcode: 0xBD, mode: c.AddrModeAbsoluteX, op: c.Op_LDA, bytes: 3, cycles: 4}
	c.ops[0xB9] = OpTableEntry{opcode: 0xB9, mode: c.AddrModeAbsoluteY, op: c.Op_LDA, bytes: 3, cycles: 4}
	c.ops[0xA1] = OpTableEntry{opcode: 0xA1, mode: c.AddrModeIndirectX, op: c.Op_LDA, bytes: 2, cycles: 6}
	c.ops[0xB1] = OpTableEntry{opcode: 0xB1, mode: c.AddrModeIndirectY, op: c.Op_LDA, bytes: 2, cycles: 5}

	// LDX
	c.ops[0xA2] = OpTableEntry{opcode: 0xA2, mode: c.AddrModeImmediate, op: c.Op_LDX, bytes: 2, cycles: 2}
	c.ops[0xA6] = OpTableEntry{opcode: 0xA6, mode: c.AddrModeZeroPage, op: c.Op_LDX, bytes: 2, cycles: 3}
	c.ops[0xB6] = OpTableEntry{opcode: 0xB6, mode: c.AddrModeZeroPageY, op: c.Op_LDX, bytes: 2, cycles: 4}
	c.ops[0xAE] = OpTableEntry{opcode: 0xAE, mode: c.AddrModeAbsolute, op: c.Op_LDX, bytes: 3, cycles: 4}
	c.ops[0xBE] = OpTableEntry{opcode: 0xBE, mode: c.AddrModeAbsoluteY, op: c.Op_LDX, bytes: 3, cycles: 4}

	// LDY
	c.ops[0xA0] = OpTableEntry{opcode: 0xA0, mode: c.AddrModeImmediate, op: c.Op_LDY, bytes: 2, cycles: 2}
	c.ops[0xA4] = OpTableEntry{opcode: 0xA4, mode: c.AddrModeZeroPage, op: c.Op_LDY, bytes: 2, cycles: 3}
	c.ops[0xB4] = OpTableEntry{opcode: 0xB4, mode: c.AddrModeZeroPageX, op: c.Op_LDY, bytes: 2, cycles: 4}
	c.ops[0xAC] = OpTableEntry{opcode: 0xAC, mode: c.AddrModeAbsolute, op: c.Op_LDY, bytes: 3, cycles: 4}
	c.ops[0xBC] = OpTableEntry{opcode: 0xBC, mode: c.AddrModeAbsoluteX, op: c.Op_LDY, bytes: 3, cycles: 4}

	// STA
	c.ops[0x85] = OpTableEntry{opcode: 0x85, mode: c.AddrModeZeroPage, op: c.Op_STA, bytes: 2, cycles: 3}
	c.ops[0x95] = OpTableEntry{opcode: 0x95, mode: c.AddrModeZeroPageX, op: c.Op_STA, bytes: 2, cycles: 4}
	c.ops[0x8D] = OpTableEntry{opcode: 0x8D, mode: c.AddrModeAbsolute, op: c.Op_STA, bytes: 3, cycles: 4}
	c.ops[0x9D] = OpTableEntry{opcode: 0x9D, mode: c.AddrModeAbsoluteX, op: c.Op_STA, bytes: 3, cycles: 5}
	c.ops[0x99] = OpTableEntry{opcode: 0x99, mode: c.AddrModeAbsoluteY, op: c.Op_STA, bytes: 3, cycles: 5}
	c.ops[0x81] = OpTableEntry{opcode: 0x81, mode: c.AddrModeIndirectX, op: c.Op_STA, bytes: 2, cycles: 6}
	c.ops[0x91] = OpTableEntry{opcode: 0x91, mode: c.AddrModeIndirectY, op: c.Op_STA, bytes: 2, cycles: 6}

	// STX
	c.ops[0x86] = OpTableEntry{opcode: 0x86, mode: c.AddrModeZeroPage, op: c.Op_STX, bytes: 2, cycles: 3}
	c.ops[0x96] = OpTableEntry{opcode: 0x96, mode: c.AddrModeZeroPageY, op: c.Op_STX, bytes: 2, cycles: 4}
	c.ops[0x8E] = OpTableEntry{opcode: 0x8E, mode: c.AddrModeAbsolute, op: c.Op_STX, bytes: 3, cycles: 4}

	// STY
	c.ops[0x84] = OpTableEntry{opcode: 0x84, mode: c.AddrModeZeroPage, op: c.Op_STY, bytes: 2, cycles: 3}
	c.ops[0x94] = OpTableEntry{opcode: 0x94, mode: c.AddrModeZeroPageX, op: c.Op_STY, bytes: 2, cycles: 4}
	c.ops[0x8C] = OpTableEntry{opcode: 0x8C, mode: c.AddrModeAbsolute, op: c.Op_STY, bytes: 3, cycles: 4}

	// SEC
	c.ops[0x38] = OpTableEntry{opcode: 0x38, mode: c.AddrModeImplied, op: c.Op_SEC, bytes: 1, cycles: 2}

	// CLC
	c.ops[0x18] = OpTableEntry{opcode: 0x18, mode: c.AddrModeImplied, op: c.Op_CLC, bytes: 1, cycles: 2}

	// TAY
	c.ops[0xA8] = OpTableEntry{opcode: 0xA8, mode: c.AddrModeImplied, op: c.Op_TAY, bytes: 1, cycles: 2}

	// TYA
	c.ops[0x98] = OpTableEntry{opcode: 0x98, mode: c.AddrModeImplied, op: c.Op_TYA, bytes: 1, cycles: 2}

	// SBC
	c.ops[0xE9] = OpTableEntry{opcode: 0xE9, mode: c.AddrModeImmediate, op: c.Op_SBC, bytes: 2, cycles: 2}
	c.ops[0xE5] = OpTableEntry{opcode: 0xE5, mode: c.AddrModeZeroPage, op: c.Op_SBC, bytes: 2, cycles: 3}
	c.ops[0xF5] = OpTableEntry{opcode: 0xF5, mode: c.AddrModeZeroPageX, op: c.Op_SBC, bytes: 2, cycles: 4}
	c.ops[0xED] = OpTableEntry{opcode: 0xED, mode: c.AddrModeAbsolute, op: c.Op_SBC, bytes: 3, cycles: 4}
	c.ops[0xFD] = OpTableEntry{opcode: 0xFD, mode: c.AddrModeAbsoluteX, op: c.Op_SBC, bytes: 3, cycles: 4}
	c.ops[0xF9] = OpTableEntry{opcode: 0xF9, mode: c.AddrModeAbsoluteY, op: c.Op_SBC, bytes: 3, cycles: 4}
	c.ops[0xE1] = OpTableEntry{opcode: 0xE1, mode: c.AddrModeIndirectX, op: c.Op_SBC, bytes: 2, cycles: 6}
	c.ops[0xF1] = OpTableEntry{opcode: 0xE1, mode: c.AddrModeIndirectY, op: c.Op_SBC, bytes: 2, cycles: 5}

	// ADC
	c.ops[0x69] = OpTableEntry{opcode: 0x69, mode: c.AddrModeImmediate, op: c.Op_ADC, bytes: 2, cycles: 2}
	c.ops[0x65] = OpTableEntry{opcode: 0x65, mode: c.AddrModeZeroPage, op: c.Op_ADC, bytes: 2, cycles: 3}
	c.ops[0x75] = OpTableEntry{opcode: 0x75, mode: c.AddrModeZeroPageX, op: c.Op_ADC, bytes: 2, cycles: 4}
	c.ops[0x6D] = OpTableEntry{opcode: 0x6D, mode: c.AddrModeAbsolute, op: c.Op_ADC, bytes: 3, cycles: 4}
	c.ops[0x7D] = OpTableEntry{opcode: 0x7D, mode: c.AddrModeAbsoluteX, op: c.Op_ADC, bytes: 3, cycles: 4}
	c.ops[0x79] = OpTableEntry{opcode: 0x79, mode: c.AddrModeAbsoluteY, op: c.Op_ADC, bytes: 3, cycles: 4}
	c.ops[0x61] = OpTableEntry{opcode: 0x61, mode: c.AddrModeIndirectX, op: c.Op_ADC, bytes: 2, cycles: 6}
	c.ops[0x71] = OpTableEntry{opcode: 0x71, mode: c.AddrModeIndirectY, op: c.Op_ADC, bytes: 2, cycles: 5}

	// DEY
	c.ops[0x88] = OpTableEntry{opcode: 0x88, mode: c.AddrModeImplied, op: c.Op_DEY, bytes: 1, cycles: 2}

	// BNE
	c.ops[0xD0] = OpTableEntry{opcode: 0xD0, mode: c.AddrModeRelative, op: c.Op_BNE, bytes: 2, cycles: 2}
}

func (c *Cpu) Op_LDA(addr uint16) {
	val := c.readMem(addr)
	c.setNegative(val)
	c.setZero(val)
	c.a = val
}

func (c *Cpu) Op_LDX(addr uint16) {
	val := c.readMem(addr)
	c.setNegative(val)
	c.setZero(val)
	c.x = val
}

func (c *Cpu) Op_LDY(addr uint16) {
	val := c.readMem(addr)
	c.setNegative(val)
	c.setZero(val)
	c.y = val
}

func (c *Cpu) Op_STA(addr uint16) {
	c.writeMem(addr, c.a)
}

func (c *Cpu) Op_STX(addr uint16) {
	c.writeMem(addr, c.x)
}

func (c *Cpu) Op_STY(addr uint16) {
	c.writeMem(addr, c.x)
}

func (c *Cpu) Op_SEC(addr uint16) {
	c.c = 1
}

func (c *Cpu) Op_CLC(addr uint16) {
	c.c = 0
}

func (c *Cpu) Op_TAY(addr uint16) {
	val := c.a
	c.setNegative(val)
	c.setZero(val)
	c.y = val
}

func (c *Cpu) Op_TYA(addr uint16) {
	val := c.y
	c.setNegative(val)
	c.setZero(val)
	c.a = val
}

func (c *Cpu) Op_SBC(addr uint16) {
	var tmp uint16 = uint16(c.a) - uint16(c.readMem(addr)) - uint16(1^c.c)
	val := uint8(tmp & 0xFF)
	c.setNegative(val)
	c.setZero(val)
	c.setCarry(tmp)
	c.setOverflow(c.a, c.readMem(addr), val)
	c.a = val
}

func (c *Cpu) Op_ADC(addr uint16) {
	var tmp uint16 = uint16(c.a) + uint16(c.readMem(addr)) + uint16(c.c)
	val := uint8(tmp & 0xFF)
	c.setNegative(val)
	c.setZero(val)
	c.setCarry(tmp)
	c.setOverflow(c.a, c.readMem(addr), val)
	c.a = val
}

func (c *Cpu) Op_DEY(addr uint16) {
	val := c.y - 1
	c.setNegative(val)
	c.setZero(val)
	c.y = val
}

func (c *Cpu) Op_BNE(addr uint16) {
	if c.z != 0 {
		return
	}
	c.pc = addr
}
