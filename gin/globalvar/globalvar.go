package globalvar

import (
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}
var (
	// Route 全局Route
	Route *gin.Engine

	//Logger 全局Logger
	Logger *logrus.Logger
)

// init 定制化gin的参数可以放到这里
func init() {
	//  gin相关
	gin.DisableConsoleColor() // 禁止控制台日志颜色

	gin.SetMode(gin.ReleaseMode)

	Route = gin.Default()
	Logger = logrus.New()

	Route.Use(loggerToFile())

}

// LoggerToFile 日志记录到文件
func loggerToFile() gin.HandlerFunc {

	//写入文件
	logFile, err := os.OpenFile("./gin.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		panic("系统初始化日志时出现错误：" + err.Error())
	}

	//设置输出
	Logger.Out = logFile

	//设置日志级别
	Logger.SetLevel(logrus.DebugLevel)

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		// 这里在windows上非首次启动会报错，linux上没问题。
		// 具体原因是windows上用WithLinkName创建软链接和linux上实现不一样，windows会把gin.log修改成只读，导致后续启动时logFile是空指针。
		// windows上可考虑把所有日志文件删除，或者修改gin.log为非只读。
		(*logFile).Name()+".%Y%m%d",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName((*logFile).Name()),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 新增 Hook
	Logger.AddHook(lfHook)

	//设置日志格式
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		Logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
		}).Debug()
	}
}
