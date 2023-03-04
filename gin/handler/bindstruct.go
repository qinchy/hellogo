package handler

import (
	"github.com/gin-gonic/gin"
	myTypes "github.com/qinchy/hellogo/gin/types"
)

// GetDataB 处理Form绑定嵌套对象的方法
func GetDataB(c *gin.Context) {
	var b myTypes.StructB
	c.Bind(&b)
	c.JSON(200, gin.H{
		"a": b.NestedStruct,
		"b": b.FieldB,
	})
}

// GetDataC 处理嵌套结构体指针的方法
func GetDataC(c *gin.Context) {
	var b myTypes.StructC
	c.Bind(&b)
	c.JSON(200, gin.H{
		"a": b.NestedStructPointer,
		"c": b.FieldC,
	})
}

// GetDataD 处理匿名结构体的方法
func GetDataD(c *gin.Context) {
	var d myTypes.StructD
	c.Bind(&d)
	c.JSON(200, gin.H{
		"x": d.NestedAnonyStruct,
		"d": d.FieldD,
	})
}
