package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/common"
	"tiktok/dao"
	"tiktok/dto"
)

type UserListResponse struct {
	Response
	UserList []dao.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")

	fmt.Println("Relation:")
	if user, exist := usersLoginInfo[token]; exist {
		params := &dto.RelationInput{}
		if err := params.GetValidParams(c); err != nil { //获得有效参数  参数是否有错
			out := &dto.Response{StatusCode: common.ParamsErr, StatusMsg: common.ParamsErrMsg}
			c.JSON(http.StatusOK, out)
			return
		}
		fmt.Println(token)

		follow := &dao.Follow{}
		params.UserAID = user.ID
		err := follow.RelationCheck(params)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "to_user_id doesn't exist"})
			return
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	params := &dto.FollowListInput{}
	if err := params.GetValidParams(c); err != nil { //获得有效参数  参数是否有错
		out := &dto.Response{StatusCode: common.ParamsErr, StatusMsg: common.ParamsErrMsg}
		c.JSON(http.StatusOK, out)
		return
	}

	user, exist := usersLoginInfo[token]
	if !exist {
		out := &dto.Response{StatusCode: common.ParamsErrExist, StatusMsg: common.ParamsErrMsg}
		c.JSON(http.StatusOK, out)
		return
	}

	out := user.GetUsersList(params)
	c.JSON(http.StatusOK, out)
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	params := &dto.FollowListInput{}
	if err := params.GetValidParams(c); err != nil { //获得有效参数  参数是否有错
		out := &dto.Response{StatusCode: common.ParamsErr, StatusMsg: common.ParamsErrMsg}
		c.JSON(http.StatusOK, out)
		return
	}

	user, exist := usersLoginInfo[token]
	if !exist {
		out := &dto.Response{StatusCode: common.ParamsErrExist, StatusMsg: common.ParamsErrMsg}
		c.JSON(http.StatusOK, out)
		return
	}

	out := user.GetFollowerList(params)
	c.JSON(http.StatusOK, out)
}
