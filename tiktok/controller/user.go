package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"tiktok/common"
	"tiktok/dao"
	"tiktok/dto"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

func Register(c *gin.Context) {
	params := &dto.RegisterInput{}
	out := &dto.UserLoginResponse{}
	if err := params.GetValidParams(c); err != nil {
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	user := &dao.User{}
	users, err := user.Register(params)
	if err != nil {
		out.ResponseError(common.ParamsErrExist, err.Error())
		out.UserId = int64(users.Model.ID)
		c.JSON(http.StatusOK, out)
		return
	}
	token := SetToken(params.Username, *user)
	out.Response = dto.Response{StatusCode: common.SuccessCode, StatusMsg: ""}
	out.Token = token
	out.UserId = int64(user.Model.ID)
	c.JSON(http.StatusOK, out)
}

func Login(c *gin.Context) {
	params := &dto.LoginInput{}
	out := &dto.UserLoginResponse{}
	if err := params.GetValidParams(c); err != nil {
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}

	user := &dao.User{}
	users, err := user.LoginCheck(params)
	if err != nil {
		out.ResponseError(common.ParamsErrExist, err.Error())
		c.JSON(http.StatusOK, out)
		return
	}

	token := SetToken(params.Username, *users)
	out.Response = dto.Response{StatusCode: common.SuccessCode, StatusMsg: ""}
	out.Token = token
	out.UserId = int64(users.Model.ID)
	c.JSON(http.StatusOK, out)

}

func UserInfo(c *gin.Context) {
	params := &dto.UserInfoInput{}
	out := &dto.UserOutput{}
	if err := params.GetValidParams(c); err != nil {
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	//token := c.Query("token")
	uid, _ := strconv.Atoi(params.UserID)
	loginUser, err := CheckToken(params.Token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't login!"})
		return
	}
	if loginUser.ID == uint(uid) {
		c.JSON(http.StatusOK, dto.UserOutput{
			Response: dto.Response{StatusCode: 0},
			User: dto.User{
				Id:            int64(loginUser.Model.ID),
				Name:          loginUser.Name,
				FollowCount:   loginUser.FollowCount,
				FollowerCount: loginUser.FollowerCount,
				IsFollow:      false,
			},
		})
		return
	}

	user := &dao.User{}
	userinfo, err := user.GetUserInfo(uint(uid))
	if err != nil {
		out.ResponseError(common.ParamsErrExist, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
	}

	userIdString := fmt.Sprintf("%010d#", loginUser.Model.ID)
	isFollow := strings.Contains(userinfo.FollowerList, userIdString)
	c.JSON(http.StatusOK, dto.UserOutput{
		Response: dto.Response{StatusCode: 0},
		User: dto.User{
			Id:            int64(userinfo.Model.ID),
			Name:          userinfo.Name,
			FollowCount:   userinfo.FollowCount,
			FollowerCount: userinfo.FollowerCount,
			IsFollow:      isFollow,
		},
	})

	out.ResponseError(common.ParamsErrExist, common.ParamsErrMsg)
	c.JSON(http.StatusOK, out)
	return
}
