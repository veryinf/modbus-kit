package slave

import (
	"encoding/binary"
	"github.com/veryinf/modbus-kit/common"
)

type RequestHandler struct {
	DeviceInfo *DeviceInfo
	store      *MemoryDataStore
}

// HandleRequest 处理 Modbus 请求
func (s *RequestHandler) HandleRequest(request *common.ProtocolDataUnit) (response *common.ProtocolDataUnit, err error) {
	response = &common.ProtocolDataUnit{
		FunctionCode: request.FunctionCode,
	}

	switch request.FunctionCode {
	case common.FuncCodeReadCoils:
		response = s.handleReadCoils(request)
	case common.FuncCodeReadDiscreteInputs:
		response = s.handleReadDiscreteInputs(request)
	case common.FuncCodeReadHoldingRegisters:
		response = s.handleReadHoldingRegisters(request)
	case common.FuncCodeReadInputRegisters:
		response = s.handleReadInputRegisters(request)
	case common.FuncCodeWriteSingleCoil:
		response = s.handleWriteSingleCoil(request)
	case common.FuncCodeWriteSingleRegister:
		response = s.handleWriteSingleRegister(request)
	case common.FuncCodeWriteMultipleCoils:
		response = s.handleWriteMultipleCoils(request)
	case common.FuncCodeWriteMultipleRegisters:
		response = s.handleWriteMultipleRegisters(request)
	case common.FuncCodeReadDeviceIdentification:
		response = s.handleReadDeviceIdentification(request)
	default:
		response = &common.ProtocolDataUnit{
			FunctionCode: request.FunctionCode | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalFunction},
		}
		return response, nil
	}

	return response, nil
}

// handleReadCoils 处理读取线圈请求 (功能码 0x01)
func (s *RequestHandler) handleReadCoils(request *common.ProtocolDataUnit) *common.ProtocolDataUnit {
	if len(request.Data) < 4 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeReadCoils | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataAddress},
		}
	}

	address := binary.BigEndian.Uint16(request.Data[0:2])
	quantity := binary.BigEndian.Uint16(request.Data[2:4])

	if quantity < 1 || quantity > 2000 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeReadCoils | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataValue},
		}
	}
	bitVector := common.NewBitVector(uint(quantity))
	for i := uint16(0); i < quantity; i++ {
		value := s.store.Read(PointTypeCoil, address+i)
		bitVector.Set(uint(i), value != 0)
	}
	bytes := bitVector.ToBytes()

	// 构建响应: [字节数] [数据...]
	responseData := make([]byte, 1+len(bytes))
	responseData[0] = byte(len(bytes))
	copy(responseData[1:], bytes)

	return &common.ProtocolDataUnit{
		FunctionCode: common.FuncCodeReadCoils,
		Data:         responseData,
	}
}

// handleReadDiscreteInputs 处理读取离散输入请求 (功能码 0x02)
func (s *RequestHandler) handleReadDiscreteInputs(request *common.ProtocolDataUnit) *common.ProtocolDataUnit {
	if len(request.Data) < 4 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeReadDiscreteInputs | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataAddress},
		}
	}

	address := binary.BigEndian.Uint16(request.Data[0:2])
	quantity := binary.BigEndian.Uint16(request.Data[2:4])

	if quantity < 1 || quantity > 2000 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeReadDiscreteInputs | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataValue},
		}
	}

	bitVector := common.NewBitVector(uint(quantity))
	for i := uint16(0); i < quantity; i++ {
		value := s.store.Read(PointTypeDiscreteInput, address+i)
		bitVector.Set(uint(i), value != 0)
	}
	bytes := bitVector.ToBytes()

	responseData := make([]byte, 1+len(bytes))
	responseData[0] = byte(len(bytes))
	copy(responseData[1:], bytes)

	return &common.ProtocolDataUnit{
		FunctionCode: common.FuncCodeReadDiscreteInputs,
		Data:         responseData,
	}
}

// handleReadHoldingRegisters 处理读取保持寄存器请求 (功能码 0x03)
func (s *RequestHandler) handleReadHoldingRegisters(request *common.ProtocolDataUnit) *common.ProtocolDataUnit {
	if len(request.Data) < 4 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeReadHoldingRegisters | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataAddress},
		}
	}

	address := binary.BigEndian.Uint16(request.Data[0:2])
	quantity := binary.BigEndian.Uint16(request.Data[2:4])

	if quantity < 1 || quantity > 125 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeReadHoldingRegisters | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataValue},
		}
	}

	registers := make([]*common.Register, quantity)
	for i := uint16(0); i < quantity; i++ {
		registers[i] = common.NewRegisterFromUInt16(s.store.Read(PointTypeHoldingRegister, address+i))
	}

	bytes := common.RegistersToBytes(registers)

	responseData := make([]byte, 1+len(*bytes))
	responseData[0] = byte(len(*bytes))
	copy(responseData[1:], *bytes)

	return &common.ProtocolDataUnit{
		FunctionCode: common.FuncCodeReadHoldingRegisters,
		Data:         responseData,
	}
}

