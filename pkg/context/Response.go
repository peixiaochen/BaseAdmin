package context

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

const (
	Success                   = 0
	CodeClientError           = 10400
	CodeClientNoLogin         = 10401
	CodeClientPermissionError = 10402
	CodeServerError           = 10500
)

var MsgFlags = map[uint32]string{
	Success:                   "success",
	CodeServerError:           "server fail",
	CodeClientError:           "client fail",
	CodeClientPermissionError: "user is not to access this uri",
	CodeClientNoLogin:         "admin user not login",
}

func (r *Response) GetMsg(code uint32) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[CodeServerError]
}
func (r *Response) ServerJson(c *gin.Context) {
	if r.Msg == "" {
		r.Msg = r.GetMsg(r.Code)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": r.Code,
		"msg":  r.Msg,
		"data": r.Data,
	})
	return
}
