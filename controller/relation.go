package controller

import (
	"net/http"
	"strconv"

	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

type RelationActionRequest struct {
	UserID     uint `binding:"required" form:"user_id"`
	ToUserID   uint `binding:"required" form:"to_user_id"`
	ActionType uint `binding:"required,min=1,max=2" form:"action_type"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(ctx *gin.Context) {
	var request RelationActionRequest
	ctx.ShouldBindQuery(&request)
	if request.ActionType == 1 {
		_, err := crud.SubscribeUser(sql.DB, request.UserID, request.ToUserID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
	} else {
		_, err := crud.CancelSubscribeUser(sql.DB, request.UserID, request.ToUserID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, Response{
				StatusCode: 2,
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

// FollowList all users have same follow list
func FollowList(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Query("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "user_id不是数字",
		})
		return
	}
	users := crud.GetUserSubscribersByID(sql.DB, uint(userID))
	ctx.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: UsersModelChange(users),
	})
}

// FollowerList all users have same follower list
func FollowerList(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Query("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "user_id不是数字",
		})
		return
	}
	users := crud.GetUserFollowersByID(sql.DB, uint(userID))
	modelUsers := UsersModelChange(users)
	for _, modelUser := range modelUsers {
		modelUser.IsFollow = crud.IsUserFollow(
			sql.DB, uint(userID), uint(modelUser.Id),
		)
	}
	ctx.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: modelUsers,
	})
}
