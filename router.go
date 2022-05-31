package main

import (
	"gitee.com/Whitroom/imitate-tiktok/controller"
	"gitee.com/Whitroom/imitate-tiktok/middlewares"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// 不需要token校验的接口
	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.POST("/user/register/", controller.Register) // 完成
	apiRouter.POST("/user/login/", controller.Login)       // 完成

	// extra apis - I
	apiRouter.GET("/comment/list/", controller.CommentList)

	// 需要token校验的接口
	// basic apis
	apiRouter.POST("/publish/action/", controller.Publish)
	auth := apiRouter.Group("/", middlewares.AuthUser())
	auth.GET("/publish/list/", controller.PublishList)
	auth.GET("/user/", controller.UserInfo) // 完成

	// extra apis - I
	auth.POST("/favorite/action/", controller.FavoriteAction) // 完成
	auth.GET("/favorite/list/", controller.FavoriteList)      // 完成
	auth.POST("/comment/action/", controller.CommentAction)

	// extra apis - II
	auth.POST("/relation/action/", controller.RelationAction)     // 完成
	auth.GET("/relation/follow/list/", controller.FollowList)     // 完成
	auth.GET("/relation/follower/list/", controller.FollowerList) // 完成
}
