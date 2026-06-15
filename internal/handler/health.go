package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wchabir/taskflow/pkg/response"
	"net/http"
)

func Health(c *gin.Context) {
	response.OK(c, http.StatusOK, gin.H{"status": "ok"})
}
