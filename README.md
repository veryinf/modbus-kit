# Modbus Kit

A high-performance Modbus toolkit written in Go, supporting Modbus TCP and RTU over TCP protocols, with complete Master and Slave functionality.

## Features

### Supported Protocols
- ✅ Modbus TCP
- ✅ RTU over TCP

### Master Functions
- Read Coils (Function Code 01)
- Read Discrete Inputs (Function Code 02)
- Read Holding Registers (Function Code 03)
- Read Input Registers (Function Code 04)
- Write Single Coil (Function Code 05)
- Write Single Register (Function Code 06)
- Write Multiple Coils (Function Code 15)
- Write Multiple Registers (Function Code 16)
- Read Device Identification (Function Code 43)

### Slave Functions
- Respond to all Master-supported function codes
- Memory data storage
- Device identification configuration

### Technical Features
- Built on high-performance network library [gnet](https://github.com/panjf2000/gnet)
- Pure Go implementation, no C dependencies
- Support for concurrent processing
- Complete error handling mechanism
- Compliant with Modbus protocol specifications

## Installation

```bash
go get -u github.com/veryinf/modbus-kit
```

## Quick Start

### Run Examples

The project provides a unified example entry point in the `example/main.go` file. You can run it using:

```bash
go run ./example/main.go
```

This will present a menu where you can select which example to run:
- Modbus TCP Slave
- Modbus TCP Master
- RTU over TCP Slave
- RTU over TCP Master

### Create Modbus TCP Slave

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
    // Configure logging
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    slog.SetDefault(logger)

    // Create device information
    deviceInfo := &slave.DeviceInfo{
        Title: "Example Modbus Device",
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

    // Create memory data store
    store := slave.NewMemoryDataStore()

    // Set initial register values
    for i := 0; i < 10; i++ {
        store.Write(slave.PointTypeHoldingRegister, uint16(i), uint16(i*100))
    }

    // Create Modbus TCP Slave instance
    slaveDevice := slave.NewModbusTCPSlave(1, deviceInfo, store)
    tcpServer := common.NewNetServer()
    tcpServer.Enroll(&slaveDevice.ModbusDevice)
    
    logger.Info("Starting Modbus TCP Slave server on tcp://0.0.0.0:502")
    err := gnet.Run(tcpServer, "tcp://0.0.0.0:502", gnet.WithMulticore(true))
    if err != nil {
        logger.Error("Server error", "error", err)
    }
}
```

### Create Modbus TCP Master

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
    // Configure logging
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    slog.SetDefault(logger)

    // Create TCP client
    client := common.NewTCPClient("localhost:502")
    
    // Create Modbus TCP Master instance
    tcpMaster := master.NewModbusTCPMaster(&client)

    // Read holding registers
    registers, err := tcpMaster.ReadHoldingRegisters(1, 0, 5)
    if err != nil {
        logger.Error("Read holding registers error", "error", err)
        return
    }
    
    for i, reg := range registers {
        logger.Info("Read holding register", "register", i, "value", reg.Value())
    }

    // Write single register
    err = tcpMaster.WriteSingleRegister(1, 0, 12345)
    if err != nil {
        logger.Error("Write single register error", "error", err)
        return
    }
    logger.Info("Write single register success")
    
    // Read again to verify
    registers, err = tcpMaster.ReadHoldingRegisters(1, 0, 1)
    if err != nil {
        logger.Error("Read holding registers error", "error", err)
        return
    }
    logger.Info("Read updated holding register", "value", registers[0].GetValue())
}
```

### Create RTU over TCP Slave

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
    // Configure logging
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    slog.SetDefault(logger)

    // Create device information
    deviceInfo := &slave.DeviceInfo{
        Title: "Example RTU Device",
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

    // Create memory data store
    store := slave.NewMemoryDataStore()

    // Set initial values
    for i := 0; i < 10; i++ {
        store.Write(slave.PointTypeHoldingRegister, uint16(i), uint16(i*100))
    }

    // Create Modbus RTU over TCP Slave instance
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

### Create RTU over TCP Master

```go
package main

import (
    "log/slog"
    "os"

    "github.com/veryinf/modbus-kit/common"
    "github.com/veryinf/modbus-kit/master"
)

func main() {
    // Configure logging
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    slog.SetDefault(logger)

    // Create TCP client
    client := common.NewTCPClient("localhost:502")
    
    // Create RTU over TCP Master instance
    rtuMaster := master.NewModbusRTUOverTCPMaster(&client)

    // Read holding registers
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

## API Documentation

### Master API

#### Create Master Instance

```go
// TCP Master
client := common.NewTCPClient("localhost:502")
tcpMaster := master.NewModbusTCPMaster(&client)

// RTU over TCP Master
rtuMaster := master.NewModbusRTUOverTCPMaster(&client)
```

#### Read Coils (Function Code 01)

```go
bitVector, err := master.ReadCoils(slaveId byte, address uint16, quantity uint16)
```

#### Read Discrete Inputs (Function Code 02)

```go
bitVector, err := master.ReadDiscreteInputs(slaveId byte, address uint16, quantity uint16)
```

#### Read Holding Registers (Function Code 03)

```go
registers, err := master.ReadHoldingRegisters(slaveId byte, address uint16, quantity uint8)
```

#### Read Input Registers (Function Code 04)

```go
registers, err := master.ReadInputRegisters(slaveId byte, address uint16, quantity uint8)
```

#### Write Single Coil (Function Code 05)

```go
err := master.WriteSingleCoil(slaveId byte, address uint16, state bool)
```

#### Write Single Register (Function Code 06)

```go
err := master.WriteSingleRegister(slaveId byte, address uint16, value uint16)
```

#### Write Multiple Coils (Function Code 15)

```go
err := master.WriteMultipleCoils(slaveId byte, address uint16, values []bool)
```

#### Write Multiple Registers (Function Code 16)

```go
err := master.WriteMultipleRegisters(slaveId byte, address uint16, registers []*common.Register)
```

#### Read Device Identification (Function Code 43)

```go
info, err := master.ReadDeviceIdentification(slaveId byte)
```

### Slave API

#### Create Slave Instance

```go
// Create memory data store
store := slave.NewMemoryDataStore()

// Create device information
deviceInfo := &slave.DeviceInfo{
    Title: "My Device",
    Identification: &common.DeviceIdentification{
        VendorName: "My Company",
        ProductCode: "PRO-001",
        ProductVersion: "1.0.0",
    },
}

// Create TCP Slave
slaveDevice := slave.NewModbusTCPSlave(1, deviceInfo, store)
```

#### Start TCP Server

```go
tcpServer := common.NewNetServer()
tcpServer.Enroll(&slaveDevice.ModbusDevice)
err := gnet.Run(tcpServer, "tcp://0.0.0.0:502", gnet.WithMulticore(true))
```

## Project Structure

```
modbus-kit/
├── common/           # Common types and utilities
│   ├── bit_vector.go # Bit vector implementation
│   ├── crc.go        # CRC checksum
│   ├── data_frame.go # Data frame processing
│   ├── mbap_frame.go # MBAP frame processing
│   ├── mbap_message.go # MBAP message processing
│   ├── register.go   # Register implementation
│   ├── rtu_frame.go  # RTU frame processing
│   ├── rtu_message.go # RTU message processing
│   ├── tcp_client.go # TCP client
│   ├── tcp_server.go # TCP server
│   └── types.go      # Type definitions and constants
├── example/          # Example applications
│   ├── main.go           # Main example file
│   ├── tcp_master.go     # Modbus TCP Master example
│   ├── tcp_slave.go      # Modbus TCP Slave example
│   ├── rtu_over_tcp_master.go # RTU over TCP Master example
│   └── rtu_over_tcp_slave.go  # RTU over TCP Slave example
├── master/           # Master functionality
│   ├── modbus_master.go # Core Master implementation
│   ├── rtu_over_tcp.go # RTU over TCP Master
│   └── tcp.go        # TCP Master
├── slave/            # Slave functionality
│   ├── modbus_slave.go # Core Slave implementation
│   ├── request_handler.go # Request handling
│   ├── rtu_over_tcp.go # RTU over TCP Slave
│   ├── store.go      # Data storage
│   └── tcp.go        # TCP Slave
├── README.md         # English README
├── README_CN.md      # Chinese README
├── go.mod            # Go module definition
└── go.sum            # Dependency checksums
```

## Technology Stack

- **Go Version**: 1.25+
- **Network Library**: [gnet/v2](https://github.com/panjf2000/gnet/v2) - High-performance event-driven network framework

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

If you have any questions or suggestions, please contact us through:

- Create an [Issue](https://github.com/yourusername/modbus-kit/issues)

## Acknowledgments

- [gnet](https://github.com/panjf2000/gnet) - High-performance network library
- Modbus Protocol Specification (http://www.modbus.org/specs.php)
