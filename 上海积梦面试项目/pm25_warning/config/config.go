package config

import (
	logs "pm25_warning/log"
	"strconv"

	"github.com/Unknwon/goconfig"
	"go.uber.org/zap"
)

var Config *goconfig.ConfigFile

type Header struct {
	//邮件信息头、内容
	From    string
	To      string
	Subject string
	Msg     string
}
type Dialer struct {
	//拨号配置信息
	Host     string
	Port     int
	Username string
	Password string
}

//初始化配置文件
func InitConfig() {
	con, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		logs.Loggers.Error("error", zap.String("status", "[config]读配置文件失败"))
	}
	logs.Loggers.Info("success", zap.String("status", "[config]读配置文件成功"))
	Config = con
}

//获取邮件头信息
func GetHeader() Header {
	sec, _ := Config.GetSection("header")

	mySetting := Header{}
	mySetting.From = sec["from"]
	mySetting.To = sec["to"]
	mySetting.Subject = sec["subject"]
	mySetting.Msg = sec["msg"]

	return mySetting
}

//获取拨号信息
func GetDialer() Dialer {
	sec, _ := Config.GetSection("dialer")

	mySetting := Dialer{}
	mySetting.Host = sec["host"]
	mySetting.Port, _ = strconv.Atoi(sec["port"])
	mySetting.Username = sec["username"]
	mySetting.Password = sec["password"]

	return mySetting
}

//获取城市
func GetCity() string {
	value, err := Config.GetValue("base", "city")

	if err != nil {
		logs.Loggers.Error("error", zap.String("status", "[config]获取"+value+"失败"))
	}

	return value
}

//获取比较标准
func GetStandard() float64 {
	value, err := Config.GetValue("base", "standard")

	if err != nil {
		logs.Loggers.Error("error", zap.String("status", "[config]获取"+value+"失败"))
	}
	result, _ := strconv.ParseFloat(value, 64)

	return result
}
