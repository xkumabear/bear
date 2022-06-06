package dto

import "github.com/gin-gonic/gin"

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type RegisterInput struct {
	Name     string `json:"name" form:"name"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password" `
}

func (params *RegisterInput) GetValidParams(c *gin.Context) error {
	if err := c.ShouldBind(params); err != nil { //获得有效的参数
		return err
	}
	return nil
}

type LoginInput struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password" `
}

func (params *LoginInput) GetValidParams(c *gin.Context) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	return nil
}

type UserInfoInput struct {
	UserID string `json:"user_id" form:"user_id"`
	Token  string `json:"token" form:"token" `
}

func (params *UserInfoInput) GetValidParams(c *gin.Context) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	return nil
}

type UserLoginResponse struct {
	Response        //返回的响应 状态码   和   返回描述
	UserId   int64  `json:"user_id"`         //返回用户的id
	Token    string `json:"token,omitempty"` //用户鉴权
}

func (u *UserLoginResponse) ResponseError(statusCode int32, statusMsg string) {
	u.Response = Response{StatusCode: statusCode, StatusMsg: statusMsg}
}

type UserOutput struct {
	Response
	User User `json:"user"`
}

func (u *UserOutput) ResponseError(statusCode int32, statusMsg string) {
	u.Response = Response{StatusCode: statusCode, StatusMsg: statusMsg}
}
