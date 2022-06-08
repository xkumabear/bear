package dto

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"tiktok/common"
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

type PublishListInput struct {
	Token  string `json:"token" form:"token"`
	UserID int64  `json:"user_id" form:"user_id" `
}

func (params *PublishListInput) GetValidParams(c *gin.Context) error {
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

type PublishListOutput struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
}

func (u *PublishListOutput) ResponseError(statusCode int32, statusMsg string) {
	u.Response = Response{StatusCode: statusCode, StatusMsg: statusMsg}
}

func (u *PublishListOutput) ResponseSuccess(VideoList *[]Video) {
	u.Response = Response{StatusCode: common.SuccessCode, StatusMsg: ""}
	u.VideoList = *VideoList
}
