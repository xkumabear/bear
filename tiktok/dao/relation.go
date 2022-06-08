package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type relation struct {
	gorm.Model
	Id            int64
	Name          string
	Username      string
	Password      string
	FollowCount   int64
	FollowerCount int64
	FollowList    string
}
