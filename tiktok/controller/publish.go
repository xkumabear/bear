package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"tiktok/common"
	"tiktok/dao"
	"tiktok/dto"
)

//type VideoListResponse struct {
//	Response
//	VideoList []Video `json:"video_list"`
//}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	params := &dto.PublishInput{}
	out := &dto.PublishOutput{}
	if err := params.GetValidParams(c); err != nil {
		out.ResponseError(common.ParamsErr, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	if _, exist := usersLoginInfo[params.Token]; !exist {
		out.ResponseError(common.ParamsErrExist, common.ParamsErrMsg)
		c.JSON(http.StatusOK, out)
		return
	}
	fmt.Println(params)

	//data, err := c.FormFile("data")
	//if err != nil {
	//	c.JSON(http.StatusOK, Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}

	filename := filepath.Base(params.Data.Filename)
	user := usersLoginInfo[params.Token]
	finalName := fmt.Sprintf("%d_%s", user.Model.ID, filename)
	saveFile := filepath.Join("./public/", finalName)
	finalUrl := fmt.Sprintf("%s/public/%s", common.Url, finalName)
	video := &dao.Video{User: user, PlayUrl: finalUrl, Title: params.Title}
	video.Upload()
	if err := c.SaveUploadedFile(params.Data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	//c.JSON(http.StatusOK, VideoListResponse{
	//	Response: Response{
	//		StatusCode: 0,
	//	},
	//	VideoList: DemoVideos,
	//})
}
