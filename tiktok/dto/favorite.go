package dto

import (
	"github.com/gin-gonic/gin"
	"tiktok/common"
)

type FavoriteInput struct {
	Token      string `json:"token" form:"token"`
	VideoID    int64  `json:"video_id" form:"video_id" `
	ActionType int64  `json:"action_type" form:"action_type" `
}

func (params *FavoriteInput) GetValidParams(c *gin.Context) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	return nil
}

type FavoriteOutput struct {
	Response
}

func (u *FavoriteOutput) ResponseSuccess() {
	u.Response = Response{StatusCode: common.SuccessCode, StatusMsg: ""}
}

func (u *FavoriteOutput) ResponseError(statusCode int32, statusMsg string) {
	u.Response = Response{StatusCode: statusCode, StatusMsg: statusMsg}
}

type FavoriteListInput struct {
	UserID string `json:"user_id" form:"user_id" `
	Token  string `json:"token" form:"token"`
}

func (params *FavoriteListInput) GetValidParams(c *gin.Context) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	return nil
}

type FavoriteListOutput struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
}

func (u *FavoriteListOutput) ResponseError(statusCode int32, statusMsg string) {
	u.Response = Response{StatusCode: statusCode, StatusMsg: statusMsg}
}

func (u *FavoriteListOutput) ResponseSuccess(VideoList []Video) {
	u.Response = Response{StatusCode: common.SuccessCode, StatusMsg: ""}
	u.VideoList = VideoList
}
