package common

import "fmt"

type Register struct {
	byte1 byte
	byte2 byte
}

func NewRegister(buffer []byte) *Register {
	return &Register{
		byte1: buffer[0],
		byte2: buffer[1],
	}
}

func NewRegisterFromUInt16(value uint16) *Register {
	return &Register{
		byte1: byte(value >> 8),
		byte2: byte(value),
	}
}

func NewRegisters(buffer []byte) []*Register {
	registers := make([]*Register, len(buffer)/2)
	for i := 0; i < len(buffer)/2; i++ {
		registers[i] = NewRegister(buffer[i*2 : i*2+2])
	}
	return registers
}

func RegistersToBytes(registers []*Register) *[]byte {
	value := make([]byte, len(registers)*2)
	for i := 0; i < len(registers); i++ {
		value[i*2] = registers[i].byte1
		value[i*2+1] = registers[i].byte2
	}
	return &value
}

func (r *Register) ToHexString() string {
	return fmt.Sprintf("%02x%02x", r.byte1, r.byte2)
}

func (r *Register) ToBytes() []byte {
	return []byte{r.byte1, r.byte2}
}

// Value 将 Register 转换为 uint16 值
func (r *Register) Value() uint16 {
	return uint16(r.byte1)<<8 | uint16(r.byte2)
}
