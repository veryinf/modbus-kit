package common

import "fmt"

type RTUMessage struct {
}

// Encode encodes PDU in a RTU frame:
func (m RTUMessage) Encode(slaveId byte, pdu *ProtocolDataUnit) (messageData []byte, err error) {
	frame, err := NewRTUFrame(slaveId, pdu)
	if err != nil {
		return
	}
	messageData = frame.ToBytes()
	return
}

func (m RTUMessage) Decode(messageData []byte) (pdu *ProtocolDataUnit, err error) {
	frame, err := NewRTUFrameFromBytes(messageData)
	if err != nil {
		return
	}
	pdu = frame.PDU
	return
}

func (m RTUMessage) Verify(requestData []byte, responseData []byte) (err error) {
	length := len(responseData)
	// Minimum size (including address, function and CRC)
	if length < rtuMinSize {
		err = fmt.Errorf("modbus: response length '%v' does not meet minimum '%v'", length, rtuMinSize)
		return
	}
	// Slave address must match
	if responseData[0] != requestData[0] {
		err = fmt.Errorf("modbus: response slave id '%v' does not match request '%v'", responseData[0], requestData[0])
		return
	}
	return
}
