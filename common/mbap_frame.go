package common

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	mbapProtocolIdentifier uint16 = 0x0000
	mbapHeaderSize                = 7
	tcpMaxLength                  = 260
)

type MBAPFrame struct {
	TransactionId uint16
	ProtocolId    uint16
	Length        uint16
	UnitId        byte
	PDU           *ProtocolDataUnit
}

func (f *MBAPFrame) ToBytes() []byte {
	messageData := make([]byte, mbapHeaderSize+1+len(f.PDU.Data))
	binary.BigEndian.PutUint16(messageData, f.TransactionId)
	binary.BigEndian.PutUint16(messageData[2:], f.ProtocolId)
	length := uint16(1 + 1 + len(f.PDU.Data))
	binary.BigEndian.PutUint16(messageData[4:], length)
	messageData[6] = f.UnitId
	copy(messageData[mbapHeaderSize:], f.PDU.ToBytes())
	return messageData
}

func (f *MBAPFrame) ReadFromConn(conn net.Conn) error {
	// 读取 MBAP 头
	var data [tcpMaxLength]byte
	if _, err := io.ReadFull(conn, data[:mbapHeaderSize]); err != nil {
		return err
	}
	f.TransactionId = binary.BigEndian.Uint16(data[:2])
	f.ProtocolId = binary.BigEndian.Uint16(data[2:4])
	f.Length = binary.BigEndian.Uint16(data[4:6])
	f.UnitId = data[6]
	if f.Length <= 0 {
		return fmt.Errorf("modbus: length in response header '%v' must not be zero", f.Length)
	}
	if f.Length > (tcpMaxLength - (mbapHeaderSize - 1)) {
		return fmt.Errorf("modbus: length in response header '%v' must not greater than '%v'", f.Length, tcpMaxLength-mbapHeaderSize+1)
	}
	// 确定结束位置
	endIndex := f.Length + mbapHeaderSize - 1
	if _, err := io.ReadFull(conn, data[mbapHeaderSize:endIndex]); err != nil {
		return err
	}
	f.PDU = &ProtocolDataUnit{
		FunctionCode: data[mbapHeaderSize],
		Data:         data[mbapHeaderSize+1 : endIndex],
	}
	return nil
}

func NewMBAPFrame(transactionId uint16, slaveId byte, pdu *ProtocolDataUnit) *MBAPFrame {
	frame := &MBAPFrame{
		TransactionId: transactionId,
		ProtocolId:    mbapProtocolIdentifier,
		UnitId:        slaveId,
		Length:        uint16(1 + 1 + len(pdu.Data)),
		PDU:           pdu,
	}
	return frame
}

func NewMBAPFrameFromBytes(messageData []byte) (frame *MBAPFrame, err error) {
	transactionId := binary.BigEndian.Uint16(messageData[:2])
	protocolId := binary.BigEndian.Uint16(messageData[2:4])
	length := binary.BigEndian.Uint16(messageData[4:6])
	unitId := messageData[6]
	payloadLength := len(messageData) - mbapHeaderSize
	if payloadLength <= 0 || payloadLength != int(length-1) {
		err = fmt.Errorf("modbus: length in response '%v' does not match pdu data length '%v'", length-1, payloadLength)
		return
	}
	frame = &MBAPFrame{
		TransactionId: transactionId,
		ProtocolId:    protocolId,
		Length:        length,
		UnitId:        unitId,
		PDU: &ProtocolDataUnit{
			FunctionCode: messageData[mbapHeaderSize],
			Data:         messageData[mbapHeaderSize+1:],
		},
	}
	return
}
