package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Modbus Kit 示例程序")
	fmt.Println("================")
	fmt.Println("1. 运行 Modbus TCP Slave")
	fmt.Println("2. 运行 Modbus TCP Master")
	fmt.Println("3. 运行 RTU over TCP Slave")
	fmt.Println("4. 运行 RTU over TCP Master")
	fmt.Println("5. 退出")
	fmt.Println("================")

	fmt.Print("请选择要运行的示例: ")
	var choice int
	_, err := fmt.Scanf("%d", &choice)
	if err != nil {
		fmt.Println("输入错误，请输入数字")
		return
	}

	switch choice {
	case 1:
		fmt.Println("运行 Modbus TCP Slave...")
		RunTCPSlave()
	case 2:
		fmt.Println("运行 Modbus TCP Master...")
		RunTCPMaster()
	case 3:
		fmt.Println("运行 RTU over TCP Slave...")
		RunRTUOverTCPSlave()
	case 4:
		fmt.Println("运行 RTU over TCP Master...")
		RunRTUOverTCPMaster()
	case 5:
		fmt.Println("退出程序")
		os.Exit(0)
	default:
		fmt.Println("无效的选择，请重新运行程序")
	}
}