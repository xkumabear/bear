package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"tiktok/common"
	"tiktok/dao"
	"tiktok/dto"
	"time"
)

// Feed douyin/feed
func Feed(c *gin.Context) {
	params := &dto.FeedInput{}
	out := &dto.FeedOutput{}
	if err := params.GetValidParams(c); err != nil {
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}

	video := &dao.Video{}
	videoList, err := video.VideoList(params)
	if err != nil {
		out.ResponseError(common.SqlFindErr, common.SqlFindErrMsg)
	}
	var outVideoList []dto.Video
	user, err := CheckToken(params.Token)
	if err != nil {
		for _, item := range *videoList {
			outVideoList = append(outVideoList, dto.Video{
				Id: int64(item.Model.ID),
				Author: dto.User{
					Id:            int64(item.User.Model.ID),
					Name:          item.User.Name,
					FollowCount:   item.User.FollowCount,
					FollowerCount: item.User.FollowerCount,
					IsFollow:      false,
				},
				PlayUrl:       item.PlayUrl,
				CoverUrl:      item.CoverUrl,
				FavoriteCount: item.FavoriteCount,
				CommentCount:  item.CommentCount,
				IsFavorite:    false,
				Title:         item.Title,
			})
		}
		var nextTime int64 = 0
		if l := len((*videoList)); l > 0 {
			nextTime = (*videoList)[0].CreatedAt.Unix()
		} else {
			nextTime = time.Now().Unix()

		}
		out.ResponseSuccess(&outVideoList, nextTime)
		c.JSON(http.StatusOK, out)
		return
	}

	//userIdString := strconv.FormatInt(int64(user.Model.ID), 10) + "#"
	userIdString := fmt.Sprintf("%010d#", user.Model.ID)
	for _, item := range *videoList {
		isFollow := strings.Contains(item.User.FollowerList, userIdString)
		isFavorite := strings.Contains(item.FavoriteList, userIdString)
		outVideoList = append(outVideoList, dto.Video{
			Id: int64(item.Model.ID),
			Author: dto.User{
				Id:            int64(item.User.Model.ID),
				Name:          item.User.Name,
				FollowCount:   item.User.FollowCount,
				FollowerCount: item.User.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:       item.PlayUrl,
			CoverUrl:      item.CoverUrl,
			FavoriteCount: item.FavoriteCount,
			CommentCount:  item.CommentCount,
			IsFavorite:    isFavorite,
			Title:         item.Title,
		})
	}
	var nextTime int64
	if l := len((*videoList)); l > 0 {
		nextTime = (*videoList)[0].CreatedAt.Unix()
	} else {
		nextTime = time.Now().Unix()

	}

	out.ResponseSuccess(&outVideoList, nextTime)
	c.JSON(http.StatusOK, out)
}
