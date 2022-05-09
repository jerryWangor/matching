package handler

import (
	"github.com/gin-gonic/gin"
	"matching/model"
	"net/http"
)

type OpenMatch struct {
	model.TimeSign
	Symbol string `form:"symbol"`
}

func OpenMatching(c *gin.Context) {
	var openmath OpenMatch
	if c.ShouldBind(&openmath) != nil {
		c.JSON(http.StatusOK, gin.H{"request": "msg"})
	}
}