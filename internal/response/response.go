package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Code 				int 	`json:"code"`
	Message 		string `json:"message"`
	Data 				any		`json:"data"`
	TraceID 		string `json:"trace_id"`
	RequestID 	string `json:"request_id"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{
		Code: 0,
		Message: "success",
		Data: data,
		TraceID: GetTraceID(c),
		RequestID: GetRequestID(c),
	})
}

func Fail(c *gin.Context, status int, code int, message string, data any) {
	c.JSON(status, Envelope{
		Code: code,
		Message: message,
		Data: data,
		TraceID: GetTraceID(c),
		RequestID: GetRequestID(c),
	})
}

func GetRequestID(c *gin.Context) string {
	value,_ := c.Get("request_id")
	id,_ := value.(string)
	return id;
}

func GetTraceID(c *gin.Context) string {
	value,_ := c.Get("trace_id")
	id,_ := value.(string)
	return id;
}