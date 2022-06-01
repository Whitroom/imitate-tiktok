package controller

import (
	"net/http"

	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

type RelationActionRequest struct {
	ToUserID   uint `binding:"required" form:"to_user_id"`
	ActionType uint `binding:"required,min=1,max=2" form:"action_type"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(ctx *gin.Context) {
	var request RelationActionRequest
	if !BindAndValid(ctx, &request) {
		return
	}
	user := GetUserFromCtx(ctx)
	if request.ToUserID == user.ID {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 2,
			StatusMsg:  "无法关注或取关自己",
		})
		return
	}
	if request.ActionType == 1 {
		_, err := crud.SubscribeUser(user.ID, request.ToUserID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, Response{
				StatusCode: 3,
				StatusMsg:  err.Error(),
			})
			return
		}
	} else {
		_, err := crud.CancelSubscribeUser(user.ID, request.ToUserID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, Response{
				StatusCode: 4,
				StatusMsg:  err.Error(),
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "操作成功",
	})
}

func FollowList(ctx *gin.Context) {
	userID := QueryIDAndValid(ctx, "user_id")
	if userID == 0 {
		return
	}
	users := crud.GetUserSubscribersByID(userID)
	responseUsers := UsersModelChange(users)
	ctx.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: responseUsers,
	})
}

// FollowerList all users have same follower list
func FollowerList(ctx *gin.Context) {
	userID := QueryIDAndValid(ctx, "user_id")
	if userID == 0 {
		return
	}
	users := crud.GetUserFollowersByID(uint(userID))
	responseUsers := UsersModelChange(users)
	for i := 0; i < len(responseUsers); i++ {
		responseUsers[i].IsFollow = crud.IsUserFollow(
			uint(responseUsers[i].ID), uint(userID),
		)
	}
	ctx.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: responseUsers,
	})
}
