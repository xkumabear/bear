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
		u, err := follow.RelationCheck(params)
		usersLoginInfo[token] = *u
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "to_user_id doesn't exist"})
			return
		}

		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow-list
func FollowList(c *gin.Context) {
	fmt.Println("FollowList:")

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

	out := &dto.FollowOutput{}
	fmt.Println(user)
	userList, err := user.GetUsersList(params)
	if err != nil {
		c.JSON(http.StatusOK, out)
		return
	}

	var outList []dto.User
	userAIdString := strconv.Itoa(int(user.Model.ID)) //获得 A 的id
	sA := "#" + userAIdString + "#"

	for _, u := range *userList {
		isFollow := strings.Contains(u.FollowList, sA)
		outList = append(outList, dto.User{
			Id:            int64(u.Model.ID),
			Name:          u.Username,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FollowerCount,
			IsFollow:      isFollow,
		})
	}

	out.UserList = outList
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
	out := &dto.FollowOutput{}
	userList, err := user.GetFollowerList(params)
	if err != nil {
		c.JSON(http.StatusOK, out)
		return
	}

	var outList []dto.User
	userAIdString := strconv.Itoa(int(user.Model.ID)) //获得 A 的id
	sA := "#" + userAIdString + "#"

	for _, u := range *userList {
		isFollow := strings.Contains(u.FollowerList, sA)
		outList = append(outList, dto.User{
			Id:            int64(u.Model.ID),
			Name:          u.Username,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FollowerCount,
			IsFollow:      isFollow,
		})
	}

	out.UserList = outList
	c.JSON(http.StatusOK, out)
}
