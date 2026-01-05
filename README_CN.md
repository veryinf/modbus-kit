# Modbus Kit

一个用Go语言编写的高性能Modbus工具包，支持Modbus TCP和RTU over TCP协议，提供完整的Master和Slave功能实现。

## 功能特性

### 支持的协议
- ✅ Modbus TCP
- ✅ RTU over TCP

### Master功能
- 读线圈 (Function Code 01)
- 读离散输入 (Function Code 02)
- 读保持寄存器 (Function Code 03)
- 读输入寄存器 (Function Code 04)
- 写单个线圈 (Function Code 05)
- 写单个寄存器 (Function Code 06)
- 写多个线圈 (Function Code 15)
- 写多个寄存器 (Function Code 16)
- 读取设备标识 (Function Code 43)

### Slave功能
- 响应所有Master支持的功能码
- 内存数据存储
- 设备标识信息配置

### 技术特点
- 基于高性能网络库 [gnet](https://github.com/panjf2000/gnet) 实现
- 纯Go语言实现，无C依赖
- 支持并发处理
- 完整的错误处理机制
- 符合Modbus协议规范

## 安装

```bash
go get -u github.com/yourusername/modbus-kit
```

## 快速开始

### 运行示例

项目在 `example/main.go` 文件中提供了统一的示例入口点。您可以使用以下命令运行它：

```bash
go run ./example/main.go
```

这将显示一个菜单，您可以选择运行哪个示例：
- Modbus TCP Slave
- Modbus TCP Master
- RTU over TCP Slave
- RTU over TCP Master

### 创建Modbus TCP Slave

```go
package main

import (
    "log/slog"
    "modbus-kit/common"
    "modbus-kit/slave"
    "github.com/panjf2000/gnet/v2"
)

func main() {
    // 配置日志
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    slog.SetDefault(logger)

    // 创建设备信息
    deviceInfo := &slave.DeviceInfo{
        Title: "第一个设备",
        Identification: &common.DeviceIdentification{
            VendorName:          "ModbusKit Inc.",
            ProductCode:         "MK-001",
            ProductVersion:      "1.0.0",
            VendorUrl:           "https://github.com/veryinf/modbus-kit",
            ProductName:         "Modbus Slave Device",
            ModelName:           "MSD-100",
            UserApplicationName: "Modbus Test Application",
        },
    }

    // 创建内存数据存储
    store := slave.NewMemoryDataStore()

    // 设置初始值
    for i := 0; i < 10; i++ {
        store.Write(slave.PointTypeHoldingRegister, uint16(i), uint16(i*100))
    }

    slaveDevice := slave.NewModbusTCPSlave(1, deviceInfo, store)
    tcpServer := common.NewNetServer()
    tcpServer.Enroll(&slaveDevice.ModbusDevice)
    err := gnet.Run(tcpServer, "tcp://0.0.0.0:502", gnet.WithMulticore(true))
    if err != nil {
        logger.Error("server error", "error", err)
    }
}
```

### 创建Modbus TCP Master

```go
package main

import (
    "log/slog"
    "os"
    "time"

    "github.com/veryinf/modbus-kit/common"
    "github.com/veryinf/modbus-kit/master"
)

func main() {
    // 配置日志
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    slog.SetDefault(logger)

    client := common.NewTCPClient("localhost:502")
    tcpMaster := master.NewModbusTCPMaster(&client)
    
    // 读保持寄存器
    registers, err := tcpMaster.ReadHoldingRegisters(1, 0, 5)
    if err != nil {
        logger.Error("read holding registers error", "error", err)
        return
    }
    
    for i, reg := range registers {
        logger.Info("read holding register", "register", i, "value", reg.Value())
    }
    
    // 写单个寄存器
    err = tcpMaster.WriteSingleRegister(1, 0, 12345)
    if err != nil {
        logger.Error("write single register error", "error", err)
        return
    }
    logger.Info("write single register success")
    
    // 再次读取以验证
    registers, err = tcpMaster.ReadHoldingRegisters(1, 0, 1)
    if err != nil {
        logger.Error("read holding registers error", "error", err)
        return
    }
    logger.Info("read updated holding register", "value", registers[0].Value())
}
```

### 创建RTU over TCP Slave

```go
package main

import (
    "log/slog"
    "os"

    "github.com/veryinf/modbus-kit/common"
    "github.com/veryinf/modbus-kit/slave"
    "github.com/panjf2000/gnet/v2"
)

func main() {
    // 配置日志
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    slog.SetDefault(logger)

    // 创建设备信息
    deviceInfo := &slave.DeviceInfo{
        Title: "RTU设备示例",
        Identification: &common.DeviceIdentification{
            VendorName:          "ModbusKit Inc.",
            ProductCode:         "MK-002",
            ProductVersion:      "1.0.0",
            VendorUrl:           "https://github.com/veryinf/modbus-kit",
            ProductName:         "RTU over TCP Device",
            ModelName:           "RTD-100",
            UserApplicationName: "Modbus RTU Test Application",
        },
    }

    // 创建内存数据存储
    store := slave.NewMemoryDataStore()

    // 设置初始值
    for i := 0; i < 10; i++ {
        store.Write(slave.PointTypeHoldingRegister, uint16(i), uint16(i*100))
    }

    // 创建RTU over TCP Slave实例
    slaveDevice := slave.NewModbusRTUOverTCPSlave(1, deviceInfo, store)
    tcpServer := common.NewNetServer()
    tcpServer.Enroll(&slaveDevice.ModbusDevice)
    
    logger.Info("Starting RTU over TCP Slave server on tcp://0.0.0.0:502")
    err := gnet.Run(tcpServer, "tcp://0.0.0.0:502", gnet.WithMulticore(true))
    if err != nil {
        logger.Error("Server error", "error", err)
    }
}
```

### 创建RTU over TCP Master

```go
package main

import (
    "log/slog"
    "os"

    "github.com/veryinf/modbus-kit/common"
    "github.com/veryinf/modbus-kit/master"
)

func main() {
    // 配置日志
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    slog.SetDefault(logger)

    // 创建TCP客户端
    client := common.NewTCPClient("localhost:502")
    
    // 创建RTU over TCP Master实例
    rtuMaster := master.NewModbusRTUOverTCPMaster(&client)

    // 读保持寄存器
    registers, err := rtuMaster.ReadHoldingRegisters(1, 0, 5)
    if err != nil {
        logger.Error("Read holding registers error", "error", err)
        return
    }
    
    for i, reg := range registers {
        logger.Info("Read holding register", "register", i, "value", reg.Value())
    }
}
```

## API文档

### Master API

#### 创建Master实例

```go
// TCP Master
client := common.NewTCPClient("localhost:502")
tcpMaster := master.NewModbusTCPMaster(&client)

// RTU over TCP Master
rtuMaster := master.NewModbusRTUOverTCPMaster(&client)
```

#### 读线圈 (Function Code 01)

```go
bitVector, err := master.ReadCoils(slaveId byte, address uint16, quantity uint16)
```

#### 读离散输入 (Function Code 02)

```go
bitVector, err := master.ReadDiscreteInputs(slaveId byte, address uint16, quantity uint16)
```

#### 读保持寄存器 (Function Code 03)

```go
registers, err := master.ReadHoldingRegisters(slaveId byte, address uint16, quantity uint8)
```

#### 读输入寄存器 (Function Code 04)

```go
registers, err := master.ReadInputRegisters(slaveId byte, address uint16, quantity uint8)
```

#### 写单个线圈 (Function Code 05)

```go
err := master.WriteSingleCoil(slaveId byte, address uint16, state bool)
```

#### 写单个寄存器 (Function Code 06)

```go
err := master.WriteSingleRegister(slaveId byte, address uint16, value uint16)
```

#### 写多个线圈 (Function Code 15)

```go
err := master.WriteMultipleCoils(slaveId byte, address uint16, values []bool)
```

#### 写多个寄存器 (Function Code 16)

```go
err := master.WriteMultipleRegisters(slaveId byte, address uint16, registers []*common.Register)
```

#### 读取设备标识 (Function Code 43)

```go
info, err := master.ReadDeviceIdentification(slaveId byte)
```

### Slave API

#### 创建Slave实例

```go
// 创建内存数据存储
store := slave.NewMemoryDataStore()

// 创建设备信息
deviceInfo := &slave.DeviceInfo{
    Title: "My Device",
    Identification: &common.DeviceIdentification{
        VendorName: "My Company",
        ProductCode: "PRO-001",
        ProductVersion: "1.0.0",
    },
}

// 创建TCP Slave
slaveDevice := slave.NewModbusTCPSlave(1, deviceInfo, store)

// 创建设备信息
deviceInfo := &slave.DeviceInfo{
    Title: "My Device",
    Identification: &common.DeviceIdentification{
        VendorName: "My Company",
        ProductCode: "PRO-001",
        ProductVersion: "1.0.0",
    },
}

// 创建TCP Slave
slaveDevice := slave.NewModbusTCPSlave(1, deviceInfo, store)
```

#### 启动TCP Server

```go
tcpServer := common.NewNetServer()
tcpServer.Enroll(&slaveDevice.ModbusDevice)
err := gnet.Run(tcpServer, "tcp://0.0.0.0:502", gnet.WithMulticore(true))
```

## 项目结构

```
modbus-kit/
├── common/           # 通用类型和工具
│   ├── bit_vector.go # 位向量实现
│   ├── crc.go        # CRC校验
│   ├── data_frame.go # 数据帧处理
│   ├── mbap_frame.go # MBAP帧处理
│   ├── mbap_message.go # MBAP消息处理
│   ├── register.go   # 寄存器实现
│   ├── rtu_frame.go  # RTU帧处理
│   ├── rtu_message.go # RTU消息处理
│   ├── tcp_client.go # TCP客户端
│   ├── tcp_server.go # TCP服务器
│   └── types.go      # 类型定义和常量
├── example/          # 示例应用
│   ├── main.go           # 主示例入口
│   ├── tcp_master.go     # Modbus TCP Master示例
│   ├── tcp_slave.go      # Modbus TCP Slave示例
│   ├── rtu_over_tcp_master.go # RTU over TCP Master示例
│   └── rtu_over_tcp_slave.go  # RTU over TCP Slave示例
├── master/           # Master功能
│   ├── modbus_master.go # 核心Master实现
│   ├── rtu_over_tcp.go # RTU over TCP Master
│   └── tcp.go        # TCP Master
├── slave/            # Slave功能
│   ├── modbus_slave.go # 核心Slave实现
│   ├── request_handler.go # 请求处理
│   ├── rtu_over_tcp.go # RTU over TCP Slave
│   ├── store.go      # 数据存储
│   └── tcp.go        # TCP Slave
├── README.md         # 英文README
├── README_CN.md      # 中文README
├── go.mod            # Go模块定义
└── go.sum            # 依赖校验
```

## 技术栈

- **Go版本**: 1.25+
- **网络库**: [gnet/v2](https://github.com/panjf2000/gnet/v2) - 高性能事件驱动网络框架

## 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 仓库
2. 创建你的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

如有问题或建议，请通过以下方式联系：

- 创建 [Issue](https://github.com/yourusername/modbus-kit/issues)
- 发送邮件到：your.email@example.com

## 致谢

- [gnet](https://github.com/panjf2000/gnet) - 高性能网络库
- Modbus协议规范 (http://www.modbus.org/specs.php)
