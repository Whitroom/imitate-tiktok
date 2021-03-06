package main

import (
	"gitee.com/Whitroom/imitate-tiktok/controller"
	"gitee.com/Whitroom/imitate-tiktok/middlewares"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// 不需要token校验的接口
	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.GET("/user/", controller.UserInfo)

	// extra apis - I
	apiRouter.GET("/comment/list/", controller.CommentList)

	// 需要token校验的接口
	// basic apis
	apiRouter.POST("/publish/action/", controller.Publish)

	auth := apiRouter.Group("/", middlewares.AuthUser())
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
