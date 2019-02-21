package middleware

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"go-eladmin/common"
	"go-eladmin/log_writer"
	"go-eladmin/model"
	"net/http"
	"strings"
	"time"

	log "github.com/flywithbug/log4go"
	"github.com/gin-gonic/gin"
)

//对照表 	rid		id		m		c 			l		p
// 		xReqId  userId 	method 	statusCode	latency	path
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := log_writer.GetLog()
		l.ReSet()
		l.ResponseCode = http.StatusOK
		// Start timer
		start := time.Now()
		l.RequestId = c.Request.Header.Get(common.KeyContextRequestId)
		if l.RequestId == "" {
			l.RequestId = GenReqID()
		}
		l.UUID = c.Request.Header.Get(common.KeyUUID)
		c.Header(common.KeyContextRequestId, l.RequestId)
		c.Set(common.KeyContextRequestId, l.RequestId)
		l.StartTime = start.UnixNano()
		l.ClientIp = c.ClientIP()
		l.Method = c.Request.Method
		path := c.Request.URL.String()
		if l.Method != "GET" && l.Method != "OPTIONS" {
			l.Path = path
		} else {
			paths := strings.Split(path, "?")
			if len(paths) == 2 {
				l.Para = paths[1]
			}
			l.Path = paths[0]
		}
		methodColor := colorForMethod(l.Method)
		log.NeverShow(l, "[GIN] [%s] [Start]\t%s %s %s|\t%s|\t%s",
			l.RequestId,
			methodColor, l.Method, reset,
			l.Path,
			l.ClientIp)
		//----====----
		c.Next()

		end := time.Now()
		l.EndTime = end.UnixNano()
		l.Latency = end.Sub(start)
		l.StatusCode = c.Writer.Status()
		statusColor := colorForStatus(l.StatusCode)
		//comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
		l.UserId = common.UserId(c)
		l.Info = fmt.Sprintf("[GIN] [%s] [Completed]\t%s|%5d\t%3d|\t3d\t%13v|\t%s\t%s",
			l.RequestId,
			l.Method,
			l.UserId,
			l.StatusCode,
			l.Latency,
			l.Path,
			//comment,
		)

		para := common.Para(c)
		if para != nil {
			l.Para = para
		}
		resI := common.Response(c)
		if resI != nil {
			aRes, ok := resI.(*model.Response)
			if ok {
				l.ResponseCode = aRes.Code
				//TODO 数据返回不正确时记录response
				if aRes.Code != 200 {
					l.Response = aRes
				}
			}
		}
		log.InfoExt(l, "[GIN] [%s] [Completed]\t%s %s %s|\t%8v|\t%s|\t%5d|\t%s%3d%s|\t%s",
			l.RequestId,
			methodColor, l.Method, reset,
			l.Latency,
			l.Path,
			l.UserId,
			statusColor, l.StatusCode, reset,
			comment,
		)
	}
}

var pid = uint32(time.Now().UnixNano() % 4294967291)

// GenReqID is a random generate string func
func GenReqID() string {
	var b [12]byte
	binary.LittleEndian.PutUint32(b[:], pid)
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	return base64.URLEncoding.EncodeToString(b[:])
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// ErrorLogger func
func ErrorLogger() gin.HandlerFunc {
	return ErrorLoggerT(gin.ErrorTypeAny)
}

// ErrorLoggerT func
func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		errors := c.Errors.ByType(typ)
		if len(errors) > 0 {
			c.JSON(-1, errors)
		}
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

var (
	skipPaths = []string{""}
)
