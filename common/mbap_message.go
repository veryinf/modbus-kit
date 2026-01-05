package common

import (
	"encoding/binary"
	"fmt"
	"sync/atomic"
)

// MBAPMessage Modbus-MBAP消息定义，实现 Message 接口
type MBAPMessage struct {
	transactionId uint32
}

// Encode 编码数据帧为 Modbus-MBAP 消息格式
func (m *MBAPMessage) Encode(slaveId byte, pdu *ProtocolDataUnit) (messageData []byte, err error) {
	transactionId := atomic.AddUint32(&m.transactionId, 1)
	tcpFrame := NewMBAPFrame(uint16(transactionId), slaveId, pdu)
	return tcpFrame.ToBytes(), nil
}

// Verify 校验请求和相应数据，确认（传输ID，协议ID，和UnitId）
func (m *MBAPMessage) Verify(requestData []byte, responseData []byte) (err error) {
	// Transaction id
	responseVal := binary.BigEndian.Uint16(responseData[:2])
	requestVal := binary.BigEndian.Uint16(requestData[:2])
	if responseVal != requestVal {
		err = fmt.Errorf("modbus: response transaction id '%v' does not match requestData '%v'", responseVal, requestVal)
		return
	}
	// Protocol id
	responseVal = binary.BigEndian.Uint16(responseData[2:4])
	requestVal = binary.BigEndian.Uint16(requestData[2:4])
	if responseVal != requestVal {
		err = fmt.Errorf("modbus: response protocol id '%v' does not match requestData '%v'", responseVal, requestVal)
		return
	}
	// Unit id (1 byte)
	if responseData[6] != requestData[6] {
		err = fmt.Errorf("modbus: response unit id '%v' does not match requestData '%v'", responseData[6], requestData[6])
		return
	}
	return
}

// Decode 解包MBAP Message为 Modbus 数据帧
func (m *MBAPMessage) Decode(messageData []byte) (pdu *ProtocolDataUnit, err error) {
	frame, err := NewMBAPFrameFromBytes(messageData)
	if err != nil {
		return
	}
	pdu = frame.PDU
	return
}
