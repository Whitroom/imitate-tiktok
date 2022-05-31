package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionRequest struct {
	VideoID     uint   `form:"video_id" binding:"required"`
	ActionType  uint   `form:"action_type" binding:"required,min=1,max=2"`
	CommentText string `form:"comment_text" binding:"omitempty"`
	CommentID   uint   `form:"comment_id" binding:"omitempty"`
}

func CommentAction(ctx *gin.Context) {
	var request CommentActionRequest
	var comment *models.Comment
	err := ctx.ShouldBindQuery(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "参数绑定错误: " + err.Error(),
		})
		return
	}
	user_, _ := ctx.Get("User")
	user, _ := user_.(*models.User)
	if request.ActionType == 1 {
		comment = crud.CreateComment(sql.DB, &models.Comment{
			UserID:  user.ID,
			VideoID: request.VideoID,
			Content: request.CommentText,
		})
	} else {
		crud.DeleteComment(sql.DB, request.CommentID)
	}
	ctx.JSON(http.StatusOK, CommentModelChange(*comment))
}

func CommentList(ctx *gin.Context) {
	videoID, err := strconv.ParseUint(ctx.Query("video_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "参数绑定错误",
		})
	}
	comments := crud.GetComments(sql.DB, uint(videoID))
	fmt.Println(comments)
	ctx.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: CommentsModelChange(comments),
	})
}
