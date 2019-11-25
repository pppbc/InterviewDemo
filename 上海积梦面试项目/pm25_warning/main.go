package main

import (
	"pm25_warning/config"
	"pm25_warning/control"
	_ "time"
)

func main() {
	//初始化配置文件
	config.InitConfig()

	//执行
	control.Run()

}
