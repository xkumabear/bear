package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"tiktok/common"
	"tiktok/dto"
	"time"
)

type Video struct {
	gorm.Model
	UserID        int64
	User          User `gorm:"ForeignKey:UserID"`
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64
	CommentCount  int64
	FavoriteList  string
	Title         string
}

func (v *Video) conn() *gorm.DB {
	db, err := gorm.Open(common.DRIVER, common.DSN)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Video{})
	return db
}

func (v *Video) VideoList(params *dto.FeedInput) (*[]Video, error) {
	db := v.conn()
	defer db.Close()
	var videoList []Video
	timeLayout := "2006-01-02 15:04:05"
	timeStr := time.Now().Format(timeLayout)
	if params.LatestTime != 0 {
		timeStr = time.Unix(params.LatestTime, 0).Format(timeLayout)
	}

	err := db.Model(videoList).Where("created_at < ?", timeStr).Preload("User").Find(&videoList).Error
	if err != nil {
		return &videoList, err
	}
	return &videoList, nil

}
func (v *Video) PublishVideoList(params *dto.PublishListInput) (*[]Video, error) {
	db := v.conn()
	defer db.Close()
	var videoList []Video

	err := db.Model(videoList).Where("user_id = ?", params.UserID).Preload("User").Find(&videoList).Error
	if err != nil {
		return &videoList, err
	}

	return &videoList, nil
}

func (v *Video) Find(db *gorm.DB, search *Video) (*Video, error) {
	fmt.Println(search)
	var video Video
	err := db.Where(search).Find(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, err
}

func (v *Video) Save(db *gorm.DB) error {
	return db.Save(v).Error
}

func (v *Video) Upload() error {
	db := v.conn()
	defer db.Close()
	return v.Save(db)
}