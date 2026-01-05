package master

import (
	"encoding/binary"
	"fmt"
	"github.com/veryinf/modbus-kit/common"
)

type ModbusMaster struct {
	message   common.Message
	transport common.Transport
}

// NewModbusMaster 创建一个新的 ModbusMaster 对象
func NewModbusMaster(message common.Message, transport common.Transport) *ModbusMaster {
	return &ModbusMaster{message: message, transport: transport}
}

// ReadCoils
// Request:
//
//	Function code         : 1 byte (0x01)
//	Starting address      : 2 bytes
//	Quantity of coils     : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x01)
//	Byte count            : 1 byte
//	Coil status           : N* bytes (=N or N+1)
func (c *ModbusMaster) ReadCoils(slaveId byte, address uint16, quantity uint16) (bitVector *common.BitVector, err error) {
	if quantity < 1 || quantity > 2000 {
		err = fmt.Errorf("modbus: quantity '%v' is out of range [1, 2000]", quantity)
		return
	}
	request := &common.ProtocolDataUnit{FunctionCode: common.FuncCodeReadCoils}
	request.LoadData(address, quantity)
	response, err := c.send(slaveId, request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	length := len(response.Data) - 1
	if count != length {
		err = fmt.Errorf("modbus: response data size '%v' does not match count '%v'", length, count)
		return
	}
	bitVector = common.NewBitVector(uint(quantity))
	bitVector.Load(response.Data[1:])
	return
}

// ReadDiscreteInputs
// Request:
//
//	Function code         : 1 byte (0x02)
//	Starting address      : 2 bytes
//	Quantity of inputs    : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x02)
//	Byte count            : 1 byte
//	Input status          : N* bytes (=N or N+1)
func (c *ModbusMaster) ReadDiscreteInputs(slaveId byte, address uint16, quantity uint16) (bitVector *common.BitVector, err error) {
	if quantity < 1 || quantity > 2000 {
		err = fmt.Errorf("modbus: quantity '%v' is out of range [1, 2000]", quantity)
		return
	}
	request := &common.ProtocolDataUnit{FunctionCode: common.FuncCodeReadDiscreteInputs}
	request.LoadData(address, quantity)
	response, err := c.send(slaveId, request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	length := len(response.Data) - 1
	if count != length {
		err = fmt.Errorf("modbus: response data size '%v' does not match count '%v'", length, count)
		return
	}
	bitVector = common.NewBitVector(uint(quantity))
	bitVector.Load(response.Data[1:])
	return
}

// ReadHoldingRegisters
// Request:
//
//	Function code         : 1 byte (0x03)
//	Starting address      : 2 bytes
//	Quantity of registers : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x03)
//	Byte count            : 1 byte
//	Register value        : Nx2 bytes
func (c *ModbusMaster) ReadHoldingRegisters(slaveId byte, address uint16, quantity uint8) (registers []*common.Register, err error) {
	if quantity < 1 || quantity > 125 {
		err = fmt.Errorf("modbus: quantity '%v' is out of range [1, 125]", quantity)
		return
	}
	request := &common.ProtocolDataUnit{FunctionCode: common.FuncCodeReadHoldingRegisters}
	request.LoadData(address, uint16(quantity))
	response, err := c.send(slaveId, request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	length := len(response.Data) - 1
	if count != length {
		err = fmt.Errorf("modbus: response data size '%v' does not match count '%v'", length, count)
		return
	}
	registers = common.NewRegisters(response.Data[1:])
	return
}

// ReadInputRegisters
// Request:
//
//	Function code         : 1 byte (0x04)
//	Starting address      : 2 bytes
//	Quantity of registers : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x04)
//	Byte count            : 1 byte
//	Input registers       : N bytes
func (c *ModbusMaster) ReadInputRegisters(slaveId byte, address uint16, quantity uint8) (registers []*common.Register, err error) {
	if quantity < 1 || quantity > 125 {
		err = fmt.Errorf("modbus: quantity '%v' is out of range [1, 125]", quantity)
		return
	}
	request := &common.ProtocolDataUnit{FunctionCode: common.FuncCodeReadInputRegisters}
	request.LoadData(address, uint16(quantity))
	response, err := c.send(slaveId, request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	length := len(response.Data) - 1
	if count != length {
		err = fmt.Errorf("modbus: response data size '%v' does not match count '%v'", length, count)
		return
	}
	registers = common.NewRegisters(response.Data[1:])
	return
}

// WriteSingleCoil
// Request:
//
//	Function code         : 1 byte (0x05)
//	Output address        : 2 bytes
//	Output value          : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x05)
//	Output address        : 2 bytes
//	Output value          : 2 bytes
func (c *ModbusMaster) WriteSingleCoil(slaveId byte, address uint16, state bool) (err error) {
	// The requested ON/OFF state can only be 0xFF00 and 0x0000
	value := uint16(0x0000)
	if state {
		value = 0xFF00
	}
	request := &common.ProtocolDataUnit{FunctionCode: common.FuncCodeWriteSingleCoil}
	request.LoadData(address, value)
	response, err := c.send(slaveId, request)
	if err != nil {
		return
	}
	// Fixed response length
	if len(response.Data) != 4 {
		err = fmt.Errorf("modbus: response data size '%v' does not match expected '%v'", len(response.Data), 4)
		return
	}
	respValue := binary.BigEndian.Uint16(response.Data)
	if address != respValue {
		err = fmt.Errorf("modbus: response address '%v' does not match request '%v'", respValue, address)
		return
	}
	respValue = binary.BigEndian.Uint16(response.Data[2:])
	if value != respValue {
		err = fmt.Errorf("modbus: response value '%v' does not match request '%v'", respValue, value)
		return
	}
	return
}

// WriteSingleRegister
// Request:
//
//	Function code         : 1 byte (0x06)
//	Register address      : 2 bytes
//	Register value        : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x06)
//	Register address      : 2 bytes
//	Register value        : 2 bytes
func (c *ModbusMaster) WriteSingleRegister(slaveId byte, address, value uint16) (err error) {
	request := &common.ProtocolDataUnit{FunctionCode: common.FuncCodeWriteSingleRegister}
	request.LoadData(address, value)
	response, err := c.send(slaveId, request)
	if err != nil {
		return
	}
	// Fixed response length
	if len(response.Data) != 4 {
		err = fmt.Errorf("modbus: response data size '%v' does not match expected '%v'", len(response.Data), 4)
		return
	}
	respValue := binary.BigEndian.Uint16(response.Data)
	if address != respValue {
		err = fmt.Errorf("modbus: response address '%v' does not match request '%v'", respValue, address)
		return
	}
	respValue = binary.BigEndian.Uint16(response.Data[2:])
	if value != respValue {
		err = fmt.Errorf("modbus: response value '%v' does not match request '%v'", respValue, value)
		return
	}
	return
}

// WriteMultipleCoils
// Request:
//
//	Function code         : 1 byte (0x0F)
//	Starting address      : 2 bytes
//	Quantity of outputs   : 2 bytes
//	Byte count            : 1 byte
//	Outputs value         : N* bytes
//
// Response:
//
//	Function code         : 1 byte (0x0F)
//	Starting address      : 2 bytes
//	Quantity of outputs   : 2 bytes
func (c *ModbusMaster) WriteMultipleCoils(slaveId byte, address uint16, values []bool) (err error) {
	quantity := uint16(len(values))
	if quantity < 1 || quantity > 1968 {
		err = fmt.Errorf("modbus: quantity '%v' is out of range [1, 1968]", quantity)
		return
	}
	bv := common.NewBitVectorFromBooleans(values)
	request := &common.ProtocolDataUnit{FunctionCode: common.FuncCodeWriteMultipleCoils}
	outputValues := bv.ToBytes()
	request.LoadData(address, quantity).Append(byte(len(outputValues))).Append(outputValues...)
	response, err := c.send(slaveId, request)
	if err != nil {
		return
	}
	// Fixed response length
	if len(response.Data) != 4 {
		err = fmt.Errorf("modbus: response data size '%v' does not match expected '%v'", len(response.Data), 4)
		return
	}
	respValue := binary.BigEndian.Uint16(response.Data)
	if address != respValue {
		err = fmt.Errorf("modbus: response address '%v' does not match request '%v'", respValue, address)
		return
	}
	respValue = binary.BigEndian.Uint16(response.Data[2:])
	if quantity != respValue {
		err = fmt.Errorf("modbus: response quantity '%v' does not match request '%v'", respValue, quantity)
		return
	}
	return
}

// WriteMultipleRegisters
// Request:
//
//	Function code         : 1 byte (0x10)
//	Starting address      : 2 bytes
//	Quantity of outputs   : 2 bytes
//	Byte count            : 1 byte
//	Registers value       : N* bytes
//
// Response:
//
//	Function code         : 1 byte (0x10)
//	Starting address      : 2 bytes
//	Quantity of registers : 2 bytes
func (c *ModbusMaster) WriteMultipleRegisters(slaveId byte, address uint16, registers []*common.Register) (err error) {
	quantity := uint8(len(registers))
	if quantity < 1 || quantity > 123 {
		err = fmt.Errorf("modbus: quantity '%v' is out of range [1, 123]", quantity)
		return
	}
	value := *common.RegistersToBytes(registers)
	request := &common.ProtocolDataUnit{FunctionCode: common.FuncCodeWriteMultipleRegisters}
	request.LoadData(address, uint16(quantity)).Append(byte(len(value))).Append(value...)
	response, err := c.send(slaveId, request)
	if err != nil {
		return
	}
	// Fixed response length
	if len(response.Data) != 4 {
		err = fmt.Errorf("modbus: response data size '%v' does not match expected '%v'", len(response.Data), 4)
		return
	}
	respValue := binary.BigEndian.Uint16(response.Data)
	if address != respValue {
		err = fmt.Errorf("modbus: response address '%v' does not match request '%v'", respValue, address)
		return
	}
	respValue = binary.BigEndian.Uint16(response.Data[2:])
	if uint16(quantity) != respValue {
		err = fmt.Errorf("modbus: response quantity '%v' does not match request '%v'", respValue, quantity)
		return
	}
	return
}

// ReadDeviceIdentification
// Request:
//
//		Function code         : 1 byte (0x2B)
//		Entity Identifier     : 1 byte (0x0E)
//		Subfunction code      : 1 byte (0x01 0x02)
//	 Object Identifier     : 1 byte
//
// Response:
//
//	Function code         : 1 byte (0x04)
//	Byte count            : 1 byte
//	Input registers       : N bytes
func (c *ModbusMaster) ReadDeviceIdentification(slaveId byte) (info *common.DeviceIdentification, err error) {
	request := &common.ProtocolDataUnit{
		FunctionCode: common.FuncCodeReadDeviceIdentification,
		Data:         []byte{0x0E, 0x01, 0x00},
	}
	response, err := c.send(slaveId, request)
	if err != nil {
		return
	}
	if len(response.Data) < 4 {
		err = fmt.Errorf("modbus: response data size '%v' is less than expected '%v'", len(response.Data), 4)
		return
	}
	if response.Data[2] != 0x01 || response.Data[3] != 0x00 || response.Data[4] != 0x00 {
		err = fmt.Errorf("modbus: currently, read device identification only support read 1 request")
		return
	}
	info = &common.DeviceIdentification{}
	offset := 6
	for offset < len(response.Data) {
		objectIdentifier := response.Data[offset]
		length := response.Data[offset+1]
		text := response.Data[offset+2 : offset+2+int(length)]
		valText := string(text)
		offset += int(length) + 2
		switch objectIdentifier {
		case 0x00:
			info.VendorName = valText
		case 0x01:
			info.ProductCode = valText
		case 0x02:
			info.ProductVersion = valText
		case 0x03:
			info.VendorUrl = valText
		case 0x04:
			info.ProductName = valText
		case 0x05:
			info.ModelName = valText
		case 0x06:
			info.UserApplicationName = valText
		}
	}

	return
}

// 发送请求并检查可能的异常
func (c *ModbusMaster) send(slaveId byte, request *common.ProtocolDataUnit) (response *common.ProtocolDataUnit, err error) {
	requestData, err := c.message.Encode(slaveId, request)
	if err != nil {
		return
	}
	responseData, err := c.transport.Send(requestData)
	if err != nil {
		return
	}
	if err = c.message.Verify(requestData, responseData); err != nil {
		return
	}
	response, err = c.message.Decode(responseData)
	if err != nil {
		return
	}
	// Check correct function code returned (exception)
	if response.FunctionCode != request.FunctionCode {
		err = responseError(response)
		return
	}
	if response.Data == nil || len(response.Data) == 0 {
		// Empty response
		err = fmt.Errorf("modbus: response data is empty")
		return
	}
	return
}

func responseError(response *common.ProtocolDataUnit) error {
	mbError := &common.Error{FunctionCode: response.FunctionCode}
	if response.Data != nil && len(response.Data) > 0 {
		mbError.ExceptionCode = response.Data[0]
	}
	return mbError
}
