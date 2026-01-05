package main

import (
	"log/slog"
	"os"

	"github.com/veryinf/modbus-kit/common"
	"github.com/veryinf/modbus-kit/slave"
	"github.com/panjf2000/gnet/v2"
)

func RunRTUOverTCPSlave() {
	// 配置日志
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	logger.Info("正在启动RTU over TCP Slave服务...")

	// 1. 创建设备信息
	deviceInfo := &slave.DeviceInfo{
		Title: "示例RTU设备",
		Identification: &common.DeviceIdentification{
			VendorName:          "ModbusKit Inc.",
			ProductCode:         "MK-002",
			ProductVersion:      "1.0.0",
			VendorUrl:           "https://github.com/veryinf/modbus-kit",
			ProductName:         "RTU over TCP设备",
			ModelName:           "RTD-100",
			UserApplicationName: "Modbus RTU测试应用",
		},
	}

	// 2. 创建内存数据存储
	store := slave.NewMemoryDataStore()

	// 3. 设置初始数据
	// 设置保持寄存器初始值
	for i := 0; i < 10; i++ {
		store.Write(slave.PointTypeHoldingRegister, uint16(i), uint16(i*200))
	}

	// 设置输入寄存器初始值
	for i := 0; i < 10; i++ {
		store.Write(slave.PointTypeInputRegister, uint16(i), uint16(i*75))
	}

	// 设置线圈初始值（前5个为true，后5个为false）
	for i := 0; i < 10; i++ {
		value := uint16(0)
		if i < 5 {
			value = uint16(1)
		}
		store.Write(slave.PointTypeCoil, uint16(i), value)
	}

	// 设置离散输入初始值（每隔一个设置为true）
	for i := 0; i < 10; i++ {
		value := uint16(0)
		if i%2 != 0 {
			value = uint16(1)
		}
		store.Write(slave.PointTypeDiscreteInput, uint16(i), value)
	}

	// 4. 创建Modbus RTU over TCP Slave实例
	slaveDevice := slave.NewModbusRTUOverTCPSlave(1, deviceInfo, store)

	// 5. 创建TCP服务器
	tcpServer := common.NewNetServer()

	// 6. 注册Slave设备
	tcpServer.Enroll(&slaveDevice.ModbusDevice)

	logger.Info("RTU over TCP Slave已配置完成，正在启动服务器...")
	logger.Info("服务器地址: tcp://0.0.0.0:502")
	logger.Info("Slave ID: 1")
	logger.Info("按Ctrl+C停止服务器")

	// 7. 启动服务器
	err := gnet.Run(tcpServer, "tcp://0.0.0.0:502", gnet.WithMulticore(true))
	if err != nil {
		logger.Error("服务器启动失败", "错误", err)
		os.Exit(1)
	}
}