// handleReadInputRegisters 处理读取输入寄存器请求 (功能码 0x04)
func (s *RequestHandler) handleReadInputRegisters(request *common.ProtocolDataUnit) *common.ProtocolDataUnit {
	if len(request.Data) < 4 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeReadInputRegisters | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataAddress},
		}
	}

	address := binary.BigEndian.Uint16(request.Data[0:2])
	quantity := binary.BigEndian.Uint16(request.Data[2:4])

	if quantity < 1 || quantity > 125 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeReadInputRegisters | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataValue},
		}
	}

	registers := make([]*common.Register, quantity)
	for i := uint16(0); i < quantity; i++ {
		registers[i] = common.NewRegisterFromUInt16(s.store.Read(PointTypeInputRegister, address+i))
	}

	bytes := common.RegistersToBytes(registers)

	responseData := make([]byte, 1+len(*bytes))
	responseData[0] = byte(len(*bytes))
	copy(responseData[1:], *bytes)

	return &common.ProtocolDataUnit{
		FunctionCode: common.FuncCodeReadInputRegisters,
		Data:         responseData,
	}
}

// handleWriteSingleCoil 处理写入单个线圈请求 (功能码 0x05)
func (s *RequestHandler) handleWriteSingleCoil(request *common.ProtocolDataUnit) *common.ProtocolDataUnit {
	if len(request.Data) < 4 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeWriteSingleCoil | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataAddress},
		}
	}

	address := binary.BigEndian.Uint16(request.Data[0:2])
	value := binary.BigEndian.Uint16(request.Data[2:4])

	if value != 0x0000 && value != 0xFF00 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeWriteSingleCoil | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataValue},
		}
	}

	s.store.Write(PointTypeCoil, address, value)

	return &common.ProtocolDataUnit{
		FunctionCode: common.FuncCodeWriteSingleCoil,
		Data:         request.Data,
	}
}

// handleWriteSingleRegister 处理写入单个寄存器请求 (功能码 0x06)
func (s *RequestHandler) handleWriteSingleRegister(request *common.ProtocolDataUnit) *common.ProtocolDataUnit {
	if len(request.Data) < 4 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeWriteSingleRegister | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataAddress},
		}
	}

	address := binary.BigEndian.Uint16(request.Data[0:2])
	value := binary.BigEndian.Uint16(request.Data[2:4])

	s.store.Write(PointTypeHoldingRegister, address, value)

	return &common.ProtocolDataUnit{
		FunctionCode: common.FuncCodeWriteSingleRegister,
		Data:         request.Data,
	}
}

// handleWriteMultipleCoils 处理写入多个线圈请求 (功能码 0x0F)
func (s *RequestHandler) handleWriteMultipleCoils(request *common.ProtocolDataUnit) *common.ProtocolDataUnit {
	if len(request.Data) < 5 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeWriteMultipleCoils | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataAddress},
		}
	}

	address := binary.BigEndian.Uint16(request.Data[0:2])
	quantity := binary.BigEndian.Uint16(request.Data[2:4])
	byteCount := int(request.Data[4])

	if quantity < 1 || quantity > 1968 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeWriteMultipleCoils | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataValue},
		}
	}

	expectedByteCount := int((quantity + 7) / 8)
	if byteCount != expectedByteCount {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeWriteMultipleCoils | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataValue},
		}
	}

	if len(request.Data) < 5+byteCount {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeWriteMultipleCoils | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataAddress},
		}
	}

	bitVector := common.NewBitVector(uint(quantity))
	bitVector.Load(request.Data[5 : 5+byteCount])
	for i := uint16(0); i < quantity; i++ {
		var val uint16 = 0
		if bitVector.Get(uint(i)) {
			val = 1
		}
		s.store.Write(PointTypeCoil, address+i, val)
	}

	responseData := make([]byte, 4)
	binary.BigEndian.PutUint16(responseData[0:2], address)
	binary.BigEndian.PutUint16(responseData[2:4], quantity)

	return &common.ProtocolDataUnit{
		FunctionCode: common.FuncCodeWriteMultipleCoils,
		Data:         responseData,
	}
}

