package common

import (
	"fmt"
)

const (
	FuncCodeReadCoils                  = 1
	FuncCodeReadDiscreteInputs         = 2
	FuncCodeReadHoldingRegisters       = 3
	FuncCodeReadInputRegisters         = 4
	FuncCodeWriteSingleCoil            = 5
	FuncCodeWriteSingleRegister        = 6
	FuncCodeWriteMultipleCoils         = 15
	FuncCodeWriteMultipleRegisters     = 16
	FuncCodeMaskWriteRegister          = 22
	FuncCodeReadWriteMultipleRegisters = 23
	FuncCodeReadFIFOQueue              = 24
	FuncCodeReadDeviceIdentification   = 43
)

var AvailableReadFunctionCodes = []byte{
	FuncCodeReadCoils,
	FuncCodeReadDiscreteInputs,
	FuncCodeReadHoldingRegisters,
	FuncCodeReadInputRegisters,
}

var AvailableWriteFunctionCodes = []byte{
	FuncCodeWriteSingleCoil,
	FuncCodeWriteSingleRegister,
	FuncCodeWriteMultipleCoils,
	FuncCodeWriteMultipleRegisters,
}

type FrameType string

const (
	FrameTypeMBAP = "MBAP"
	FrameTypeRTU  = "RTU"
)

const (
	ExceptionCodeIllegalFunction                    = 1
	ExceptionCodeIllegalDataAddress                 = 2
	ExceptionCodeIllegalDataValue                   = 3
	ExceptionCodeServerDeviceFailure                = 4
	ExceptionCodeAcknowledge                        = 5
	ExceptionCodeServerDeviceBusy                   = 6
	ExceptionCodeMemoryParityError                  = 8
	ExceptionCodeGatewayPathUnavailable             = 10
	ExceptionCodeGatewayTargetDeviceFailedToRespond = 11
)

type DeviceIdentification struct {
	VendorName          string //厂商名称 0x00
	ProductCode         string //产品编号 0x01
	ProductVersion      string //产品编号 0x02
	VendorUrl           string //厂商网址 0x03
	ProductName         string //产品名称 0x04
	ModelName           string //模式名称 0x05
	UserApplicationName string //用户应用名称 0x06
}

// Error Modbus 错误定义
type Error struct {
	FunctionCode  byte
	ExceptionCode byte
}

// Error 实现error接口
func (e *Error) Error() string {
	var name string
	switch e.ExceptionCode {
	case ExceptionCodeIllegalFunction:
		name = "illegal function"
	case ExceptionCodeIllegalDataAddress:
		name = "illegal data address"
	case ExceptionCodeIllegalDataValue:
		name = "illegal data value"
	case ExceptionCodeServerDeviceFailure:
		name = "server device failure"
	case ExceptionCodeAcknowledge:
		name = "acknowledge"
	case ExceptionCodeServerDeviceBusy:
		name = "server device busy"
	case ExceptionCodeMemoryParityError:
		name = "memory parity error"
	case ExceptionCodeGatewayPathUnavailable:
		name = "gateway path unavailable"
	case ExceptionCodeGatewayTargetDeviceFailedToRespond:
		name = "gateway target device failed to respond"
	default:
		name = "unknown"
	}
	return fmt.Sprintf("modbus: exception '%v' (%s), function '(%v,%v)'", e.ExceptionCode, name, e.FunctionCode-0x80, e.FunctionCode)
}

// Message 消息包解析定义
type Message interface {
	Encode(slaveId byte, pdu *ProtocolDataUnit) (messageData []byte, err error)
	Decode(messageData []byte) (pdu *ProtocolDataUnit, err error)
	Verify(requestData []byte, responseData []byte) (err error)
}

// Transport 通信层定义
type Transport interface {
	Send(requestData []byte) (responseData []byte, err error)
}

type ModbusDevice struct {
	SlaveId   uint8
	FrameType FrameType
	Transport Transport
	Message   Message
}
