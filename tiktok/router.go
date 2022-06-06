package main

import (
	"github.com/gin-gonic/gin"
	"tiktok/controller"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.POST("/user/register/", controller.Register)  //用户注册
	apiRouter.POST("/user/login/", controller.Login)        //用户登陆
	apiRouter.GET("/user/", controller.UserInfo)            // 用户信息
	apiRouter.POST("/publish/action/", controller.Publish)  //发布信息
	apiRouter.GET("/publish/list/", controller.PublishList) // 发布列表
	apiRouter.GET("/feed/", controller.Feed)                //视频

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
}
