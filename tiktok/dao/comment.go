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
	UserID int64
	//User          User
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

//func (c *comment) FindVideo(db *gorm.DB, search *comment) (*comment, error) {
//	fmt.Println(search)
//	//var user User
//	//err := db.Where(search).Find(&user).Error
//	//if err != nil {
//	//	return nil, err
//	//}
//	return comment
//}

func (c *Comment) CommentAdd(param *dto.CommentActionInput) error {
	db := c.conn()
	defer db.Close()
	//videoid, err := c.FindVideo(db, &comment{VideoId: param.VideoId}) //, IsDelete: 0
	//if err == nil || user != nil {
	//	return user, errors.New("已存在该用户，不可重复注册。") //打印堆栈
	//}

	c.UserID = param.UserId
	c.CommentText = param.CommentText
	c.VideoId = param.VideoId

	err := db.Save(c).Error
	//err = c.Save(db)
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) CommentDelte(param *dto.CommentActionInput) error {
	db := c.conn()
	defer db.Close()
	//videoid, err := c.FindVideo(db, &comment{VideoId: param.VideoId}) //, IsDelete: 0
	//if err == nil || user != nil {
	//	return user, errors.New("已存在该用户，不可重复注册。") //打印堆栈
	//}

	//c.UserID = param.UserId
	//c.CommentText = param.CommentText
	//c.VideoId = param.VideoId
	c.ID = uint(param.CommentId)
	err := db.Delete(c).Error
	//err = c.Save(db)
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) VideoCommentList(param *dto.CommentListRequire) (*dto.CommentListResponse, error) {
	db := c.conn()
	defer db.Close()
	//var clr dto.CommentListResponse
	//用videoid 查库，返回CommentList
	//fmt.Println(c)
	//var cl CommentList
	var clr dto.CommentListResponse

	comments, err := c.FindCommentList(db, &CommentList{VideoId: param.VideoId})
	clr.CommentList = *comments
	fmt.Println("VideoCommentList")
	fmt.Println(clr.CommentList)
	//clr.CommentList = cl.CommentList
	//c.ID = uint(param.CommentId)
	//err := db.Delete(c).Error
	//err = c.Save(db)
	if err != nil {
		return &clr, err
	}
	return &clr, nil
}
