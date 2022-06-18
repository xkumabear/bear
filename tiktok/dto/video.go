package dto

import (
	"github.com/gin-gonic/gin"
	"tiktok/common"
)

type FeedInput struct {
	Token      string `json:"token" form:"token"`              // 验权token
	LatestTime int64  `json:"latest_time" form:"latest_time" ` // 限制返回视频的最新时间戳
}

func (params *FeedInput) GetValidParams(c *gin.Context) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	return nil
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"Title,omitempty"`
}

type FeedOutput struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

func (u *FeedOutput) ResponseError(statusCode int32, statusMsg string) {
	u.Response = Response{StatusCode: statusCode, StatusMsg: statusMsg}
}

func (u *FeedOutput) ResponseSuccess(VideoList *[]Video, NextTime int64) {
	u.Response = Response{StatusCode: common.SuccessCode, StatusMsg: ""}
	u.NextTime = NextTime
	u.VideoList = *VideoList
}
