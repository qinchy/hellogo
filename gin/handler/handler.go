package handler

import (
	"github.com/gin-gonic/gin"
	. "github.com/qinchy/hellogo/gin/globalvar"
)

// Handler 所有handler的集合都放这里
func Handler() {
	Route.GET("/ping", Ping)

	Route.GET("/somejson", SomeJson)

	Route.GET("/morejson", MoreJson)

	Route.GET("/somexml", SomeXml)

	Route.GET("/someyaml", SomeYaml)

	Route.GET("/someprotobuf", SomeProtoBuf)

	Route.LoadHTMLGlob("templates/**/*")

	Route.GET("/index", Index)

	Route.GET("/posts/index", PostIndex)

	Route.GET("/users/index", UsersIndex)

	Route.GET("/jsonp", JsonP)

	Route.POST("/loginform", LoginForm)

	// curl -k -X POST --form "name=qinchy" --form "address=hangzhou" --form "birthday=2013-04-27" --form "id=987fbc97-4bed-5078-9f07-9141ba07c9f3"  "https://localhost/bindform"
	Route.POST("/bindform", BindForm)

	// curl -k "https://localhost/getb?field_a=hello&field_b=world"
	Route.GET("/getb", GetDataB)

	// curl -k "https://localhost/getb?field_a=hello&field_c=world"
	Route.GET("/getc", GetDataC)

	// curl -k "https://localhost/getb?field_x=hello&field_d=world"
	Route.GET("/getd", GetDataD)

	// 绑定 JSON ({"user": "user", "password": "password"})
	Route.POST("/loginjson", LoginJson)

	// 绑定 XML (
	// curl --location 'https://localhost/loginXML' \
	// --header 'Content-Type: application/xml' \
	// --data '<?xml version="1.0" encoding="UTF-8"?>
	// <root>
	//	 <user>user</user>
	//	 <password>password</password>
	// </root>'
	Route.POST("/loginxml", LoginXml)

	Route.POST("/postform", PostForm)

	// 提供 unicode 实体
	Route.GET("/json", Json)

	Route.GET("/SecureJson", SecureJson)

	// curl -k -X POST "https://localhost/postformwithquery?id=11&page=1"
	Route.POST("/postformwithquery", PostFormWithQuery)

	// 映射查询字符串或表单参数
	// curl -k -X POST --location "https://localhost/postmultiformwithquery?ids\[a\]=11&ids\[b\]=22" --header "Content-Type: application/x-www-form-urlencoded" -d "names[first]=thinkerou&names[second]=tianou"
	Route.POST("/postmultiformwithquery", PostMultiFormWithQuery)

	// 提供字面字符
	Route.GET("/purejson", PureJson)

	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	// curl -k -X POST https://localhost/singleupload  -F "file=@D:\Source_Code\go\src\github.com\qinchy\hellogo\cmd\main.go"   -H "Content-Type: multipart/form-data"
	Route.MaxMultipartMemory = 8 << 20 // 8 MiB
	Route.POST("/singleupload", SingleUpload)

	// curl -k -X POST https://localhost/multiupload  -F "upload[]=@C:\Users\Administrator\AppData\Local\Temp\GoLand\___go_build_github_com_qinchy_hellogo_cmd.exe"   -F "upload[]=@D:\Source_Code\go\bin\hellogo\go_build_github_com_qinchy_hellogo.exe"   -H "Content-Type: multipart/form-data"
	Route.POST("/multiupload", MultiUpload)

	Route.GET("/fetchfromreader", FetchFromReader)

	//  =================使用 BasicAuth 中间件==================
	// 路由组使用 gin.BasicAuth() 中间件
	// gin.Accounts 是 map[string]string 的一种快捷方式
	authorized := Route.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	// /admin/secrets 端点
	// 触发 "localhost:443/admin/secrets
	authorized.GET("/secrets", Getting)
	//  =================使用 BasicAuth 中间件==================

	// 任意协议的请求到testting，均调用startPage函数
	Route.Any("/testing", StartPage)

	// 当在中间件或 handler 中启动新的 Goroutine 时，不能使用原始的上下文，必须使用只读副本。
	Route.GET("/longasync", LongAsync)

	Route.GET("/longsync", LongSync)

	Route.GET("/:name/:id", GetDataByUri)

	// curl -k "https://localhost/bookable?check_in=2023-04-16&check_out=2023-04-17"
	// curl -k "https://localhost/bookable?check_in=2023-03-08&check_out=2023-03-09"
	Route.GET("/bookable", GetBookable)
}
