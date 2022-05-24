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
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)

	// extra apis - I
	apiRouter.GET("/comment/list/", controller.CommentList)

	// 需要token校验的接口
	// basic apis
	auth := apiRouter.Group("/", middlewares.AuthUser())
	auth.POST("/publish/action/", controller.Publish)
	auth.GET("/publish/list/", controller.PublishList)

	// extra apis - I
	auth.POST("/favorite/action/", controller.FavoriteAction)
	auth.GET("/favorite/list/", controller.FavoriteList)
	auth.POST("/comment/action/", controller.CommentAction)

	// extra apis - II
	auth.POST("/relation/action/", controller.RelationAction)
	auth.GET("/relation/follow/list/", controller.FollowList)
	auth.GET("/relation/follower/list/", controller.FollowerList)
}
