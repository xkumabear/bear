package controller

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
	"tiktok/common"
	"tiktok/dao"
	"time"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author" gorm:"foreignkey:authorId"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	dao.User
	availableTime time.Time
}

var UsersLoginInfo = make(map[string]User)

func CheckToken(token string) (dao.User, error) {
	now := time.Now()
	if user, exist := UsersLoginInfo[token]; exist && now.Before(user.availableTime) {
		return user.User, nil
	}
	return dao.User{}, errors.New("")
}

func SetToken(username string, user dao.User) string {

	token := GetSaltString(common.TokenKey, username)
	now := time.Now()
	mm, _ := time.ParseDuration("60m")
	deadline := now.Add(mm)
	UsersLoginInfo[token] = User{user, deadline}
	return token
}

func UpdateTokenInfo(token string, user dao.User) error {
	userInfo, exist := UsersLoginInfo[token]
	if !exist {
		return errors.New("token is not exist!")
	}
	UsersLoginInfo[token] = User{user, userInfo.availableTime}
	return nil
}

func GetSaltString(salt, token string) string {
	s1 := sha256.New()
	s1.Write([]byte(token))
	str1 := fmt.Sprintf("%x", s1.Sum(nil))
	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))
	return fmt.Sprintf("%x", s2.Sum(nil))
}

func ReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		println(inFileName)
		panic(err)
	}
	return buf
}
