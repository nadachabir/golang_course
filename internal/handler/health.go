package handler

import(
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/wchabir/taskflow/pkg/response"

)


func Health(c *gin.Context){
	response.OK(c, http.StatusOK, gin.H{"status": "ok"})
}