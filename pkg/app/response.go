package app

import "github.com/gin-gonic/gin"

type Gin struct {
	C *gin.Context
}

type Response struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin) Response(httpCode int, msg string, data interface{}) {
	if msg == "" {
		g.C.JSON(httpCode, data)
		return
	}

	g.C.JSON(httpCode, Response{
		Msg:  msg,
		Data: data,
	})
}
