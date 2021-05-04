package v1

import (
	"awesomeProject/models"
	"awesomeProject/pkg/app"
	"awesomeProject/pkg/util"
	"awesomeProject/services/tag_service"
	"github.com/gin-gonic/gin"

	"net/http"
)

func GetTaggedArticles(c *gin.Context) {
	appG := app.Gin{C: c}
	tagName := c.Param("tagName")
	dateStr := c.Param("date")

	if !models.ValidateTagName(tagName) {
		appG.Response(http.StatusBadRequest, "invalid tag name", nil)
		return
	}

	date, parseErr := util.ConvertCondensedDateString(dateStr)
	if parseErr != nil {
		appG.Response(http.StatusBadRequest, "invalid date", nil)
		return
	}

	as := tag_service.GetArticles(tagName, date)
	summary := as.GetTagSummary(tagName)
	appG.Response(http.StatusOK, "", summary)
}
