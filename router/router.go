package router

import (
	v1 "awesomeProject/router/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/v1")
	apiV1.POST("/articles", v1.AddArticle)
	apiV1.GET("/articles/:id", v1.GetArticle)
	apiV1.GET("/tags/:tagName/:date", v1.GetTaggedArticles)

	return r
}
