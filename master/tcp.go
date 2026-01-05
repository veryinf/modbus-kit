// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package master

import (
	"github.com/veryinf/modbus-kit/common"
	"net"
)

// NewModbusTCPMasterWithAddress 使用默认处理程序创建 TcpClient
func NewModbusTCPMasterWithAddress(address string) *ModbusMaster {
	message := &common.MBAPMessage{}
	tcpClient := common.NewTCPClient(address)
	transport := &TCPTransport{
		client: &tcpClient,
	}
	return NewModbusMaster(message, transport)
}

func NewModbusTCPMaster(client *common.TCPClient) *ModbusMaster {
	message := &common.MBAPMessage{}
	transport := &TCPTransport{
		client: client,
	}
	return NewModbusMaster(message, transport)
}

// TCPTransport Modbus-TCP 传输定义,实现 Transport 接口.
type TCPTransport struct {
	client *common.TCPClient
}

// Send 发送数据到服务器，并确保响应长度大于头部长度
func (t *TCPTransport) Send(requestData []byte) (responseData []byte, err error) {
	err = t.client.Send(requestData, func(conn net.Conn) error {
		frame := &common.MBAPFrame{}
		if e := frame.ReadFromConn(conn); e != nil {
			return e
		}
		responseData = frame.ToBytes()
		return nil
	})
	if err != nil {
		return
	}
	return
}
