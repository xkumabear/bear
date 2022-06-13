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
	params := &dto.RelationInput{}
	if err := params.GetValidParams(c); err != nil { //获得有效参数  参数是否有错
		out := &dto.Response{StatusCode: common.ParamsErr, StatusMsg: common.ParamsErrMsg}
		c.JSON(http.StatusOK, out)
		return
	}

	fmt.Println("Relation:")
	user, err := CheckToken(params.Token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	
	if user.ID == params.UserBID {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
		return
	}

	users := &dao.User{}
	err = users.RelationCheck(user.ID, params)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "relation illegal!"})
		return
	}
	//err = UpdateTokenInfo(params.Token, *users)
	//if err != nil {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "token doesn't exist"})
	//	return
	//}

	c.JSON(http.StatusOK, Response{StatusCode: 0})

}

// FollowList all users have same follow-list
func FollowList(c *gin.Context) {
	fmt.Println("FollowList:")

	params := &dto.FollowListInput{}
	if err := params.GetValidParams(c); err != nil { //获得有效参数  参数是否有错
		out := &dto.Response{StatusCode: common.ParamsErr, StatusMsg: common.ParamsErrMsg}
		c.JSON(http.StatusOK, out)
		return
	}
	user, err := CheckToken(params.Token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
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
	userIdString := strconv.FormatInt(int64(user.Model.ID), 10) + "#"

	for _, u := range *userList {
		isFollow := strings.Contains(u.FollowList, userIdString)
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
	params := &dto.FollowListInput{}
	if err := params.GetValidParams(c); err != nil { //获得有效参数  参数是否有错
		out := &dto.Response{StatusCode: common.ParamsErr, StatusMsg: common.ParamsErrMsg}
		c.JSON(http.StatusOK, out)
		return
	}

	user, err := CheckToken(params.Token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
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
