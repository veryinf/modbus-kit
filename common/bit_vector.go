package common

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
func (bv *BitVector) Set(index uint, state bool) {
	if index >= bv.size {
		panic("index out of range")
	}
	// 计算位所在的word索引和bit位置
	word := index / 64
	bit := index % 64
	if state {
		bv.bits[word] |= 1 << bit
	} else {
		bv.bits[word] &= ^(1 << bit)
	}
}

// Get 获取指定索引位置的位的状态
func (bv *BitVector) Get(index uint) bool {
	if index >= bv.size {
		panic("index out of range")
	}
	// 计算位所在的word索引和bit位置
	word := index / 64
	bit := index % 64
	return (bv.bits[word] & (1 << bit)) != 0
}

func (bv *BitVector) Size() uint {
	return bv.size
}

func (bv *BitVector) Load(dataBuffer []byte) {
	bitIndex := uint(0)
	for _, b := range dataBuffer {
		if bitIndex >= bv.size {
			break
		}
		for i := uint(0); i < 8; i++ {
			if bitIndex >= bv.size {
				break
			}
			bit := (b >> i) & 1
			bv.Set(bitIndex, bit == 1)
			bitIndex++
		}
	}
}

func (bv *BitVector) ToString() string {
	str := make([]byte, bv.size)
	for i := uint(0); i < bv.size; i++ {
		if bv.Get(i) {
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
		if bv.Get(i) {
			byteIndex := i / 8
			bitInByte := i % 8
			data[byteIndex] |= 1 << bitInByte
		}
	}

	return data
}
