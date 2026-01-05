package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/veryinf/modbus-kit/common"
	"github.com/veryinf/modbus-kit/master"
)

func RunTCPMaster() {
	// 配置日志
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	// 创建TCP客户端
	client := common.NewTCPClient("localhost:502")

	// 创建Modbus TCP Master实例
	tcpMaster := master.NewModbusTCPMaster(&client)

	// Slave ID
	slaveID := byte(1)

	logger.Info("开始Modbus TCP Master测试")

	// 1. 读取线圈 (Function Code 01)
	logger.Info("测试: 读取线圈")
	coils, err := tcpMaster.ReadCoils(slaveID, 0, 10)
	if err != nil {
		logger.Error("读取线圈失败", "错误", err)
	} else {
		for i := 0; i < 10; i++ {
			logger.Info(fmt.Sprintf("线圈 %d: %t", i, coils.Get(uint(i))))
		}
	}

	time.Sleep(500 * time.Millisecond)

	// 2. 读取离散输入 (Function Code 02)
	logger.Info("测试: 读取离散输入")
	discreteInputs, err := tcpMaster.ReadDiscreteInputs(slaveID, 0, 10)
	if err != nil {
		logger.Error("读取离散输入失败", "错误", err)
	} else {
		for i := 0; i < 10; i++ {
			logger.Info(fmt.Sprintf("离散输入 %d: %t", i, discreteInputs.Get(uint(i))))
		}
	}

	time.Sleep(500 * time.Millisecond)

	// 3. 读取保持寄存器 (Function Code 03)
	logger.Info("测试: 读取保持寄存器")
	holdingRegisters, err := tcpMaster.ReadHoldingRegisters(slaveID, 0, 5)
	if err != nil {
		logger.Error("读取保持寄存器失败", "错误", err)
	} else {
		for i, reg := range holdingRegisters {
			logger.Info(fmt.Sprintf("保持寄存器 %d: %d (0x%s)", i, reg.Value(), reg.ToHexString()))
		}
	}

	time.Sleep(500 * time.Millisecond)

	// 4. 读取输入寄存器 (Function Code 04)
	logger.Info("测试: 读取输入寄存器")
	inputRegisters, err := tcpMaster.ReadInputRegisters(slaveID, 0, 5)
	if err != nil {
		logger.Error("读取输入寄存器失败", "错误", err)
	} else {
		for i, reg := range inputRegisters {
			logger.Info(fmt.Sprintf("输入寄存器 %d: %d (0x%s)", i, reg.Value(), reg.ToHexString()))
		}
	}

	time.Sleep(500 * time.Millisecond)

	// 5. 写入单个线圈 (Function Code 05)
	logger.Info("测试: 写入单个线圈")
	err = tcpMaster.WriteSingleCoil(slaveID, 0, true)
	if err != nil {
		logger.Error("写入单个线圈失败", "错误", err)
	} else {
		logger.Info("写入单个线圈成功")
	}

	time.Sleep(500 * time.Millisecond)

	// 6. 写入单个寄存器 (Function Code 06)
	logger.Info("测试: 写入单个寄存器")
	err = tcpMaster.WriteSingleRegister(slaveID, 0, 12345)
	if err != nil {
		logger.Error("写入单个寄存器失败", "错误", err)
	} else {
		logger.Info("写入单个寄存器成功")
	}

	time.Sleep(500 * time.Millisecond)

	// 7. 写入多个线圈 (Function Code 15)
	logger.Info("测试: 写入多个线圈")
	coilValues := []bool{true, false, true, false, true, false, true, false}
	err = tcpMaster.WriteMultipleCoils(slaveID, 0, coilValues)
	if err != nil {
		logger.Error("写入多个线圈失败", "错误", err)
	} else {
		logger.Info("写入多个线圈成功")
	}

	time.Sleep(500 * time.Millisecond)

	// 8. 写入多个寄存器 (Function Code 16)
	logger.Info("测试: 写入多个寄存器")
	registerValues := []*common.Register{
		common.NewRegisterFromUInt16(100),
		common.NewRegisterFromUInt16(200),
		common.NewRegisterFromUInt16(300),
	}
	err = tcpMaster.WriteMultipleRegisters(slaveID, 1, registerValues)
	if err != nil {
		logger.Error("写入多个寄存器失败", "错误", err)
	} else {
		logger.Info("写入多个寄存器成功")
	}

	time.Sleep(500 * time.Millisecond)

	// 9. 读取设备标识 (Function Code 43)
	logger.Info("测试: 读取设备标识")
	deviceInfo, err := tcpMaster.ReadDeviceIdentification(slaveID)
	if err != nil {
		logger.Error("读取设备标识失败", "错误", err)
	} else {
		logger.Info(fmt.Sprintf("供应商名称: %s", deviceInfo.VendorName))
		logger.Info(fmt.Sprintf("产品代码: %s", deviceInfo.ProductCode))
		logger.Info(fmt.Sprintf("产品版本: %s", deviceInfo.ProductVersion))
		logger.Info(fmt.Sprintf("产品名称: %s", deviceInfo.ProductName))
		logger.Info(fmt.Sprintf("型号名称: %s", deviceInfo.ModelName))
	}

	logger.Info("Modbus TCP Master测试完成")
}
