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
	FollowerList  string `gorm:"DEFAULT:''"` // 粉丝列表  #id#
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

func (u *User) Search(db *gorm.DB, id uint) (*User, error) {

	var user User
	fmt.Println("id:", id)
	err := db.Model(u).Where("id=?", id).Find(&user).Error

	if err != nil { //有错误
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(id, user)
	return &user, nil // 空  找到了

}

func (u *User) GetUsersList(param *dto.FollowListInput) (*[]User, error) {
	db := u.conn()
	defer db.Close()

	follows := u.FollowList //得到关注列表 字符串
	userIds := strings.Split(follows, "#")
	var userList []User

	size := 0
	fmt.Println("userIds:", userIds)
	for i := 0; i < len(userIds); i++ {
		fmt.Println("i:", i)

		x, _ := strconv.Atoi(userIds[i])
		id := uint(x)

		user, err := u.Search(db, id)
		if err != nil {
			continue
		}
		fmt.Println("user", user)
		userList = append(userList, *user)
		size++
	}

	var err error = nil
	if size == 0 {
		err = errors.New("没有关注人信息")
	}

	return &userList, err
}

func (u *User) GetFollowerList(param *dto.FollowListInput) (*[]User, error) {
	db := u.conn()
	defer db.Close()

	followers := u.FollowerList //得到关注列表 字符串
	userIds := strings.Split(followers, "#")
	var userList []User

	size := 0
	for i := 0; i < len(userIds); i++ {
		x, _ := strconv.Atoi(userIds[i])
		id := uint(x)

		user, err := u.Search(db, id)
		if err != nil {
			continue
		}
		fmt.Println(user)
		userList = append(userList, *user)
		size++
	}

	var err error = nil
	if size == 0 {
		err = errors.New("没有粉丝信息")
	}

	return &userList, err
}
