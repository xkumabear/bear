package dao

import (
	"crypto/sha256"
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
	saltPassword := GetSaltString(common.Salt, param.Password)
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
	saltPassword := GetSaltString(common.Salt, param.Password)
	if user.Password != saltPassword {
		return nil, errors.New("密码错误！")
	}
	return user, nil
}

func GetSaltString(salt, password string) string {
	s1 := sha256.New()
	s1.Write([]byte(password))
	str1 := fmt.Sprintf("%x", s1.Sum(nil))
	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))
	return fmt.Sprintf("%x", s2.Sum(nil))
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

	user, err := u.Search(db, param.UserID)
	if err != nil {
		return nil, err
	}

	follows := user.FollowList //得到关注列表 字符串
	users := strings.Split(follows, "#")
	var userList []User

	for _, item := range users {
		x, _ := strconv.Atoi(item)
		userInfo, err := u.Search(db, uint(x))
		if err != nil {
			continue
		}
		userList = append(userList, *userInfo)
	}

	if len(users) == 0 {
		err = errors.New("没有关注人信息")
	}

	return &userList, err
}

func (u *User) GetFollowerList(param *dto.FollowListInput) (*[]User, error) {
	db := u.conn()
	defer db.Close()

	user, err := u.Search(db, param.UserID)
	if err != nil {
		return nil, err
	}

	followers := user.FollowerList //得到关注列表 字符串
	users := strings.Split(followers, "#")
	var userList []User

	for _, item := range users {
		x, _ := strconv.Atoi(item)
		userInfo, err := u.Search(db, uint(x))
		if err != nil {
			continue
		}
		userList = append(userList, *userInfo)
	}

	if len(users) == 0 {
		err = errors.New("没有关注人信息")
	}

	return &userList, err
}

func (u *User) RelationCheck(userid uint, param *dto.RelationInput) error {
	db := u.conn()
	defer db.Close()
	// 根据 id 查找用户 A B
	//var u *User
	var userA User
	err := db.Model(userA).Where("id = ?", userid).Find(&userA).Error
	if err != nil {
		fmt.Println("userA:", err)
		return err
	}

	var userB User
	err = db.Model(userB).Where("id = ? ", param.UserBID).Find(&userB).Error
	if err != nil {
		fmt.Println("userB:", err)
		return err
	}

	if param.ActionType == "1" { // 关注  查找 A 是否 关注 B  查找 B 是否 关注 A
		fmt.Println("get action:")
		userBid := strconv.FormatInt(int64(param.UserBID), 10) + "#"
		isExist := strings.Contains(userA.FollowList, userBid)
		if isExist {
			fmt.Println("have action!")
			return nil
			//return errors.New("relation is exist!")
		}
		userA.FollowList += userBid
		userA.FollowCount++
		userAid := strconv.FormatInt(int64(userid), 10) + "#"
		userB.FollowerList += userAid
		userB.FollowerCount++
		db.Save(userA)
		db.Save(userB)

	}

	if param.ActionType == "2" { //取消关注
		fmt.Println("don't have action:")
		userBid := strconv.FormatInt(int64(param.UserBID), 10) + "#"
		isExist := strings.Contains(userA.FollowList, userBid)
		if !isExist {
			fmt.Println("have not action!")
			return nil
			//return errors.New("relation is not exist!")
		}
		userA.FollowList = strings.Replace(userA.FollowList, userBid, "", -1)
		userA.FollowCount--
		userAid := strconv.FormatInt(int64(userid), 10) + "#"
		userB.FollowerList = strings.Replace(userB.FollowerList, userAid, "", -1)
		userB.FollowerCount--
		db.Save(userA)
		db.Save(userB)
	}

	return nil

}
