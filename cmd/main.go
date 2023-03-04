package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	myHandler "github.com/qinchy/hellogo/handler"
	"github.com/qinchy/hellogo/pkg/proto"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	//"github.com/qinchy/hellogo/pkg/write"
	//"github.com/qinchy/hellogo/pkg/scheduler"
	//"github.com/qinchy/hellogo/pkg/write"
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

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/someJSON", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}

		// 输出 : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.AsciiJSON(http.StatusOK, data)
	})

	r.GET("/moreJSON", func(c *gin.Context) {
		// 你也可以使用一个结构体
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		// 注意 msg.Name 在 JSON 中变成了 "user"
		// 将输出：{"user": "Lena", "Message": "hey", "Number": 123}
		c.JSON(http.StatusOK, msg)
	})

	r.GET("/someXML", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/someYAML", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/someProtoBuf", func(c *gin.Context) {
		reps := []int64{int64(1), int64(2)}
		label := "test"
		// protobuf 的具体定义写在 pkg/proto 文件中。
		data := &proto.Test{
			Label: &label,
			Reps:  reps,
		}
		// 请注意，数据在响应中变为二进制数据
		// 将输出被 proto.Test protobuf 序列化了的数据
		c.ProtoBuf(http.StatusOK, data)
	})

	r.LoadHTMLGlob("templates/**/*")

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	r.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})

	r.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
		})
	})

	r.GET("/JSONP", func(c *gin.Context) {
		data := map[string]interface{}{
			"foo": "bar",
		}

		// /JSONP?callback=x
		// 将输出：x({\"foo\":\"bar\"})
		c.JSONP(http.StatusOK, data)
	})

	r.POST("/loginForm", func(c *gin.Context) {
		// 你可以使用显式绑定声明绑定 multipart form：
		// c.ShouldBindWith(&form, binding.Form)
		// 或者简单地使用 ShouldBind 方法自动绑定：
		var form LoginForm
		// 在这种情况下，将自动选择合适的绑定
		if c.ShouldBind(&form) == nil {
			if form.User == "user" && form.Password == "password" {
				c.JSON(200, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(401, gin.H{"status": "unauthorized"})
			}
		}
	})

	// curl -k -X POST --form "name=qinchy" --form "address=hangzhou" --form "birthday=2013-04-27" --form "id=987fbc97-4bed-5078-9f07-9141ba07c9f3"  "https://localhost/bindForm"
	r.POST("/bindForm", func(c *gin.Context) {
		var person Person
		// 如果是 `GET` 请求，只使用 `Form` 绑定引擎（`query`）。
		// 如果是 `POST` 请求，首先检查 `content-type` 是否为 `JSON` 或 `XML`，然后再使用 `Form`（`form-data`）。
		// 查看更多：https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L88
		if c.ShouldBind(&person) == nil {
			log.Println(person.Name)
			log.Println(person.Address)
			log.Println(person.Birthday)
			c.String(200, "Success")
		} else {
			c.String(400, "Error")
		}
	})

	// curl -k "https://localhost/getb?field_a=hello&field_b=world"
	r.GET("/getb", myHandler.GetDataB)
	// curl -k "https://localhost/getb?field_a=hello&field_c=world"
	r.GET("/getc", myHandler.GetDataC)
	// curl -k "https://localhost/getb?field_x=hello&field_d=world"
	r.GET("/getd", myHandler.GetDataD)

	// 绑定 JSON ({"user": "user", "password": "password"})
	r.POST("/loginJSON", func(c *gin.Context) {
		var json LoginForm
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.User != "user" || json.Password != "password" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// 绑定 XML (
	//	<?xml version="1.0" encoding="UTF-8"?>
	//	<root>
	//		<user>user</user>
	//		<password>password</password>
	//	</root>)
	// curl --location 'https://localhost/loginXML' \
	// --header 'Content-Type: application/xml' \
	// --data '<?xml version="1.0" encoding="UTF-8"?>
	// <root>
	//	 <user>user</user>
	//	 <password>password</password>
	// </root>'
	r.POST("/loginXML", func(c *gin.Context) {
		var xml LoginForm
		if err := c.ShouldBindXML(&xml); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if xml.User != "user" || xml.Password != "password" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	r.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})

	// 提供 unicode 实体
	r.GET("/json", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})

	r.GET("/SecureJSON", func(c *gin.Context) {
		names := []string{"lena", "austin", "foo"}

		// 将输出：while(1);["lena","austin","foo"]
		c.SecureJSON(http.StatusOK, names)
	})

	// curl -k -X POST "https://localhost/post1?id=11&page=1"
	r.POST("/post1", func(c *gin.Context) {
		// 请求示例
		// POST /post?id=1234&page=1 HTTP/1.1
		// Content-Type: application/x-www-form-urlencoded
		//name=manu&message=this_is_great

		// 从url中获取字段
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		// 从表单中获取字段
		name := c.PostForm("name")
		message := c.PostForm("message")

		fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
		c.JSON(200, gin.H{
			"id":      id,
			"page":    page,
			"name":    name,
			"message": message,
		})
	})

	// 映射查询字符串或表单参数
	// curl -k -X POST --location "https://localhost/post2?ids\[a\]=11&ids\[b\]=22" --header "Content-Type: application/x-www-form-urlencoded" -d "names[first]=thinkerou&names[second]=tianou"
	r.POST("/post2", func(c *gin.Context) {
		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")

		log.Printf("ids: %v; names: %v", ids, names)
		c.JSON(200, gin.H{
			"ids":  ids,
			"name": names,
		})
	})

	// 提供字面字符
	r.GET("/purejson", func(c *gin.Context) {
		// 这里的sleep是用来模拟优雅关机的，当请求来后，会处理完这个请求后再关闭服务器
		time.Sleep(10 * time.Second)
		c.PureJSON(200, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})

	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	// curl -k -X POST https://localhost/singleupload  -F "file=@D:\Source_Code\go\src\github.com\qinchy\hellogo\cmd\main.go"   -H "Content-Type: multipart/form-data"
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/singleupload", func(c *gin.Context) {
		// 单文件
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		// 这里的当前目录是整个工程的路径，比如当前是$GOPATH/src/github.com/qinchy/hellogo
		// 感觉这个目录是go.mod中module的路径
		dst := "./" + file.Filename
		log.Printf("dst:%v", dst)
		// 上传文件至指定的完整文件路径
		c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	// curl -k -X POST https://localhost/multiupload  -F "upload[]=@C:\Users\Administrator\AppData\Local\Temp\GoLand\___go_build_github_com_qinchy_hellogo_cmd.exe"   -F "upload[]=@D:\Source_Code\go\bin\hellogo\go_build_github_com_qinchy_hellogo.exe"   -H "Content-Type: multipart/form-data"
	r.POST("/multiupload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			dst := "./" + file.Filename + ".bak"
			// 上传文件至指定目录
			c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	r.GET("/someDataFromReader", func(c *gin.Context) {
		response, err := http.Get("https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png")
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})

	//  =================使用 BasicAuth 中间件==================
	// 路由组使用 gin.BasicAuth() 中间件
	// gin.Accounts 是 map[string]string 的一种快捷方式
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	// /admin/secrets 端点
	// 触发 "localhost:443/admin/secrets
	authorized.GET("/secrets", getting)
	//  =================使用 BasicAuth 中间件==================

	// 任意协议的请求到testting，均调用startPage函数
	r.Any("/testing", startPage)

	// 当在中间件或 handler 中启动新的 Goroutine 时，不能使用原始的上下文，必须使用只读副本。
	r.GET("/long_async", func(c *gin.Context) {
		// 创建在 goroutine 中使用的副本
		cCp := c.Copy()
		go func() {
			// 用 time.Sleep() 模拟一个长任务。
			time.Sleep(5 * time.Second)

			// 请注意您使用的是复制的上下文 "cCp"，这一点很重要
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})

	r.GET("/long_sync", func(c *gin.Context) {
		// 用 time.Sleep() 模拟一个长任务。
		time.Sleep(5 * time.Second)

		// 因为没有使用 goroutine，不需要拷贝上下文
		log.Println("Done! in path " + c.Request.URL.Path)
	})

	r.GET("/:name/:id", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
	})

	// 传统启动服务器
	// r.RunTLS(":443", "./cert/server.pem", "./cert/server.key")

	// 增加优雅停机feature
	srv := &http.Server{
		Addr:    ":443",
		Handler: r,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		},
	}

	// 协程启动服务器
	go func() {
		if err := srv.ListenAndServeTLS("./cert/server.pem", "./cert/server.key"); err != nil && err != http.ErrServerClosed {
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

func getting(c *gin.Context) {
	// 获取用户，它是由 BasicAuth 中间件设置的
	user := c.MustGet(gin.AuthUserKey).(string)
	if secret, ok := secrets[user]; ok {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}

// startPage
func startPage(c *gin.Context) {
	var person Person
	if c.ShouldBindQuery(&person) == nil {
		log.Println("====== Only Bind By Query String ======")
		log.Println(person.Name)
		log.Println(person.Address)
	}
	c.String(200, "Success")
}

type LoginForm struct {
	User     string `form:"user" json:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type Person struct {
	ID       string    `form:"id" uri:"id" binding:"required,uuid"`
	Name     string    `form:"name" uri:"name" binding:"required"`
	Address  string    `form:"address" uri:"address"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}
