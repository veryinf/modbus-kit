package slave

import (
	"github.com/veryinf/modbus-kit/common"
)

type DeviceInfo struct {
	Title          string
	Identification *common.DeviceIdentification
}

type ModbusSlave struct {
	common.ModbusDevice
	DeviceInfo *DeviceInfo
	Store      *MemoryDataStore
}

// NewModbusSlave 创建一个新的 ModbusSlave 对象
func NewModbusSlave(slaveInfo common.ModbusDevice, deviceInfo *DeviceInfo, store *MemoryDataStore) *ModbusSlave {
	slave := ModbusSlave{
		DeviceInfo: deviceInfo, Store: store,
	}
	slave.ModbusDevice = slaveInfo
	return &slave
}
