package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	. "github.com/qinchy/hellogo/gin/globalvar"
	"github.com/qinchy/hellogo/gin/proto"
	"github.com/qinchy/hellogo/gin/types"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

// Ping 访问/ping的处理器
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// StartPage 从Query中获取数据并组装到结构体的处理器
func StartPage(c *gin.Context) {
	var person types.Person
	if c.ShouldBindQuery(&person) == nil {
		Logger.WithFields(logrus.Fields{
			"Name":     person.Name,
			"Address":  person.Address,
			"Birthday": person.Birthday,
		}).Info("请求参数")
	}
	c.String(200, "Success")
}

// Getting 基础认证通过后获取数据的处理器
func Getting(c *gin.Context) {
	// 获取用户，它是由 BasicAuth 中间件设置的
	user := c.MustGet(gin.AuthUserKey).(string)
	if secret, ok := Secrets[user]; ok {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}

// SomeJson 返回JSON的处理器
func SomeJson(c *gin.Context) {
	data := map[string]interface{}{
		"lang": "GO语言",
		"tag":  "<br>",
	}

	// 输出 : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
	c.AsciiJSON(http.StatusOK, data)
}

// MoreJson 返回结构体的JSON的处理器
func MoreJson(c *gin.Context) {
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
}

// SomeXml 返回XML的处理器
func SomeXml(c *gin.Context) {
	c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
}

// SomeYaml 返回yaml的处理器
func SomeYaml(c *gin.Context) {
	c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
}

// SomeProtoBuf 返回ProtoBuf的处理器
func SomeProtoBuf(c *gin.Context) {
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
}

// Index 加载根目录前端模板
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Main website",
	})
}

// PostIndex 加载post目录前端模板
func PostIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
		"title": "Posts",
	})
}

// UsersIndex 加载users目录前端模板
func UsersIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
		"title": "Users",
	})
}

// JsonP 返回JSONP数据
func JsonP(c *gin.Context) {
	data := map[string]interface{}{
		"foo": "bar",
	}

	// /JSONP?callback=x
	// 将输出：x({\"foo\":\"bar\"})
	c.JSONP(http.StatusOK, data)
}

// Json 返回JSON数据
func Json(c *gin.Context) {
	c.JSON(200, gin.H{
		"html": "<b>Hello, world!</b>",
	})
}

// SecureJson 返回SecureJson数据
func SecureJson(c *gin.Context) {
	names := []string{"lena", "austin", "foo"}

	// 将输出：while(1);["lena","austin","foo"]
	c.SecureJSON(http.StatusOK, names)
}

// LoginForm 请求体绑定的结构体处理器
func LoginForm(c *gin.Context) {
	// 你可以使用显式绑定声明绑定 multipart form：
	// c.ShouldBindWith(&form, binding.Form)
	// 或者简单地使用 ShouldBind 方法自动绑定：
	var form types.LoginForm
	// 在这种情况下，将自动选择合适的绑定
	if c.ShouldBind(&form) == nil {
		if form.User == "user" && form.Password == "password" {
			c.JSON(200, gin.H{"status": "you are logged in"})
		} else {
			c.JSON(401, gin.H{"status": "unauthorized"})
		}
	}
}

// BindForm 请求体绑定的结构体处理器
func BindForm(c *gin.Context) {
	var person types.Person
	// 如果是 `GET` 请求，只使用 `Form` 绑定引擎（`query`）。
	// 如果是 `POST` 请求，首先检查 `content-type` 是否为 `JSON` 或 `XML`，然后再使用 `Form`（`form-data`）。
	// 查看更多：https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L88
	if c.ShouldBind(&person) == nil {
		Logger.WithFields(logrus.Fields{
			"Name":     person.Name,
			"Address":  person.Address,
			"Birthday": person.Birthday,
		}).Info("请求参数")
		c.String(200, "Success")
		return
	}

	c.String(400, "Error")
}

// LoginJson json绑定的结构体
func LoginJson(c *gin.Context) {
	var json types.LoginForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if json.User != "user" || json.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}

// LoginXml xml绑定到结构体
func LoginXml(c *gin.Context) {
	var xml types.LoginForm
	if err := c.ShouldBindXML(&xml); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if xml.User != "user" || xml.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}

// PostForm 从表单中获取数据
func PostForm(c *gin.Context) {
	message := c.PostForm("message")
	nick := c.DefaultPostForm("nick", "anonymous")

	c.JSON(200, gin.H{
		"status":  "posted",
		"message": message,
		"nick":    nick,
	})
}

// PostFormWithQuery 从表单和请求串中获取数据
func PostFormWithQuery(c *gin.Context) {
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
}

// PostMultiFormWithQuery 从请求传和表单中获取多个同名key的数据
func PostMultiFormWithQuery(c *gin.Context) {
	ids := c.QueryMap("ids")
	names := c.PostFormMap("names")

	log.Printf("ids: %v; names: %v", ids, names)
	Logger.WithFields(logrus.Fields{
		"ids":   ids,
		"names": names,
	}).Info("请求参数")
	c.JSON(200, gin.H{
		"ids":  ids,
		"name": names,
	})
}

// PureJson 返回纯粹的JSON
func PureJson(c *gin.Context) {
	// 这里的sleep是用来模拟优雅关机的，当请求来后，会处理完这个请求后再关闭服务器
	time.Sleep(10 * time.Second)
	c.PureJSON(200, gin.H{
		"html": "<b>Hello, world!</b>",
	})
}

// SingleUpload 通过表单上传单个文件
func SingleUpload(c *gin.Context) {
	// 单文件
	file, _ := c.FormFile("file")

	// 这里的当前目录是整个工程的路径，比如当前是$GOPATH/src/github.com/qinchy/hellogo
	// 感觉这个目录是go.mod中module的路径
	dst := "./" + file.Filename
	// 上传文件至指定的完整文件路径
	err := c.SaveUploadedFile(file, dst)
	if err != nil {
		Logger.Fatalf("上传文件时出现异常：%s", err.Error())
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

// MultiUpload 通过表单上传多个文件
func MultiUpload(c *gin.Context) {
	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]

	for _, file := range files {
		dst := "./" + file.Filename + ".bak"
		// 上传文件至指定目录
		err := c.SaveUploadedFile(file, dst)
		if err != nil {
			Logger.Fatalf("上传文件时出现异常：%s", err.Error())
		}
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

// FetchFromReader 从reader中获取数据
func FetchFromReader(c *gin.Context) {
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
}

// LongAsync 在处理器中使用协程，需要使用上下文的副本
func LongAsync(c *gin.Context) {
	// 创建在 goroutine 中使用的副本
	cCp := c.Copy()
	go func() {
		// 用 time.Sleep() 模拟一个长任务。
		time.Sleep(5 * time.Second)

		// 请注意您使用的是复制的上下文 "cCp"，这一点很重要
		Logger.WithFields(logrus.Fields{
			"Path": cCp.Request.URL.Path,
		}).Info("Done! in path ")
	}()
}

// LongSync 同步方法使用原始上下文
func LongSync(c *gin.Context) {
	// 用 time.Sleep() 模拟一个长任务。
	time.Sleep(5 * time.Second)

	// 因为没有使用 goroutine，不需要拷贝上下文
	Logger.WithFields(logrus.Fields{
		"Path": c.Request.URL.Path,
	}).Info("Done! in path ")
}

// GetDataByUri 通过uri访问资源
func GetDataByUri(c *gin.Context) {
	var person types.Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
}

// GetBookable 是否可以订购，校验
func GetBookable(c *gin.Context) {
	var b types.Booking
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
