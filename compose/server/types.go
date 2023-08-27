package server

import (
	"github.com/gin-gonic/gin"
)

type Func struct {
	HTTPMethod   string
	RelativePath string
	Handlers     []gin.HandlerFunc
}
