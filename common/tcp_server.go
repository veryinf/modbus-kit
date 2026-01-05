package common

import (
	"log/slog"

	"github.com/panjf2000/gnet/v2"
)

type connectionContext struct {
	SlaveId   uint8
	FrameType FrameType
}

type NetServer struct {
	gnet.BuiltinEventEngine
	devices []*ModbusDevice
}

func NewNetServer() *NetServer {
	return &NetServer{
		devices: make([]*ModbusDevice, 0),
	}
}

func (s *NetServer) Enroll(device *ModbusDevice) {
	for _, dev := range s.devices {
		if dev.SlaveId == device.SlaveId && dev.FrameType == device.FrameType {
			panic("device already exists")
		}
	}
	s.devices = append(s.devices, device)
}

func (s *NetServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	slog.Info("connection opened", "remote", c.RemoteAddr())
	return nil, gnet.None
}

func (s *NetServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	slog.Warn("connection closed", "error", err)
	return gnet.None
}

func (s *NetServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	deviceContext := c.Context()
	if deviceContext == nil {
		//自动检测协议类型
		if mbapFrame, err := NewMBAPFrameFromBytes(buf); err == nil {
			deviceContext = &connectionContext{
				SlaveId:   mbapFrame.UnitId,
				FrameType: FrameTypeMBAP,
			}
		}
		if rtuFrame, err := NewRTUFrameFromBytes(buf); err == nil {
			deviceContext = &connectionContext{
				SlaveId:   rtuFrame.SlaveId,
				FrameType: FrameTypeRTU,
			}
		}
		if deviceContext != nil {
			c.SetContext(deviceContext)
		}
	}
	if deviceContext != nil {
		ctx := deviceContext.(*connectionContext)
		for _, device := range s.devices {
			if device.SlaveId == ctx.SlaveId && device.FrameType == ctx.FrameType {
				responseData, err := device.Transport.Send(buf)
				if err != nil {
					slog.Warn("handle request data error", "error", err)
					break
				}
				_, err = c.Write(responseData)
				if err != nil {
					slog.Warn("write response data error", "error", err)
					break
				}
			}
		}
	}
	return gnet.None
}
