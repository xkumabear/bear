package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/common"
	"tiktok/dao"
	"tiktok/dto"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

var usersLoginInfo = map[string]dao.User{}

func Register(c *gin.Context) {
	fmt.Println("Register:")
	//username := c.Query("username")
	//password := c.Query("password")
	params := &dto.RegisterInput{}                   //创建一个输出对象
	out := &dto.UserLoginResponse{}                  //创建一个返回响应 对象
	if err := params.GetValidParams(c); err != nil { //获得有效参数  参数是否有错
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	//需更改策略
	token := params.Username + params.Password //得到用户鉴权
	user := &dao.User{}                        //创建  用户表 的对象
	users, err := user.Register(params)        //注册对象  存在就返回 该对象 和 错误描述； 不存在 就注册并返回该对象 和 空

	if err != nil {
		out.ResponseError(common.ParamsErrExist, err.Error())
		out.UserId = int64(users.ID) //   ？？？users.ID
		c.JSON(http.StatusOK, out)
		return
	}

	out.Response = dto.Response{StatusCode: common.SuccessCode, StatusMsg: ""}
	out.Token = token
	out.UserId = int64(user.ID)
	c.JSON(http.StatusOK, out)
	fmt.Println("register success")
	return
}

func Login(c *gin.Context) {
	fmt.Println("Login:")
	params := &dto.LoginInput{}
	out := &dto.UserLoginResponse{}
	if err := params.GetValidParams(c); err != nil { //参数是否有错
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	token := params.Username + params.Password
	user := &dao.User{}
	users, err := user.LoginCheck(params)
	if err != nil {

		out.ResponseError(common.ParamsErrExist, err.Error())
		c.JSON(http.StatusOK, out)
		return
	}
	usersLoginInfo[token] = *users
	out.Response = dto.Response{StatusCode: common.SuccessCode, StatusMsg: ""}
	out.Token = token
	out.UserId = int64(users.ID)
	c.JSON(http.StatusOK, out)
	fmt.Println("login success")

}

func UserInfo(c *gin.Context) {
	fmt.Println("UserInfo:")
	params := &dto.UserInfoInput{}
	out := &dto.UserOutput{}
	if err := params.GetValidParams(c); err != nil {
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	//token := c.Query("token")
	if user, exist := usersLoginInfo[params.Token]; exist {
		c.JSON(http.StatusOK, dto.UserOutput{
			Response: dto.Response{StatusCode: 0},
			User: dto.User{
				Id:            int64(user.Model.ID),
				Name:          user.Username,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      true,
			},
		})

		return
	}
	out.ResponseError(common.ParamsErrExist, common.ParamsErrMsg)
	c.JSON(http.StatusOK, out)

	return
}
