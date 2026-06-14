package response

import (
	"github.com/gin-gonic/gin"
)

type Envelope struct {
    Data  any       `json:"data,omitempty"`
    Error *APIError `json:"error,omitempty"`
}
type APIError struct{ Message string `json:"message"` }

func OK(c *gin.Context, status int, data any) {
	c.JSON(status, Envelope{Data: data})
}
func Error(c *gin.Context, status int, msg string) {
	c.JSON(status, Envelope{Error: &APIError{Message: msg}})
}
