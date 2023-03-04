package main

import (
	"context"
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"github.com/qinchy/hellogo/gin/globalvar"
	myHandler "github.com/qinchy/hellogo/gin/handler"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Println("开始启动gin服务器...")
	//go read.ReadFile()
	//go write.WriteFile()
	//go scheduler.PrintTimeEveryMinute()

	//  gin相关
	f, _ := os.Create("gin.log")
	// 改写日志到控制台和文件
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	globalvar.Route.GET("/ping", myHandler.Ping)

	globalvar.Route.GET("/somejson", myHandler.SomeJson)

	globalvar.Route.GET("/morejson", myHandler.MoreJson)

	globalvar.Route.GET("/somexml", myHandler.SomeXml)

	globalvar.Route.GET("/someyaml", myHandler.SomeYaml)

	globalvar.Route.GET("/someprotobuf", myHandler.SomeProtoBuf)

	globalvar.Route.LoadHTMLGlob("templates/**/*")

	globalvar.Route.GET("/index", myHandler.Index)

	globalvar.Route.GET("/posts/index", myHandler.PostIndex)

	globalvar.Route.GET("/users/index", myHandler.UsersIndex)

	globalvar.Route.GET("/jsonp", myHandler.JsonP)

	globalvar.Route.POST("/loginform", myHandler.LoginForm)

	// curl -k -X POST --form "name=qinchy" --form "address=hangzhou" --form "birthday=2013-04-27" --form "id=987fbc97-4bed-5078-9f07-9141ba07c9f3"  "https://localhost/bindForm"
	globalvar.Route.POST("/bindform", myHandler.BindForm)

	// curl -k "https://localhost/getb?field_a=hello&field_b=world"
	globalvar.Route.GET("/getb", myHandler.GetDataB)

	// curl -k "https://localhost/getb?field_a=hello&field_c=world"
	globalvar.Route.GET("/getc", myHandler.GetDataC)

	// curl -k "https://localhost/getb?field_x=hello&field_d=world"
	globalvar.Route.GET("/getd", myHandler.GetDataD)

	// 绑定 JSON ({"user": "user", "password": "password"})
	globalvar.Route.POST("/loginJSON", myHandler.LoginJson)

	// 绑定 XML (
	// curl --location 'https://localhost/loginXML' \
	// --header 'Content-Type: application/xml' \
	// --data '<?xml version="1.0" encoding="UTF-8"?>
	// <root>
	//	 <user>user</user>
	//	 <password>password</password>
	// </root>'
	globalvar.Route.POST("/loginxml", myHandler.LoginXml)

	globalvar.Route.POST("/postform", myHandler.PostForm)

	// 提供 unicode 实体
	globalvar.Route.GET("/json", myHandler.Json)

	globalvar.Route.GET("/SecureJson", myHandler.SecureJson)

	// curl -k -X POST "https://localhost/post1?id=11&page=1"
	globalvar.Route.POST("/postformwithquery", myHandler.PostFormWithQuery)

	// 映射查询字符串或表单参数
	// curl -k -X POST --location "https://localhost/post2?ids\[a\]=11&ids\[b\]=22" --header "Content-Type: application/x-www-form-urlencoded" -d "names[first]=thinkerou&names[second]=tianou"
	globalvar.Route.POST("/postmultiformwithquery", myHandler.PostMultiFormWithQuery)

	// 提供字面字符
	globalvar.Route.GET("/purejson", myHandler.PureJson)

	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	// curl -k -X POST https://localhost/singleupload  -F "file=@D:\Source_Code\go\src\github.com\qinchy\hellogo\cmd\main.go"   -H "Content-Type: multipart/form-data"
	globalvar.Route.MaxMultipartMemory = 8 << 20 // 8 MiB
	globalvar.Route.POST("/singleupload", myHandler.SingleUpload)

	// curl -k -X POST https://localhost/multiupload  -F "upload[]=@C:\Users\Administrator\AppData\Local\Temp\GoLand\___go_build_github_com_qinchy_hellogo_cmd.exe"   -F "upload[]=@D:\Source_Code\go\bin\hellogo\go_build_github_com_qinchy_hellogo.exe"   -H "Content-Type: multipart/form-data"
	globalvar.Route.POST("/multiupload", myHandler.MultiUpload)

	globalvar.Route.GET("/fetchfromreader", myHandler.FetchFromReader)

	//  =================使用 BasicAuth 中间件==================
	// 路由组使用 gin.BasicAuth() 中间件
	// gin.Accounts 是 map[string]string 的一种快捷方式
	authorized := globalvar.Route.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	// /admin/secrets 端点
	// 触发 "localhost:443/admin/secrets
	authorized.GET("/secrets", myHandler.Getting)
	//  =================使用 BasicAuth 中间件==================

	// 任意协议的请求到testting，均调用startPage函数
	globalvar.Route.Any("/testing", myHandler.StartPage)

	// 当在中间件或 handler 中启动新的 Goroutine 时，不能使用原始的上下文，必须使用只读副本。
	globalvar.Route.GET("/longasync", myHandler.LongAsync)

	globalvar.Route.GET("/longsync", myHandler.LongSync)

	globalvar.Route.GET("/:name/:id", myHandler.GetDataByUri)

	// 传统启动服务器
	// r.RunTLS(":443", "./cert/server.pem", "./cert/server.key")

	// 增加优雅停机feature
	srv := &http.Server{
		Addr:    ":443",
		Handler: globalvar.Route,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		},
	}

	// 协程启动服务器
	go func() {
		if err := srv.ListenAndServeTLS("./gin/cert/server.pem", "./gin/cert/server.key"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败，错误原因: %s\n", err)
		}
	}()

	log.Println("服务器启动完成")

	// 定义一个关闭服务器接受信号的通道
	quit := make(chan os.Signal)
	// 这个通道只接收os.Interrupt信号
	signal.Notify(quit, os.Interrupt)
	// 如果从通道中接收信号，就调用srv的shutdown优雅的关闭服务器
	<-quit
	log.Println("接收到关闭信号，服务器关闭中...")

	// 定义一个在后台5秒钟关闭的context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭出现异常，错误原因：%s\n", err)
	}
	log.Println("服务器正常停止")
}
