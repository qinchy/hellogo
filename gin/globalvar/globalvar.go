package globalvar

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

var Secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}
var (
	// Route 全局Route
	Route *gin.Engine
)

// init 定制化gin的参数可以放到这里
func init() {
	//  gin相关
	f, _ := os.OpenFile("gin.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0766)
	// 改写日志到控制台和文件
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 返回什么格式,日志格式就是什么样子
	var formatter = func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("客户端IP:%s,请求时间:[%s],请求方式:%s,请求地址:%s,http协议版本:%s,请求状态码:%d,响应时间:%s,客户端:%s，错误信息:%s\n",
			param.ClientIP,
			param.TimeStamp.Format("2006年01月02日 15:03:04"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}
	gin.LoggerWithFormatter(formatter)

	Route = gin.Default()
}
