package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Err       error       `json:"err"`
	SucessMsg string      `json:"msg"`
	Data      interface{} `json:"data"`
}

func New(err error, sucess_msg string, data interface{}) *Response {
	return &Response{
		err,
		sucess_msg,
		data,
	}
}

func (r *Response) Do(c *gin.Context) {
	err := r.Err
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "internal error:" + err.Error(),
		})
		return
	}

	data := r.Data

	//success
	if data != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  r.SucessMsg,
			"data": data,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  r.SucessMsg,
		})
	}
}
