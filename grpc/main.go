package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)
type RespBody struct {
	Ip string `json:"ip"`
	Port uint32 `json:"port"`
	Protocol string `json:"protocol"`
}
var reqInfo RespBody


func test(Protocol,Ip string,Port uint32) (resp uint32,reason string){
	fmt.Println("获取到的信息：", Protocol, Ip, Port)
	return 200,"OK"
}

func main(){
	r := gin.Default()
	r.POST("/ping", func(c *gin.Context) {
		c.BindJSON(&reqInfo)
		resp,reason  := test(reqInfo.Protocol,reqInfo.Ip,reqInfo.Port)
		c.JSON(http.StatusOK,gin.H{
			"status_code":resp,
			"reason":reason,
		})
	})
	r.Run(":8088")
}