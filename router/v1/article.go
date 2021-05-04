package v1

import (
	"awesomeProject/models"
	"awesomeProject/pkg/app"
	"awesomeProject/services/article_service"
	"awesomeProject/services/tag_service"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"time"
)

type AddArticleForm struct {
	Id    int       `form:"id" json:"id" binding:"required"`
	Title string    `form:"title" json:"title" binding:"required,min=1,max=100"`
	Date  time.Time `form:"date" json:"date" binding:"required"`
	Body  string    `form:"body" json:"body" binding:"required"`
	Tags  []string  `form:"tags" json:"tags" binding:"required,gt=0"`
}

func AddArticle(c *gin.Context) {
	var form AddArticleForm
	var appG = app.Gin{C: c}

	httpCode, errBindAndValid := app.BindAndValid(c, &form)
	if errBindAndValid != nil {
		appG.Response(httpCode, errBindAndValid.Error(), nil)
		return
	}

	a := models.Article{
		Id:    form.Id,
		Title: form.Title,
		Date:  form.Date,
		Body:  form.Body,
		Tags:  form.Tags,
	}
	errArticle := article_service.AddArticle(&a)
	if errArticle != nil {
		appG.Response(http.StatusInternalServerError, errArticle.Error(), nil)
		return
	}

	// NOTE:

	errTags := tag_service.AddArticle(a)
	if errTags != nil {
		appG.Response(http.StatusInternalServerError, errTags.Error(), nil)
	}

	appG.Response(http.StatusOK, "article saved", nil)
}

func GetArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()

	if ok := validId(id); !ok {
		appG.Response(http.StatusBadRequest, "invalid id", nil)
		return
	}

	a, err := article_service.Get(id)
	if err != nil {
		appG.Response(http.StatusNotFound, "no article with that ID found", nil)
		return
	}

	appG.Response(http.StatusOK, "", a)
}

func validId(id int) bool {
	return id > 0
}
