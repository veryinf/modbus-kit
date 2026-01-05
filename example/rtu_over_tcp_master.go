package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/veryinf/modbus-kit/common"
	"github.com/veryinf/modbus-kit/master"
)

func RunRTUOverTCPMaster() {
	// 配置日志
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	// 创建TCP客户端
	client := common.NewTCPClient("localhost:502")

	// 创建RTU over TCP Master实例
	rtuMaster := master.NewModbusRTUOverTCPMaster(&client)

	// Slave ID
	slaveID := byte(1)

	logger.Info("开始RTU over TCP Master测试")

	// 1. 读取线圈 (Function Code 01)
	logger.Info("测试: 读取线圈")
	coils, err := rtuMaster.ReadCoils(slaveID, 0, 10)
	if err != nil {
		logger.Error("读取线圈失败", "错误", err)
	} else {
		for i := 0; i < 10; i++ {
			logger.Info(fmt.Sprintf("线圈 %d: %t", i, coils.Get(uint(i))))
		}
	}

	time.Sleep(500 * time.Millisecond)

	// 2. 读取保持寄存器 (Function Code 03)
	logger.Info("测试: 读取保持寄存器")
	holdingRegisters, err := rtuMaster.ReadHoldingRegisters(slaveID, 0, 5)
	if err != nil {
		logger.Error("读取保持寄存器失败", "错误", err)
	} else {
		for i, reg := range holdingRegisters {
			logger.Info(fmt.Sprintf("保持寄存器 %d: %d (0x%s)", i, reg.Value(), reg.ToHexString()))
		}
	}

	time.Sleep(500 * time.Millisecond)

	// 3. 写入单个线圈 (Function Code 05)
	logger.Info("测试: 写入单个线圈")
	err = rtuMaster.WriteSingleCoil(slaveID, 0, true)
	if err != nil {
		logger.Error("写入单个线圈失败", "错误", err)
	} else {
		logger.Info("写入单个线圈成功")
	}

	time.Sleep(500 * time.Millisecond)

	// 4. 写入单个寄存器 (Function Code 06)
	logger.Info("测试: 写入单个寄存器")
	err = rtuMaster.WriteSingleRegister(slaveID, 0, 9999)
	if err != nil {
		logger.Error("写入单个寄存器失败", "错误", err)
	} else {
		logger.Info("写入单个寄存器成功")
	}

	time.Sleep(500 * time.Millisecond)

	// 5. 读取设备标识 (Function Code 43)
	logger.Info("测试: 读取设备标识")
	deviceInfo, err := rtuMaster.ReadDeviceIdentification(slaveID)
	if err != nil {
		logger.Error("读取设备标识失败", "错误", err)
	} else {
		logger.Info(fmt.Sprintf("供应商名称: %s", deviceInfo.VendorName))
		logger.Info(fmt.Sprintf("产品代码: %s", deviceInfo.ProductCode))
		logger.Info(fmt.Sprintf("产品版本: %s", deviceInfo.ProductVersion))
	}

	logger.Info("RTU over TCP Master测试完成")
}
