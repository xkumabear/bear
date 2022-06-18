package controller

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
	"tiktok/common"
	"tiktok/dao"
	"tiktok/dto"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	params := &dto.PublishInput{}
	out := &dto.PublishOutput{}
	if err := params.GetValidParams(c); err != nil {
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}

	user, err := CheckToken(params.Token)
	if err != nil {
		out.ResponseError(common.ParamsErrExist, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	filename := filepath.Base(params.Data.Filename)
	finalName := fmt.Sprintf("%d_%s", user.Model.ID, filename)
	saveFile := filepath.Join("./public/", finalName)
	finalUrl := fmt.Sprintf("%s/static/%s", common.Url, finalName)

	//cover
	if err := c.SaveUploadedFile(params.Data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	reader := ReadFrameAsJpeg(saveFile, 5)
	img, err := imaging.Decode(reader)
	if err != nil {
		out.ResponseError(common.ParamsErrExist, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	name := strings.Split(filename, ".")
	finalImgName := fmt.Sprintf("%d_%s.jpeg", user.Model.ID, name[0])
	saveImgFile := filepath.Join("./public/", finalImgName)
	finalImgUrl := fmt.Sprintf("%s/static/%s", common.Url, finalImgName)
	err = imaging.Save(img, saveImgFile)
	if err != nil {
		out.ResponseError(common.ParamsErrExist, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}

	video := &dao.Video{User: user, PlayUrl: finalUrl, CoverUrl: finalImgUrl, Title: params.Title}
	video.Upload()

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	params := &dto.PublishListInput{}
	out := &dto.PublishListOutput{}
	if err := params.GetValidParams(c); err != nil {
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	user, err := CheckToken(params.Token)
	if err != nil {
		out.ResponseError(common.ParamsErrExist, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	//userIdString := strconv.FormatInt(int64(user.Model.ID), 10) + "#"
	userIdString := fmt.Sprintf("%010d#", user.Model.ID)
	video := &dao.Video{}
	videoList, err := video.PublishVideoList(params)
	if err != nil {
		out.ResponseError(common.SqlFindErr, common.SqlFindErrMsg)
	}

	var outVideoList []dto.Video
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

	out.ResponseSuccess(&outVideoList)
	c.JSON(http.StatusOK, out)
}
