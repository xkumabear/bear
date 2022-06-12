package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"strings"
	"tiktok/common"
)

type Follow struct {
	gorm.Model
	UserID       uint   `gorm:"ForeignKey:UserID;AssociationForeignKey:Id"`
	FollowList   string `gorm:"DEFAULT:''"`
	FollowerList string `gorm:"DEFAULT:''"`
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

	return nil
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
