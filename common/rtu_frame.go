package common

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/goburrow/serial"
)

const (
	rtuMaxSize       = 256
	rtuMinSize       = 4
	rtuExceptionSize = 5
)

type ConnectionType int

const (
	ConnectionTypeSerial ConnectionType = iota
	ConnectionTypeTCP
	ConnectionTypeUDP
)

type RTUFrame struct {
	SlaveId byte
	PDU     *ProtocolDataUnit
	CRC     *CRC
}

// ReadConfig 读取配置
type ReadConfig struct {
	ConnectionType ConnectionType
	BaudRate       int           // 仅串口有效，用于计算字符间超时
	ReadTimeout    time.Duration //仅串口有效，用于控制读取超时
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

func (f *RTUFrame) ToBytes() []byte {
	length := len(f.PDU.Data) + 4
	messageData := make([]byte, length)

	messageData[0] = f.SlaveId
	copy(messageData[1:], f.PDU.ToBytes())

	crc := CRC{}
	crc.Reset().PushBytes(messageData[0 : length-2])
	copy(messageData[length-2:], crc.SumBytes())
	return messageData
}

// ReadFromConn 使用配置从连接读取Modbus RTU响应
func (f *RTUFrame) ReadFromConn(requestData []byte, conn io.Reader, config ReadConfig) error {
	var frameData []byte
	var err error

	switch config.ConnectionType {
	case ConnectionTypeSerial:
		conn, ok := conn.(serial.Port)
		if !ok {
			return fmt.Errorf("modbus: connection is not a serial port")
		}
		frameData, err = f.readFromSerial(requestData, conn, config)
	case ConnectionTypeTCP:
		conn, ok := conn.(*net.TCPConn)
		if !ok {
			return fmt.Errorf("modbus: connection is not a TCP port")
		}
		frameData, err = f.readFromStream(requestData, conn)
	case ConnectionTypeUDP:
		conn, ok := conn.(*net.UDPConn)
		if !ok {
			return fmt.Errorf("modbus: connection is not a UDP port")
		}
		frameData, err = f.readFromDatagram(conn)
	default:
		return fmt.Errorf("modbus: unsupported connection type")
	}

	if err != nil {
		return err
	}

	return f.parseAndValidate(frameData, requestData)
}

func (f *RTUFrame) readFromSerial(requestData []byte, conn serial.Port, config ReadConfig) ([]byte, error) {
	expectedLength := calculateResponseLength(requestData)
	buf := make([]byte, rtuMaxSize)
	pos := 0

	// 阶段1:等待第一个字节
	deadline := time.Now().Add(config.ReadTimeout)
	for time.Now().Before(deadline) {
		n, err := conn.Read(buf[pos : pos+1])
		if err != nil {
			if err == io.EOF {
				continue
			}
			return nil, fmt.Errorf("modbus serial: read error: %w", err)
		}
		if n > 0 {
			pos++
			break
		}
	}

	if pos == 0 {
		return nil, fmt.Errorf("modbus serial: timeout waiting for response")
	}

	// 阶段2:持续读取,使用字符间超时检测帧结束
	lastByteTime := time.Now()
	charTime := time.Duration(11*1000000/config.BaudRate) * time.Microsecond
	charTimeout := charTime * 2

	for pos < expectedLength && pos < rtuMaxSize {
		n, err := conn.Read(buf[pos : pos+1])
		if err != nil {
			// 超时或其他错误
			if pos < rtuMinSize {
				return nil, fmt.Errorf("modbus serial: incomplete frame: %w", err)
			}
			// 已读到最小帧长度,认为帧结束
			break
		}
		if n > 0 {
			pos++
			lastByteTime = time.Now() // 更新最后收到字节的时间
		}
		// 检查是否超过字符间超时
		if time.Since(lastByteTime) > charTimeout {
			break
		}
	}

	if pos < rtuMinSize {
		return nil, fmt.Errorf("modbus serial: frame too short (%d bytes)", pos)
	}

	result := make([]byte, pos)
	copy(result, buf[:pos])
	return result, nil
}

func (f *RTUFrame) readFromStream(requestData []byte, conn *net.TCPConn) ([]byte, error) {
	expectedLength := calculateResponseLength(requestData)
	buf := make([]byte, rtuMaxSize)
	pos := 0

	for pos < expectedLength && pos < rtuMaxSize {
		// 精确计算还需要读多少字节,避免读到多余数据
		bytesNeeded := expectedLength - pos
		if bytesNeeded > rtuMaxSize-pos {
			bytesNeeded = rtuMaxSize - pos
		}

		n, err := conn.Read(buf[pos : pos+bytesNeeded])
		if err != nil {
			if err == io.EOF && pos >= rtuMinSize {
				break
			}
			return nil, fmt.Errorf("modbus TCP: read error at byte %d: %w", pos, err)
		}

		if n == 0 {
			break
		}

		pos += n

		if pos >= expectedLength {
			break
		}
	}

	if pos < rtuMinSize {
		return nil, fmt.Errorf("modbus TCP: incomplete frame (%d bytes)", pos)
	}

	result := make([]byte, pos)
	copy(result, buf[:pos])
	return result, nil
}

func (f *RTUFrame) readFromDatagram(conn *net.UDPConn) ([]byte, error) {
	buf := make([]byte, rtuMaxSize)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil, fmt.Errorf("modbus UDP: read error: %w", err)
	}

	if n < rtuMinSize {
		return nil, fmt.Errorf("modbus UDP: frame too short (%d bytes)", n)
	}

	result := make([]byte, n)
	copy(result, buf[:n])
	return result, nil
}

func (f *RTUFrame) parseAndValidate(frameData []byte, requestData []byte) error {
	frame, err := NewRTUFrameFromBytes(frameData)
	if err != nil {
		return err
	}

	if requestData != nil && len(requestData) > 0 {
		if frame.SlaveId != requestData[0] {
			return fmt.Errorf("modbus: slave id mismatch (got %d, expected %d)",
				frame.SlaveId, requestData[0])
		}

		expectedFunc := requestData[1]
		actualFunc := frame.PDU.FunctionCode

		if actualFunc != expectedFunc && actualFunc != expectedFunc+0x80 {
			return fmt.Errorf("modbus: unexpected function code %d", actualFunc)
		}
	}

	f.SlaveId = frame.SlaveId
	f.PDU = frame.PDU
	f.CRC = frame.CRC
	return nil
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
		length += (1 + count*2)
	case FuncCodeWriteSingleCoil,
		FuncCodeWriteMultipleCoils,
		FuncCodeWriteSingleRegister,
		FuncCodeWriteMultipleRegisters:
		length += 4
	case FuncCodeMaskWriteRegister:
		length += 6
	case FuncCodeReadFIFOQueue:
	default:
	}
	return length
}
