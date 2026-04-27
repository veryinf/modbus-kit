package common

import "fmt"

// BitVector 是一个位向量结构，用于高效存储和操作大量布尔值
type BitVector struct {
	bits []uint64
	size uint
}

// NewBitVector 创建一个新的位向量
// 参数:
//
//	size - 位向量的大小（位数）
//
// 返回值:
//
//	*BitVector - 指向新创建的位向量的指针
func NewBitVector(size uint) *BitVector {
	n := (size + 63) / 64
	return &BitVector{
		bits: make([]uint64, n),
		size: size,
	}
}

func NewBitVectorFromBooleans(values []bool) *BitVector {
	bv := NewBitVector(uint(len(values)))
	for i, value := range values {
		bv.Set(uint(i), value)
	}
	return bv
}

// Set 设置指定索引位置的状态
func (bv *BitVector) Set(index uint, state bool) error {
	if index >= bv.size {
		return fmt.Errorf("index out of range")
	}
	// 计算位所在的word索引和bit位置
	word := index / 64
	bit := index % 64
	if state {
		bv.bits[word] |= 1 << bit
	} else {
		bv.bits[word] &= ^(1 << bit)
	}
	return nil
}

// Get 获取指定索引位置的位的状态
func (bv *BitVector) Get(index uint) (bool, error) {
	if index >= bv.size {
		return false, fmt.Errorf("index out of range")
	}
	// 计算位所在的word索引和bit位置
	word := index / 64
	bit := index % 64
	return (bv.bits[word] & (1 << bit)) != 0, nil
}

func (bv *BitVector) Size() uint {
	return bv.size
}

// Load 从字节数组加载数据到位向量
// 参数:
//
// dataBuffer - 字节数组
// bigEndian - 是否使用大端序，true 为大端序，false 为小端序
func (bv *BitVector) Load(dataBuffer []byte, bigEndian bool) {
	bitIndex := uint(0)
	for _, b := range dataBuffer {
		if bitIndex >= bv.size {
			break
		}
		if bigEndian {
			// 大端序：从最高位开始
			for i := uint(0); i < 8; i++ {
				if bitIndex >= bv.size {
					break
				}
				bit := (b >> (7 - i)) & 1
				_ = bv.Set(bitIndex, bit == 1)
				bitIndex++
			}
		} else {
			// 小端序：从最低位开始
			for i := uint(0); i < 8; i++ {
				if bitIndex >= bv.size {
					break
				}
				bit := (b >> i) & 1
				_ = bv.Set(bitIndex, bit == 1)
				bitIndex++
			}
		}
	}
}

func (bv *BitVector) ToString() string {
	str := make([]byte, bv.size)
	for i := uint(0); i < bv.size; i++ {
		b, _ := bv.Get(i)
		if b {
			str[i] = '1'
		} else {
			str[i] = '0'
		}
	}
	return string(str)
}

func (bv *BitVector) ToBytes() []byte {
	numBytes := (bv.size + 7) / 8
	data := make([]byte, numBytes)

	for i := uint(0); i < bv.size; i++ {
		b, _ := bv.Get(i)
		if b {
			byteIndex := i / 8
			bitInByte := i % 8
			data[byteIndex] |= 1 << bitInByte
		}
	}

	return data
}
