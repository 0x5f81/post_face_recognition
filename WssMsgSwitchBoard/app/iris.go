package app

import (
	"github.com/go-resty/resty/v2"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"wsmsgswb/config"
)

//初始化iris框架
func InitIris() {
	app := iris.New()
	app.Logger().SetLevel("warn")
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           600,
		AllowedMethods:   []string{iris.MethodGet, iris.MethodPost, iris.MethodOptions, iris.MethodHead, iris.MethodDelete, iris.MethodPut},
		AllowedHeaders:   []string{"*"},
	}))

	app.AllowMethods(iris.MethodOptions)

	//错误输出
	app.OnAnyErrorCode(func(ctx iris.Context) {
		path := ctx.Path()
		var err error
		if strings.Contains(path, "") {
			_, err = ctx.JSON(simple.JsonErrorCode(ctx.GetStatusCode(), "Http error"))
		}
		if err != nil {
			logrus.Error(err)
		}
	})

	app.Any("/", func(i iris.Context) {
		_, _ = i.HTML("<h1></h1>")
	})

	// 前端接口
	mvc.Configure(app.Party("/api"), func(m *mvc.Application) {
		//m.Party("")
	})

	// 管理后台接口
	mvc.Configure(app.Party("/api"), func(m *mvc.Application) {
		//m.Party("")
	})

	//图片等静态资源
	app.Get("/api/img/proxy", func(i iris.Context) {
		url := i.FormValue("url")
		resp, err := resty.New().R().Get(url)
		i.Header("Content-Type", "image/jpg")
		if nil != err {
			logrus.Error(err)
		} else {
			_, _ = i.Write(resp.Body())
		}
	})
	server := &http.Server{Addr: ":" + config.Instance.HttpPort}
	// 绑定服务终止信号
	handleSignal(server)
	// 启动服务框架
	err := app.Run(iris.Server(server), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           false,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		EnableOptimizations:               true,
		TimeFormat:                        "2006-01-02 15:04:05",
		Charset:                           "UTF-8",
	}))
	// 启动失败，退出
	if err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
}

// 服务停止信号处理
func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	//侦听停止信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	//启动关闭线程
	go func() {
		//线程阻塞，等待信号
		s := <-c
		//关闭服务
		logrus.Infof("收到信号[%s]，正在结束服务...", s)
		if err := server.Close(); nil != err {
			logrus.Errorf("服务结束失败: " + err.Error())
		}
		//关闭数据库链接
		simple.CloseDB()
		//关闭进程
		logrus.Infof("已关闭服务。")
		os.Exit(0)
	}()
}
