package dto

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type PublishInput struct {
	Data  *multipart.FileHeader `json:"data" form:"data"`
	Token string                `json:"token" form:"token"`
	Title string                `json:"title" form:"title" `
}

func (params *PublishInput) GetValidParams(c *gin.Context) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	return nil
}

type PublishOutput struct {
	Response
}

func (u *PublishOutput) ResponseError(statusCode int32, statusMsg string) {
	u.Response = Response{StatusCode: statusCode, StatusMsg: statusMsg}
}
