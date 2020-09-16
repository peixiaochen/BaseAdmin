package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/peixiaochen/BaseAdmin/app/Requests"
	"github.com/peixiaochen/BaseAdmin/app/models"
	"github.com/peixiaochen/BaseAdmin/pkg/context"
)

func AdminOperationLogIndex(c *gin.Context) {
	var (
		Response Response
		Request  *Requests.AdminOperationLogIndexRequest
		Model    models.Model
	)
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Msg = err.Error()
		Response.Code = context.CodeClientError
		Response.ServerJson(c)
		return
	}
	Model = &models.AdminOperationLogModel{}
	if Data, err := Model.GetAll(map[string]interface{}{
		"Keywords":  "%" + Request.Data.Keywords + "%",
		"StartTime": Request.Data.StartTime,
		"EndTime":   Request.Data.EndTime,
	}, (Request.Data.Page-1)*uint64(Request.Data.Limit), Request.Data.Limit); err != nil {
		Response.Code = context.CodeServerError
		Response.Msg = err.Error()
	} else {
		Response.Data = Data
	}
	Response.ServerJson(c)
	return
}
