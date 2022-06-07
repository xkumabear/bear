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

type Follow struct {
	gorm.Model
	UserID       uint   `gorm:"ForeignKey:UserID;AssociationForeignKey:Id"`
	FollowList   string `gorm:"DEFAULT:''"`
	FollowerList string `gorm:"DEFAULT:''"`
	IsFollow     int64  `gorm:"DEFAULT:0"`
}

func (f *Follow) conn() *gorm.DB {
	db, err := gorm.Open(common.DRIVER, common.DSN)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Follow{})
	return db
}

func (f *Follow) Delete(db *gorm.DB) error {

	return db.Delete(f).Error
}
func (f *Follow) Find(db *gorm.DB, search *Follow) (*Follow, error) {
	//fmt.Println(search)
	var fo Follow
	err := db.Where(search).Find(&fo).Error
	if err != nil { //有错误没找到   返回   空  和 错误
		return search, err
	}
	return &fo, err // 找到了       返回   对象 和 空
}

// Add A of followList and B of followerList
func Add(fA *Follow, fB *Follow) error {
	userAIdString := strconv.Itoa(int(fA.Model.ID)) //获得 A 的id
	userBIdString := strconv.Itoa(int(fB.Model.ID)) //获得 B 的 ID

	followList := strings.Split(fA.FollowList, "#")
	//fmt.Println(followList, len(followList))
	followListA, errA := DeleteStringElement(followList, userBIdString)
	//fmt.Println(followListA, len(followListA))
	if errA == nil {
		followListA = append(followListA, userBIdString)
		fA.FollowList = strings.Join(followListA, "#")
		fA.FollowList += "#"
		fA.FollowList = "#" + fA.FollowList
		//fmt.Println(fA.FollowList, len(fA.FollowList))
	}

	fmt.Println("errA:", errA)
	followerList := strings.Split(fB.FollowerList, "#")
	followerListB, errB := DeleteStringElement(followerList, userAIdString)
	if errB == nil {
		followerListB = append(followerListB, userAIdString)
		fB.FollowerList = strings.Join(followerListB, "#")
		fB.FollowerList += "#"
		fB.FollowerList = "#" + fB.FollowerList
	}
	fmt.Println("errB:", errB)

	if errA != nil || errB != nil {
		fmt.Println("已添加")
		return errors.New("已添加")
	}

	sB := "#" + userBIdString + "#"
	isA := strings.Contains(fA.FollowerList, sB)
	isB := strings.Contains(fA.FollowList, sB)
	if isA && isB {
		fA.IsFollow = 1
		fB.IsFollow = 1
	} else {
		fA.IsFollow = 0
		fB.IsFollow = 0
	}

	return nil
}

func Delete(fA *Follow, fB *Follow) error {
	userAIdString := strconv.Itoa(int(fA.Model.ID)) //获得 A 的id
	userBIdString := strconv.Itoa(int(fB.Model.ID)) //获得 B 的 ID

	ss := "#" + userBIdString + "#"
	if is := strings.Contains(fA.FollowList, ss); !is {
		return errors.New("没有关注")
	}

	followList := strings.Split(fA.FollowList, "#")
	followListA, _ := DeleteStringElement(followList, userBIdString)
	fA.FollowList = strings.Join(followListA, "#")
	fA.FollowList += "#"
	fA.FollowList = "#" + fA.FollowList

	followerList := strings.Split(fB.FollowerList, "#")
	followerListB, _ := DeleteStringElement(followerList, userAIdString)
	fB.FollowerList = strings.Join(followerListB, "#")
	fB.FollowerList += "#"
	fB.FollowerList = "#" + fB.FollowerList

	sB := "#" + userBIdString + "#"

	isA := strings.Contains(fA.FollowerList, sB)
	isB := strings.Contains(fA.FollowList, sB)
	if isA && isB {
		fA.IsFollow = 1
		fB.IsFollow = 1
	} else {
		fA.IsFollow = 0
		fB.IsFollow = 0
	}
	return nil
}

func (f *Follow) RelationCheck(param *dto.RelationInput) (*User, error) {
	db := f.conn()
	defer db.Close()
	// 根据 id 查找用户 A B

	uA := &User{}                        //查找 A id  是否存在
	errA := uA.Search(db, param.UserAID) //查找 A ID
	fmt.Println("用户A有错？:", errA)
	if errA != nil { //没找到 A
		return uA, errA
	}
	uB := &User{}                        //在user中查找 B id  是否存在
	errB := uB.Search(db, param.UserBID) //, IsDelete: 0
	fmt.Println("用户 B 有错？:", errB)
	if errB != nil { //没找到 B
		return uA, errB
	}

	fA, errFa := f.Find(db, &Follow{UserID: param.UserAID}) //在follow查找 A的粉丝
	fmt.Println("粉丝表有fA，有错？：", errFa)
	if errFa != nil { // 有错没找到 fA  还没有关注 或 被关注 添加一条记录
		fmt.Println("fA不存在，创建fA")
		errFa = db.Create(fA).Error //
	}
	fmt.Println("粉丝表的fA:", fA)

	fB, errFb := f.Find(db, &Follow{UserID: param.UserBID}) //在follow查找 B的粉丝
	fmt.Println("粉丝表有fB，有错？:", errFb)
	if errFa != nil { // 有错没找到 fB  还没有关注 或 被关注 添加一条记录
		fmt.Println("fB不存在，创建fB")
		errFa = db.Create(fB).Error //添加后
	}
	fmt.Println("粉丝表的fB:", fB)

	if param.ActionType == "1" { // 关注  查找 A 是否 关注 B  查找 B 是否 关注 A
		fmt.Println("关注：")
		if errAdd := Add(fA, fB); errAdd == nil {
			fmt.Println("errAdd", errAdd)
			uA.FollowCount++
			uB.FollowerCount++
		}
	}

	if param.ActionType == "2" { //取消关注
		fmt.Println("取消关注：")
		if errDe := Delete(fA, fB); errDe == nil {
			uA.FollowCount--
			uB.FollowerCount--
		}
	}

	uA.FollowList = fA.FollowList
	uB.FollowerList = fB.FollowerList
	uA.IsFollow = fA.IsFollow
	uB.IsFollow = fB.IsFollow
	db.Save(uA)
	db.Save(uB)
	db.Save(fA)
	db.Save(fB)

	return uA, nil
}

func DeleteStringElement(list []string, ele string) ([]string, error) {
	result := make([]string, 0)
	var err error = nil
	for _, v := range list {
		if v == "" {
			continue
		}
		if v != ele {
			result = append(result, v)
		} else {
			err = errors.New("已存在")
		}

	}

	return result, err
}
