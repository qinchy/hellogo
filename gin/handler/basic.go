package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qinchy/hellogo/gin/globalvar"
	"github.com/qinchy/hellogo/gin/types"
	"github.com/qinchy/hellogo/pkg/proto"
	"log"
	"net/http"
	"time"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// StartPage handler
func StartPage(c *gin.Context) {
	var person types.Person
	if c.ShouldBindQuery(&person) == nil {
		log.Println("====== Only Bind By Query String ======")
		log.Println(person.Name)
		log.Println(person.Address)
	}
	c.String(200, "Success")
}

// Getting handler
func Getting(c *gin.Context) {
	// 获取用户，它是由 BasicAuth 中间件设置的
	user := c.MustGet(gin.AuthUserKey).(string)
	if secret, ok := globalvar.Secrets[user]; ok {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}

func SomeJson(c *gin.Context) {
	data := map[string]interface{}{
		"lang": "GO语言",
		"tag":  "<br>",
	}

	// 输出 : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
	c.AsciiJSON(http.StatusOK, data)
}

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

func SomeXml(c *gin.Context) {
	c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
}

func SomeYaml(c *gin.Context) {
	c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
}

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

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Main website",
	})
}

func PostIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
		"title": "Posts",
	})
}

func UsersIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
		"title": "Users",
	})
}

func JsonP(c *gin.Context) {
	data := map[string]interface{}{
		"foo": "bar",
	}

	// /JSONP?callback=x
	// 将输出：x({\"foo\":\"bar\"})
	c.JSONP(http.StatusOK, data)
}

func Json(c *gin.Context) {
	c.JSON(200, gin.H{
		"html": "<b>Hello, world!</b>",
	})
}

func SecureJson(c *gin.Context) {
	names := []string{"lena", "austin", "foo"}

	// 将输出：while(1);["lena","austin","foo"]
	c.SecureJSON(http.StatusOK, names)
}

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

func BindForm(c *gin.Context) {
	var person types.Person
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
}

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

func PostForm(c *gin.Context) {
	message := c.PostForm("message")
	nick := c.DefaultPostForm("nick", "anonymous")

	c.JSON(200, gin.H{
		"status":  "posted",
		"message": message,
		"nick":    nick,
	})
}

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

func PostMultiFormWithQuery(c *gin.Context) {
	ids := c.QueryMap("ids")
	names := c.PostFormMap("names")

	log.Printf("ids: %v; names: %v", ids, names)
	c.JSON(200, gin.H{
		"ids":  ids,
		"name": names,
	})
}

func PureJson(c *gin.Context) {
	// 这里的sleep是用来模拟优雅关机的，当请求来后，会处理完这个请求后再关闭服务器
	time.Sleep(10 * time.Second)
	c.PureJSON(200, gin.H{
		"html": "<b>Hello, world!</b>",
	})
}

func SingleUpload(c *gin.Context) {
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
}

func MultiUpload(c *gin.Context) {
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
}

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

func LongAsync(c *gin.Context) {
	// 创建在 goroutine 中使用的副本
	cCp := c.Copy()
	go func() {
		// 用 time.Sleep() 模拟一个长任务。
		time.Sleep(5 * time.Second)

		// 请注意您使用的是复制的上下文 "cCp"，这一点很重要
		log.Println("Done! in path " + cCp.Request.URL.Path)
	}()
}

func LongSync(c *gin.Context) {
	// 用 time.Sleep() 模拟一个长任务。
	time.Sleep(5 * time.Second)

	// 因为没有使用 goroutine，不需要拷贝上下文
	log.Println("Done! in path " + c.Request.URL.Path)
}

func GetDataByUri(c *gin.Context) {
	var person types.Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
}
