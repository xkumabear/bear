package dto

import (
	"github.com/gin-gonic/gin"
)

type RelationInput struct {
	Token      string `json:"token" form:"token" `
	UserBID    uint   `json:"to_user_id" form:"to_user_id" `
	ActionType string `json:"action_type" form:"action_type"`
}

func (params *RelationInput) GetValidParams(c *gin.Context) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	return nil
}

type FollowListInput struct {
	UserID uint   `json:"user_id" form:"user_id"`
	Token  string `json:"token" form:"token" `
}

func (params *FollowListInput) GetValidParams(c *gin.Context) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	return nil
}

type FollowOutput struct {
	Response
	UserList []User `json:"user_list"`
}
