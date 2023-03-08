package main

import (
	"context"
	"crypto/tls"
	. "github.com/qinchy/hellogo/gin/globalvar"
	"github.com/qinchy/hellogo/gin/handler"
	"github.com/qinchy/hellogo/pkg/read"
	"github.com/qinchy/hellogo/pkg/scheduler"
	"github.com/qinchy/hellogo/pkg/write"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	Logger.Info("开始初始化任务引擎...")
	schedule()
	Logger.Info("任务引擎初始化完成")

	Logger.Info("开始启动gin服务器...")

	// 所有请求路径及函数都放这里便于管理
	handler.Handler()

	// 传统启动服务器
	// r.RunTLS(":443", "./gin/cert/server.pem", "./gin/cert/server.key")

	// 增加优雅停机feature
	srv := &http.Server{
		Addr:    ":443",
		Handler: Route,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		},
	}

	// 协程启动服务器
	go func() {
		if err := srv.ListenAndServeTLS("./gin/cert/server.pem", "./gin/cert/server.key"); err != nil && err != http.ErrServerClosed {
			Logger.Fatalf("服务器启动失败，错误原因: %s\n", err)
		}
	}()

	Logger.Info("服务器启动完成")

	// 定义一个关闭服务器接受信号的通道
	quit := make(chan os.Signal)
	// 这个通道只接收os.Interrupt信号
	signal.Notify(quit, os.Interrupt)
	// 如果从通道中接收信号，就调用srv的shutdown优雅的关闭服务器
	<-quit
	Logger.Info("接收到关闭信号，服务器关闭中...")

	// 定义一个在后台5秒钟关闭的context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		Logger.Fatalf("服务器关闭出现异常，错误原因：%s\n", err)
	}
	Logger.Info("服务器正常停止")
}

func schedule() {
	go read.ReadFile()
	go write.WriteFile()
	go scheduler.PrintTimeEveryMinute()
}
