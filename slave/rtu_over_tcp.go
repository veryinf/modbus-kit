package slave

import "github.com/veryinf/modbus-kit/common"

func NewModbusRTUOverTCPSlave(slaveId uint8, deviceInfo *DeviceInfo, store *MemoryDataStore) *ModbusSlave {
	transport := &TCPTransport{}
	transport.RequestHandler.store = store
	transport.RequestHandler.DeviceInfo = deviceInfo
	slaveInfo := common.ModbusDevice{
		SlaveId:   slaveId,
		FrameType: common.FrameTypeRTU,
		Transport: transport,
	}
	return NewModbusSlave(slaveInfo, deviceInfo, store)
}

type RTUOverTCPTransport struct {
	RequestHandler
}

func (t *RTUOverTCPTransport) Send(requestData []byte) (responseData []byte, err error) {
	frame, err := common.NewRTUFrameFromBytes(requestData)
	if err != nil {
		return nil, err
	}
	response, err := t.HandleRequest(frame.PDU)
	if err != nil {
		return nil, err
	}
	frame.PDU = response
	return frame.ToBytes(), nil
}
