package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"tiktok/common"
	"tiktok/dto"
)

type User struct {
	gorm.Model
	Name          string `gorm:"DEFAULT:'未定义'"`
	Username      string
	Password      string
	FollowCount   int64  `gorm:"DEFAULT:0"`
	FollowerCount int64  `gorm:"DEFAULT:0"`
	FollowList    string `gorm:"DEFAULT:''"`
}

func init() {

}
func (u *User) conn() *gorm.DB {
	db, err := gorm.Open(common.DRIVER, common.DSN)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})
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
func (u *User) Save(db *gorm.DB) error {

	return db.Save(u).Error
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
	if err == nil {
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
