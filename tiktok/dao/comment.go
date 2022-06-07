package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type follow struct {
	gorm.Model
	Id       int64
	UserAID  int64 `gorm:"ForeignKey:UserID;AssociationForeignKey:Id"`
	UserBID  int64
	User     User
	Relation int
}

type follower struct {
	gorm.Model
	Id       int64
	UserAID  int64
	UserBID  int64 `gorm:"ForeignKey:UserID;AssociationForeignKey:Id"`
	User     User
	Relation int
}
