package slave

import (
	"sync"
)

type PointType string

const (
	PointTypeCoil            PointType = "coil"
	PointTypeDiscreteInput   PointType = "discrete_input"
	PointTypeHoldingRegister PointType = "holding_register"
	PointTypeInputRegister   PointType = "input_register"
)

// Point 类型
type Point struct {
	Address uint16 // 地址
	Value   uint16
	Type    PointType
}

// PointWriteCallback 事件回调函数类型
type PointWriteCallback func(event Point)

// MemoryDataStore 基于内存的数据存储实现
type MemoryDataStore struct {
	mu                  sync.RWMutex
	coils               map[uint16]bool
	discreteInputs      map[uint16]bool
	holdingRegisters    map[uint16]uint16
	inputRegisters      map[uint16]uint16
	eventWriteCallbacks []PointWriteCallback // 事件回调列表
}

// NewMemoryDataStore 创建新的内存数据存储
func NewMemoryDataStore() *MemoryDataStore {
	return &MemoryDataStore{
		coils:               make(map[uint16]bool),
		discreteInputs:      make(map[uint16]bool),
		holdingRegisters:    make(map[uint16]uint16),
		inputRegisters:      make(map[uint16]uint16),
		eventWriteCallbacks: make([]PointWriteCallback, 0),
	}
}

// AddWriteEventCallback 添加事件回调函数
func (m *MemoryDataStore) AddWriteEventCallback(callback PointWriteCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.eventWriteCallbacks = append(m.eventWriteCallbacks, callback)
}

// triggerWriteEvent 触发事件回调
func (m *MemoryDataStore) triggerWriteEvent(address uint16, value uint16, valueType PointType) {
	if len(m.eventWriteCallbacks) == 0 {
		return
	}
	event := Point{
		Address: address,
		Value:   value,
		Type:    valueType,
	}

	// 创建回调函数副本以避免在锁内执行回调
	m.mu.RLock()
	callbacks := make([]PointWriteCallback, len(m.eventWriteCallbacks))
	copy(callbacks, m.eventWriteCallbacks)
	m.mu.RUnlock()

	for _, callback := range callbacks {
		callback(event)
	}
}

// Read 根据类型直接读取单个值
func (m *MemoryDataStore) Read(pointType PointType, address uint16) uint16 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	switch pointType {
	case PointTypeCoil:
		if m.coils[address] {
			return 1
		}
		return 0
	case PointTypeDiscreteInput:
		if m.discreteInputs[address] {
			return 1
		}
		return 0
	case PointTypeHoldingRegister:
		return m.holdingRegisters[address]
	case PointTypeInputRegister:
		return m.inputRegisters[address]
	default:
		return 0
	}
}

// Write 根据类型直接写入单个值
func (m *MemoryDataStore) Write(pointType PointType, address uint16, value uint16) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch pointType {
	case PointTypeCoil:
		m.coils[address] = value != 0
	case PointTypeDiscreteInput:
		m.discreteInputs[address] = value != 0
	case PointTypeHoldingRegister:
		m.holdingRegisters[address] = value
	case PointTypeInputRegister:
		m.inputRegisters[address] = value
	}
	m.triggerWriteEvent(address, value, pointType)
}

func (m *MemoryDataStore) GetAllPoints() []Point {
	m.mu.Lock()
	defer m.mu.Unlock()
	points := make([]Point, 0)

	for address, value := range m.coils {
		point := Point{
			Address: address,
			Value:   0,
			Type:    PointTypeCoil,
		}
		if value {
			point.Value = 1
		}
		points = append(points, point)
	}

	for address, value := range m.discreteInputs {
		point := Point{
			Address: address,
			Value:   0,
			Type:    PointTypeDiscreteInput,
		}
		if value {
			point.Value = 1
		}
		points = append(points, point)
	}
	for address, value := range m.holdingRegisters {
		points = append(points, Point{
			Address: address,
			Value:   value,
			Type:    PointTypeHoldingRegister,
		})
	}
	for address, value := range m.inputRegisters {
		points = append(points, Point{
			Address: address,
			Value:   value,
			Type:    PointTypeInputRegister,
		})
	}

	return points
}
