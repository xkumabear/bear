package dto

import (
	"github.com/gin-gonic/gin"
)

//io
type CommentActionInput struct {
	UserId      int64  `json:"user_id,omitempty" form:"user_id"`
	Token       string `json:"token"  form:"token"`
	VideoId     int64  `json:"video_id,omitempty"  form:"video_id"`
	ActionType  string `json:"action_type,omitempty"  form:"action_type"`
	CommentText string `json:"comment_text,omitempty"  form:"comment_text"`
	CommentId   int64  `json:"comment_id,omitempty"  form:"comment_id"`
}

func (params *CommentActionInput) GetValidParams(c *gin.Context) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	return nil
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

type CommentListResponse struct {
	CommentList []Comment `json:"comment_list"` // 评论列表
	StatusCode  int64     `json:"status_code"`  // 状态码，0-成功，其他值-失败
	StatusMsg   string    `json:"status_msg"`   // 返回状态描述
}

type CommentListRequire struct {
	VideoId int64 `json:"video_id,omitempty"  form:"video_id"`
}

func (u *CommentActionResponse) ResponseError(statusCode int32, statusMsg string) {
	u.Response = Response{StatusCode: statusCode, StatusMsg: statusMsg}
}
