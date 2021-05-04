package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BindAndValid(c *gin.Context, form interface{}) (int, error) {
	err := c.BindJSON(form)
	if err != nil {
		return http.StatusBadRequest, errors.New("invalid parameters")
	}

	return http.StatusOK, nil
}
