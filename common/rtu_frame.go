package common

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	rtuMaxSize       = 256
	rtuMinSize       = 4
	rtuExceptionSize = 5
)

type RTUFrame struct {
	SlaveId byte
	PDU     *ProtocolDataUnit
	CRC     *CRC
}

func (f *RTUFrame) ToBytes() []byte {
	length := len(f.PDU.Data) + 4
	messageData := make([]byte, length)

	messageData[0] = f.SlaveId
	copy(messageData[1:], f.PDU.ToBytes())

	// Append crc
	crc := CRC{}
	crc.Reset().PushBytes(messageData[0 : length-2])
	copy(messageData[length-2:], crc.SumBytes())
	return messageData
}

func (f *RTUFrame) ReadFromConn(requestData []byte, conn net.Conn) error {
	bytesToRead := calculateResponseLength(requestData)
	if bytesToRead > tcpMaxLength {
		return fmt.Errorf("modbus: response length '%v' must not greater than '%v'", bytesToRead, tcpMaxLength)
	}
	delay := calculateDelay(0, len(requestData)+bytesToRead)
	time.Sleep(delay)

	var data [tcpMaxLength]byte
	if _, err := io.ReadFull(conn, data[:2]); err != nil {
		return err
	}
	f.SlaveId = data[0]
	f.PDU = &ProtocolDataUnit{
		FunctionCode: data[1],
	}
	if requestData[0] != f.SlaveId {
		return fmt.Errorf("modbus: response slave id '%v' does not match request '%v'", data[0], requestData[0])
	}
	//正确返回
	if f.PDU.FunctionCode == requestData[1] {
		if _, err := io.ReadFull(conn, data[2:bytesToRead]); err != nil {
			return err
		}
		crc := CRC{}
		crc.Reset().PushBytes(data[0 : bytesToRead-2])
		if !crc.Match(data[bytesToRead-2 : bytesToRead]) {
			checksum := uint16(data[bytesToRead-1])<<8 | uint16(data[bytesToRead-2])
			return fmt.Errorf("modbus: response crc '%v' does not match expected '%v'", checksum, crc.Value())
		}
		f.PDU.Data = data[2 : bytesToRead-2]
		f.CRC = &crc
	} else if f.PDU.FunctionCode == requestData[1]&0x80 {
		//返回异常
		if _, err := io.ReadFull(conn, data[2:rtuExceptionSize]); err != nil {
			return err
		}
		f.PDU.Data = data[2:rtuExceptionSize]
	} else {
		return fmt.Errorf("modbus: response function '%v' does not match request '%v'", data[1], requestData[1])
	}

	return nil
}

func NewRTUFrame(slaveId byte, pdu *ProtocolDataUnit) (frame *RTUFrame, err error) {
	length := len(pdu.Data) + 4
	if length > rtuMaxSize {
		err = fmt.Errorf("modbus: length of data '%v' must not be bigger than '%v'", length, rtuMaxSize)
		return
	}
	frame = &RTUFrame{
		SlaveId: slaveId,
		PDU:     pdu,
	}
	return
}
func NewRTUFrameFromBytes(messageData []byte) (frame *RTUFrame, err error) {
	length := len(messageData)
	//calculated crc
	crc := CRC{}
	crc.Reset().PushBytes(messageData[0 : length-2])
	if !crc.Match(messageData[length-2:]) {
		checksum := uint16(messageData[length-1])<<8 | uint16(messageData[length-2])
		err = fmt.Errorf("modbus: response crc '%v' does not match expected '%v'", checksum, crc.Value())
		return
	}
	frame = &RTUFrame{
		SlaveId: messageData[0],
		PDU: &ProtocolDataUnit{
			FunctionCode: messageData[1],
			Data:         messageData[2 : length-2],
		},
		CRC: &crc,
	}
	return
}

func calculateDelay(baudRate, chars int) time.Duration {
	var characterDelay, frameDelay int
	//无效波特率
	if baudRate <= 0 || baudRate > 19200 {
		characterDelay = 750
		frameDelay = 1750
	} else {
		characterDelay = 15000000 / baudRate
		frameDelay = 35000000 / baudRate
	}
	return time.Duration(characterDelay*chars+frameDelay) * time.Microsecond
}

func calculateResponseLength(requestData []byte) int {
	length := rtuMinSize
	switch requestData[1] {
	case FuncCodeReadDiscreteInputs,
		FuncCodeReadCoils:
		count := int(binary.BigEndian.Uint16(requestData[4:]))
		length += 1 + count/8
		if count%8 != 0 {
			length++
		}
	case FuncCodeReadInputRegisters,
		FuncCodeReadHoldingRegisters,
		FuncCodeReadWriteMultipleRegisters:
		count := int(binary.BigEndian.Uint16(requestData[4:]))
		length += 1 + count*2
	case FuncCodeWriteSingleCoil,
		FuncCodeWriteMultipleCoils,
		FuncCodeWriteSingleRegister,
		FuncCodeWriteMultipleRegisters:
		length += 4
	case FuncCodeMaskWriteRegister:
		length += 6
	case FuncCodeReadFIFOQueue:
		// undetermined
	default:
	}
	return length
}
