package controller

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"tiktok/common"
	"tiktok/dao"
	"tiktok/dto"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	c.JSON(http.StatusOK, token)
	//actionType := c.Query("action_type")
	params := &dto.CommentActionInput{}
	out := &dto.CommentActionResponse{}
	if err := params.GetValidParams(c); err != nil {
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	if _, exist := usersLoginInfo[params.Token]; exist {
		//发布评论
		if params.ActionType == "1" {
			//text := c.Query("comment_text")
			//dao写入数据库
			com := &dao.Comment{}
			err := com.CommentAdd(params)
			if err != nil {
				out.ResponseError(common.SqlAddErr, common.SqlAddErrMsg)
				c.JSON(http.StatusOK, out)
				return
			}

			//dto返回json
			out = &dto.CommentActionResponse{
				Response: dto.Response{
					StatusCode: 0,
					StatusMsg:  "Msg sent",
				},
				Comment: dto.Comment{
					Id: 0,
					User: dto.User{
						Id:            0,
						Name:          "wwwww",
						FollowCount:   0,
						FollowerCount: 0,
						IsFollow:      false,
					},
					Content:    "",
					CreateDate: "",
				},
			}
		}
		//删除评论
		if params.ActionType == "2" {
			//删除数据库中的某些值
			com := &dao.Comment{}
			err := com.CommentDelte(params)
			if err != nil {
				out.ResponseError(common.SqlAddErr, common.SqlAddErrMsg)
				c.JSON(http.StatusOK, out)
				return
			}

			out = &dto.CommentActionResponse{
				Response: dto.Response{
					StatusCode: 0,
					StatusMsg:  "Msg deleted",
				},
				Comment: dto.Comment{
					Id: 0,
					User: dto.User{
						Id:            0,
						Name:          "wwwww",
						FollowCount:   0,
						FollowerCount: 0,
						IsFollow:      false,
					},
					Content:    "",
					CreateDate: "",
				},
			}
		}
		c.JSON(http.StatusOK, out)
		return
	}
	c.JSON(http.StatusOK, out)
	return

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	//获取参数，校验
	token := c.Query("token")
	video_id := c.Query("video_id")

	if _, exist := usersLoginInfo[token]; exist {
		if video_id == "1" {
			text := c.Query("comment_text")

			//写入数据库

			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
				Comment: Comment{
					Id: 1,
					//User:       user,
					Content:    text,
					CreateDate: "05-01",
				}})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

	//直接赋值

	//查库赋值（先不写）

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: DemoComments,
	})
}
