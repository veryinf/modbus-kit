package common

import (
	"encoding/binary"
)

// ProtocolDataUnit 数据帧定义
type ProtocolDataUnit struct {
	FunctionCode byte
	Data         []byte
}

func (u *ProtocolDataUnit) LoadData(value ...uint16) *ProtocolDataUnit {
	data := make([]byte, 2*len(value))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	u.Data = data
	return u
}

func (u *ProtocolDataUnit) Append(value ...byte) *ProtocolDataUnit {
	u.Data = append(u.Data, value...)
	return u
}

func (u *ProtocolDataUnit) ToBytes() []byte {
	data := make([]byte, 1+len(u.Data))
	data[0] = u.FunctionCode
	copy(data[1:], u.Data)
	return data
}
