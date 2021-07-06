package main

func (c *Cpu) setNegative(val uint8) {
	if val&0x80 != 0 {
		c.n = 1
	} else {
		c.n = 0
	}
}

func (c *Cpu) setZero(val uint8) {
	if val == 0 {
		c.z = 1
	} else {
		c.z = 0
	}
}

func (c *Cpu) setCarry(val uint16) {
	if val >= 0x100 {
		c.c = 1
	} else {
		c.c = 0
	}
}

func (c *Cpu) setOverflow(left, right, val uint8) {
	if (1^(left^right))&(left^val)&0x80 != 0 {
		c.v = 1
	} else {
		c.v = 0
	}
}
