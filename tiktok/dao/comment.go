package dao

//database
import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"tiktok/common"
	"tiktok/dto"
)

type Comment struct {
	gorm.Model
	UserID      int64
	User        User
	CommentText string
	//CommentId   int64
	VideoId  int64
	IsDelete bool `gorm:"DEFAULT:false"`
}

type CommentList struct {
	gorm.Model
	CommentList []Comment `json:"comment_list"` // 评论列表
	VideoId     int64
}

//type VideoComment struct {
//	gorm.Model
//	VideoId int64
//}

//type follow struct {
//	gorm.Model
//	Id       int64
//	UserAID  int64 `gorm:"ForeignKey:UserID;AssociationForeignKey:Id"`
//	UserBID  int64
//	User     User
//	Relation int
//}
//
//type follower struct {
//	gorm.Model
//	Id       int64
//	UserAID  int64
//	UserBID  int64 `gorm:"ForeignKey:UserID;AssociationForeignKey:Id"`
//	User     User
//	Relation int
//}

func init() {

}
func (c *Comment) conn() *gorm.DB {
	db, err := gorm.Open(common.DRIVER, common.DSN)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Comment{})
	return db
}

func (c *Comment) FindCommentList(db *gorm.DB, search *CommentList) (*[]dto.Comment, error) {
	fmt.Println(search)
	//var cl CommentList
	var com []dto.Comment
	err := db.Where(search.VideoId).Find(&com).Error
	if err != nil {
		return nil, err
	}
	//cl.CommentList = comments
	fmt.Println("FindCommentList")
	fmt.Println(com)
	return &com, nil
}

//func (c *VideoComment) conn() *gorm.DB {
//	db, err := gorm.Open(common.DRIVER, common.DSN)
//	if err != nil {
//		panic(err)
//	}
//	db.AutoMigrate(&Comment{})
//	return db
//}

func (c *Comment) FindVideo(db *gorm.DB, videoid int64) (Video, error) {

	var video Video
	err := db.Where("id = ?", videoid).Find(&video).Preload("User").Error
	if err != nil {
		return video, err
	}
	return video, nil
}

func (c *Comment) CommentAdd(userid int64, param *dto.CommentActionInput) (User, error) {
	db := c.conn()
	defer db.Close()
	video, err := c.FindVideo(db, param.VideoId)
	if err != nil {
		return video.User, err
	}

	c.UserID = userid
	c.CommentText = param.CommentText
	c.VideoId = param.VideoId
	video.CommentCount++
	err = db.Save(video).Error
	if err != nil {
		return video.User, err
	}
	err = db.Save(c).Error
	//err = c.Save(db)
	if err != nil {
		return video.User, err
	}
	return video.User, nil
}

func (c *Comment) CommentDelte(userid int64, param *dto.CommentActionInput) (User, error) {
	db := c.conn()
	defer db.Close()
	video, err := c.FindVideo(db, param.VideoId)
	if err != nil {
		return video.User, err
	}
	c.UserID = userid
	c.CommentText = param.CommentText
	c.VideoId = param.VideoId
	c.ID = uint(param.CommentId)
	err = db.Delete(c).Error
	//err = c.Save(db)
	if err != nil {
		return video.User, err
	}
	return video.User, nil
}

func (c *Comment) VideoCommentList(userid int64, param dto.CommentListRequire) (*[]Comment, User, error) {
	db := c.conn()
	defer db.Close()
	var commentlist []Comment
	video, err := c.FindVideo(db, param.VideoId)
	c.UserID = userid
	c.VideoId = param.VideoId
	if err != nil {
		return nil, video.User, err
	}

	err = db.Model(commentlist).Where("video_id = ?", param.VideoId).Preload("User").Find(&commentlist).Error
	if err != nil {
		return &commentlist, video.User, err
	}
	return &commentlist, video.User, nil
}

func (c *Comment) VideoIdCheck(param Comment) (*Video, error) {
	db := c.conn()
	defer db.Close()
	var video *Video
	video = new(Video)
	video.Model.ID = uint(param.VideoId)
	err := db.Where(video).Find(&video).Error
	if err != nil {
		return nil, err
	}
	return video, nil
}
