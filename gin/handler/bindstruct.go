package handler

import (
	"github.com/gin-gonic/gin"
	myTypes "github.com/qinchy/hellogo/gin/types"
)

func GetDataB(c *gin.Context) {
	var b myTypes.StructB
	c.Bind(&b)
	c.JSON(200, gin.H{
		"a": b.NestedStruct,
		"b": b.FieldB,
	})
}

func GetDataC(c *gin.Context) {
	var b myTypes.StructC
	c.Bind(&b)
	c.JSON(200, gin.H{
		"a": b.NestedStructPointer,
		"c": b.FieldC,
	})
}

func GetDataD(c *gin.Context) {
	var d myTypes.StructD
	c.Bind(&d)
	c.JSON(200, gin.H{
		"x": d.NestedAnonyStruct,
		"d": d.FieldD,
	})
}
