package handler_common

import (
	"fmt"
	"go-eladmin/common"

	"github.com/gin-gonic/gin"
)

type StateType int

const (
	RouterTypeNormal StateType = iota
	RouterTypeNeedAuth
)

type GinHandleFunc struct {
	Handler    gin.HandlerFunc
	RouterType StateType
	Method     string
	Route      string
}

func RequestId(c *gin.Context) string {
	return fmt.Sprintf("【rid:%s】", common.XRequestId(c))
}
