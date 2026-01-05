package master

import (
	"github.com/veryinf/modbus-kit/common"
	"net"
)

func NewModbusRTUOverTCPMasterWithAddress(address string) *ModbusMaster {
	message := &common.RTUMessage{}
	tcpClient := common.NewTCPClient(address)
	transport := &RTUOverTCPTransport{
		client: &tcpClient,
	}
	return NewModbusMaster(message, transport)
}

func NewModbusRTUOverTCPMaster(client *common.TCPClient) *ModbusMaster {
	message := &common.RTUMessage{}
	transport := &RTUOverTCPTransport{
		client: client,
	}
	return NewModbusMaster(message, transport)
}

type RTUOverTCPTransport struct {
	client *common.TCPClient
}

// Send 发送数据到服务器，并确保响应长度大于头部长度
func (t *RTUOverTCPTransport) Send(requestData []byte) (responseData []byte, err error) {
	err = t.client.Send(requestData, func(conn net.Conn) error {
		message := &common.RTUFrame{}
		if e := message.ReadFromConn(requestData, conn); e != nil {
			return e
		}
		responseData = message.ToBytes()
		return nil
	})
	return
}