// handleWriteMultipleRegisters 处理写入多个寄存器请求 (功能码 0x10)
func (s *RequestHandler) handleWriteMultipleRegisters(request *common.ProtocolDataUnit) *common.ProtocolDataUnit {
	if len(request.Data) < 5 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeWriteMultipleRegisters | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataAddress},
		}
	}

	address := binary.BigEndian.Uint16(request.Data[0:2])
	quantity := binary.BigEndian.Uint16(request.Data[2:4])
	byteCount := int(request.Data[4])

	if quantity < 1 || quantity > 123 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeWriteMultipleRegisters | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataValue},
		}
	}

	expectedByteCount := int(quantity * 2)
	if byteCount != expectedByteCount {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeWriteMultipleRegisters | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataValue},
		}
	}

	if len(request.Data) < 5+byteCount {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeWriteMultipleRegisters | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataAddress},
		}
	}

	registers := common.NewRegisters(request.Data[5 : 5+byteCount])

	for i, reg := range registers {
		s.store.Write(PointTypeHoldingRegister, address+uint16(i), reg.Value())
	}

	responseData := make([]byte, 4)
	binary.BigEndian.PutUint16(responseData[0:2], address)
	binary.BigEndian.PutUint16(responseData[2:4], quantity)

	return &common.ProtocolDataUnit{
		FunctionCode: common.FuncCodeWriteMultipleRegisters,
		Data:         responseData,
	}
}

// handleReadDeviceIdentification 处理读取设备标识请求 (功能码 0x2B)
func (s *RequestHandler) handleReadDeviceIdentification(request *common.ProtocolDataUnit) *common.ProtocolDataUnit {
	if len(request.Data) < 2 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeReadDeviceIdentification | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataAddress},
		}
	}
	if s.DeviceInfo == nil || s.DeviceInfo.Identification == nil {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeReadDeviceIdentification | 0x80,
			Data:         []byte{common.ExceptionCodeServerDeviceFailure},
		}
	}

	// 解析请求参数
	meiType := request.Data[0]
	readDeviceIDCode := request.Data[1]

	// 检查MEI Type是否为设备标识请求
	if meiType != 0x0E {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeReadDeviceIdentification | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataValue},
		}
	}
	if readDeviceIDCode != 0x01 {
		return &common.ProtocolDataUnit{
			FunctionCode: common.FuncCodeReadDeviceIdentification | 0x80,
			Data:         []byte{common.ExceptionCodeIllegalDataValue},
		}
	}

	// 构建响应数据 MEI类型，设备对象ID，一致性等级，随后更多数据，随后对象ID, 对象数量
	responseData := []byte{0x0E, 0x01, 0x01, 0x00, 0x00, 0x00}

	itemCount := byte(0)
	deviceID := s.DeviceInfo.Identification
	if deviceID.VendorName != "" {
		value := []byte(deviceID.VendorName)
		responseData = append(responseData, 0x00, byte(len(value)))
		responseData = append(responseData, value...)
		itemCount++
	}
	if deviceID.ProductCode != "" {
		value := []byte(deviceID.ProductCode)
		responseData = append(responseData, 0x01, byte(len(value)))
		responseData = append(responseData, value...)
		itemCount++
	}
	if deviceID.ProductName != "" {
		value := []byte(deviceID.ProductName)
		responseData = append(responseData, 0x02, byte(len(value)))
		responseData = append(responseData, value...)
		itemCount++
	}
	if deviceID.VendorUrl != "" {
		value := []byte(deviceID.VendorUrl)
		responseData = append(responseData, 0x03, byte(len(value)))
		responseData = append(responseData, value...)
		itemCount++
	}
	if deviceID.ProductName != "" {
		value := []byte(deviceID.ProductName)
		responseData = append(responseData, 0x04, byte(len(value)))
		responseData = append(responseData, value...)
		itemCount++
	}
	if deviceID.ModelName != "" {
		value := []byte(deviceID.ModelName)
		responseData = append(responseData, 0x05, byte(len(value)))
		responseData = append(responseData, value...)
		itemCount++
	}
	if deviceID.UserApplicationName != "" {
		value := []byte(deviceID.UserApplicationName)
		responseData = append(responseData, 0x06, byte(len(value)))
		responseData = append(responseData, value...)
		itemCount++
	}
	responseData[5] = itemCount

	return &common.ProtocolDataUnit{
		FunctionCode: common.FuncCodeReadDeviceIdentification,
		Data:         responseData,
	}
}
