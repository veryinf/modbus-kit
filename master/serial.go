package master

import (
	"github.com/goburrow/serial"
	"github.com/veryinf/modbus-kit/common"
)

func NewModbusRTUMaster(config *serial.Config) *ModbusMaster {
	message := &common.RTUMessage{}
	serialClient := common.NewSerialClient(config)
	transport := &RTUTransport{
		client: &serialClient,
	}
	return NewModbusMaster(message, transport)
}

func NewModbusRTUMasterWithClient(client *common.SerialClient) *ModbusMaster {
	message := &common.RTUMessage{}
	transport := &RTUTransport{
		client: client,
	}
	return NewModbusMaster(message, transport)
}

type RTUTransport struct {
	client *common.SerialClient
}

func (t *RTUTransport) Send(requestData []byte) (responseData []byte, err error) {
	err = t.client.Send(requestData, func(port serial.Port) error {
		message := &common.RTUFrame{}
		if e := message.ReadFromConn(requestData, port, common.ReadConfig{
			ConnectionType: common.ConnectionTypeSerial,
			BaudRate:       t.client.Config.BaudRate,
			ReadTimeout:    t.client.Config.Timeout,
		}); e != nil {
			return e
		}
		responseData = message.ToBytes()
		return nil
	})
	return
}
