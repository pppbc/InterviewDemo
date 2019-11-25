package control

import (
	"io/ioutil"
	"net/http"
	"pm25_warning/config"
	logs "pm25_warning/log"

	"github.com/thedevsaddam/gojsonq"
	"go.uber.org/zap"
	gomail "gopkg.in/gomail.v2"
)

type Info struct {
	City string  `json:"city"` //城市
	Aqi  float64 `json:"aqi"`  //pm2.5值
}

//通过API获取指定城市pm2.5的实时数据
func GetPM25(info *Info) {
	token := "4ae0b38cd16a395fe0eb598fe3a8a88b2b64b075"

	//初始化路径
	url := "https://api.waqi.info/feed/" + info.City + "/?token=" + token

	//以get方式请求
	resp, err := http.Get(url)
	if err != nil {
		//panic，请求都失败了，不再考虑向下执行
		logs.Loggers.Error("error", zap.String("status", "[gatPM25]请求失败"))
		panic(err)
	}

	//记得关闭
	defer resp.Body.Close()

	//读response body里的信息
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//panic，数据都没有，没有后续
		logs.Loggers.Error("error", zap.String("status", "[gatPM25]读数据失败"))
		panic(err)
	}

	//使用gojsonq包访问嵌套结构的json文件，得到的数据为interface{}类型
	result := gojsonq.New().FromString(string(body)).Find("data.iaqi.pm25.v")

	//数据转float64，赋值给info
	info.Aqi = result.(float64)

	//成功写日志
	logs.Loggers.Info("success", zap.String("status", "[getPM25]获取数据成功"))
}

//发邮件提醒
func WarningByEmail(header config.Header, dialer config.Dialer) {
	//初始化一条新邮件
	m := gomail.NewMessage()

	//设置信息头、内容
	m.SetHeader("From", header.From)
	m.SetHeader("To", header.To)
	m.SetHeader("Subject", header.Subject)
	m.SetBody("text/html", "<p>"+header.Msg+"</p>")

	//初始化拨号连接
	d := gomail.NewDialer(dialer.Host, dialer.Port, dialer.Username, dialer.Password)

	//发送
	if err := d.DialAndSend(m); err != nil {
		//失败写日志
		logs.Loggers.Error("error", zap.String("status", "[sendEmail]发送失败"))
	} else {
		//成功写日志
		logs.Loggers.Info("success", zap.String("status", "[sendEmail]给<"+header.To+">的提醒邮件发送成功"))
	}
}

//判断空气质量
func JudgeAqi(info Info, standard float64) bool {
	//大于标准值则返回true，否则返回false
	if info.Aqi > standard {
		return true
	} else {
		return false
	}
}

//***************************************************************************
//执行
func Run() {
	//获取配置文件里的数据
	var (
		standard = config.GetStandard()
		header   = config.GetHeader()
		dialer   = config.GetDialer()
		info     = &Info{City: config.GetCity()}
	)

	//获取空气质量,传info地址，直接写入pm2.5值
	GetPM25(info)

	//判断所测城市空气质量
	if JudgeAqi(*info, standard) {

		//若返回true，则大于标准，发送邮件
		WarningByEmail(header, dialer)
	} else {

		//否则不发送信息，写入日志
		logs.Loggers.Info("success", zap.String("status", "[success]"+info.City+"的空气状况良好，无需发送提醒"))
	}
}
