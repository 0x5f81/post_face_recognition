package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Instance *Config

type Config struct {
	//yaml基本配置
	Env        string `yaml:"Env"`        //打包环境
	HttpPort   string `yaml:"HttpPort"`   //http端口
	WssPort    string `yaml:"WssPort"`    //websocket ssl 端口
	LogFile    string `yaml:"LogFile"`    //日志文件
	ShowSql    bool   `yaml:"ShowSql"`    //是否打印sql
	StaticPath string `yaml:"StaticPath"` //根路径下的静态页面等资源文件目录（相对、绝对路径都可以）
	MysqlUrl   string `yaml:"MysqlUrl"`   //数据库连接

	//微信登录配置
	WeChartLogin struct {
		AppID     string `yaml:"AppID"`     //appid
		AppSecret string `yaml:"AppSecret"` //AppSecret
	}

	//smtp配置
	Smtp struct {
		Host     string `yaml:"Host"`     //smtp 服务器访问地址
		Port     string `yaml:"Port"`     //端口号
		Username string `yaml:"Username"` //用户名
		Password string `yaml:"Password"` //密码
		SSL      bool   `yaml:"SSL"`      //是否使用ssl协议
	}
}

func Init(filename string) *Config {
	Instance = &Config{}
	if yamlFile, err := ioutil.ReadFile(filename); err != nil {
		logrus.Error(err)
	} else if err = yaml.Unmarshal(yamlFile, Instance); err != nil {
		logrus.Error(err)
	}
	return Instance
}
