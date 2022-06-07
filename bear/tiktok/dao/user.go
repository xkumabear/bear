package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"strings"
	"tiktok/common"
	"tiktok/dto"
)

type User struct {
	gorm.Model
	Name string `gorm:"DEFAULT:'未定义'"`
	//Video         []Video `gorm:"ForeignKey:UserId;AssociationForeignKey:ID"`
	Username      string
	Password      string
	FollowCount   int64  `gorm:"DEFAULT:0"`
	FollowerCount int64  `gorm:"DEFAULT:0"`
	FollowList    string `gorm:"DEFAULT:''"`
	FollowerList  string `gorm:"DEFAULT:''"`
	IsFollow      int64  `gorm:"DEFAULT:0"`
}

func init() {

}
func (u *User) conn() *gorm.DB {
	db, err := gorm.Open(common.DRIVER, common.DSN)
	if err != nil {
		panic(err)
	}
	//db.AutoMigrate(&User{})
	return db
}

func (u *User) Find(db *gorm.DB, search *User) (*User, error) {
	fmt.Println(search)
	var user User
	err := db.Where(search).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (u *User) Search(db *gorm.DB, id uint) error {

	err := db.Where("id=?", id).Find(u).Error
	//fmt.Println(u)
	if err != nil { //有错误
		return err
	}

	return nil // 空  找到了

}

func (u *User) Save(db *gorm.DB) error {

	return db.Create(u).Error
}

func (u *User) Register(param *dto.RegisterInput) (*User, error) {
	db := u.conn()
	defer db.Close()
	user, err := u.Find(db, &User{Username: param.Username}) //, IsDelete: 0
	if err == nil || user != nil {
		return user, errors.New("已存在该用户，不可重复注册。") //打印堆栈
	}

	u.Name = param.Name
	u.Username = param.Username
	saltPassword := common.MD5(param.Password)
	u.Password = saltPassword
	err = u.Save(db)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *User) LoginCheck(param *dto.LoginInput) (*User, error) {
	db := u.conn()
	defer db.Close()
	user, err := u.Find(db, &User{Username: param.Username}) //, IsDelete: 0
	if err != nil {
		return nil, errors.New("用户名错误！") //打印堆栈
	}
	saltPassword := common.MD5(param.Password)
	if user.Password != saltPassword {
		return nil, errors.New("密码错误！")
	}
	return user, nil
}

func (u *User) GetUsersList(param *dto.FollowListInput) *dto.FollowOutput {
	out := &dto.FollowOutput{}

	db := u.conn()
	defer db.Close()

	follows := u.FollowList //得到关注列表 字符串
	userIds := strings.Split(follows, "#")
	var userList []dto.User

	for i := 0; i < len(userIds); i++ {
		var user User
		x, _ := strconv.Atoi(userIds[i])
		id := uint(x)

		if err := user.Search(db, id); err != nil {
			continue
		}
		fmt.Println(user)
		userList = append(userList, dto.User{
			Id:            int64(user.Model.ID),
			Name:          user.Username,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      user.IsFollow == 1,
		})
	}

	out.UserList = userList
	return out
}

func (u *User) GetFollowerList(param *dto.FollowListInput) *dto.FollowOutput {
	out := &dto.FollowOutput{}

	db := u.conn()
	defer db.Close()

	followers := u.FollowerList //得到关注列表 字符串
	userIds := strings.Split(followers, "#")
	var userList []dto.User

	for i := 0; i < len(userIds); i++ {
		var user User
		x, _ := strconv.Atoi(userIds[i])
		id := uint(x)

		if err := user.Search(db, id); err != nil {
			continue
		}
		fmt.Println(user)
		userList = append(userList, dto.User{
			Id:            int64(user.Model.ID),
			Name:          user.Username,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      user.IsFollow == 1,
		})
	}

	out.UserList = userList
	return out
}
