package controller

import (
	"net/http"

	"gitee.com/Whitroom/imitate-tiktok/common"
	"gitee.com/Whitroom/imitate-tiktok/common/response"
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"github.com/gin-gonic/gin"
)

func RelationAction(ctx *gin.Context) {
	db := sql.GetSession()

	var request struct {
		ToUserID   uint `binding:"required" form:"to_user_id"`
		ActionType uint `binding:"required,min=1,max=2" form:"action_type"`
	}
	if !common.BindAndValid(ctx, &request) {
		return
	}
	user := common.GetUserFromCtx(ctx)
	if request.ToUserID == user.ID {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 2,
			StatusMsg:  "无法关注或取关自己",
		})
		return
	}
	if request.ActionType == 1 {
		_, err := crud.SubscribeUser(db, user.ID, request.ToUserID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, response.Response{
				StatusCode: 3,
				StatusMsg:  err.Error(),
			})
			return
		}
		if crud.IsUserFollow(db, user.ID, request.ToUserID) {
			ctx.JSON(http.StatusBadRequest, response.Response{
				StatusCode: 4,
				StatusMsg:  "已关注过用户",
			})
			return
		}
	} else {
		_, err := crud.CancelSubscribeUser(db, user.ID, request.ToUserID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, response.Response{
				StatusCode: 3,
				StatusMsg:  err.Error(),
			})
			return
		}
		if !crud.IsUserFollow(db, user.ID, request.ToUserID) {
			ctx.JSON(http.StatusBadRequest, response.Response{
				StatusCode: 4,
				StatusMsg:  "未关注过用户",
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, response.Response{
		StatusCode: 0,
		StatusMsg:  "操作成功",
	})
}

func FollowList(ctx *gin.Context) {
	db := sql.GetSession()

	userID := common.QueryIDAndValid(ctx, "user_id")
	if userID == 0 {
		return
	}
	users := crud.GetUserSubscribersByID(db, userID)
	responseUsers := common.UsersModelChange(db, users)
	ctx.JSON(http.StatusOK, response.UserListResponse{
		Response: response.Response{
			StatusCode: 0,
		},
		UserList: responseUsers,
	})
}

func FollowerList(ctx *gin.Context) {
	db := sql.GetSession()

	userID := common.QueryIDAndValid(ctx, "user_id")
	if userID == 0 {
		return
	}
	users := crud.GetUserFollowersByID(db, uint(userID))
	responseUsers := common.UsersModelChange(db, users)
	for i := 0; i < len(responseUsers); i++ {
		responseUsers[i].IsFollow = crud.IsUserFollow(
			db, uint(responseUsers[i].ID), uint(userID),
		)
	}
	ctx.JSON(http.StatusOK, response.UserListResponse{
		Response: response.Response{
			StatusCode: 0,
		},
		UserList: responseUsers,
	})
}
