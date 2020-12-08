package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
	"wsmsgswb/app"
	"wsmsgswb/config"
)

var configFile = flag.String("config", "./config/application.yaml", "配置文件路径（相对、绝对）")

func init() {
	fmt.Println("configFile=" + *configFile)
	flag.Parse()
	conf := config.Init(*configFile)

	// gorm配置
	gormConf := &gorm.Config{}

	//初始化日志
	if file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		logrus.SetOutput(io.MultiWriter(os.Stdout, file))
		if conf.ShowSql {
			gormConf.Logger = logger.New(log.New(file, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel:      logger.Info,
			})
		}
	} else {
		logrus.SetOutput(os.Stdout)
		logrus.Error(err)
	}
}

func main() {
	app.InitIris()
}
