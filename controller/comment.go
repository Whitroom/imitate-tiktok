package controller

import (
	"net/http"

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
	if BindAndValid(ctx, &request) {
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
	videoID := QueryIDAndValid(ctx, "video_id")
	if videoID == 0 {
		return
	}

	comments := crud.GetComments(sql.DB, videoID)

	ctx.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: CommentsModelChange(comments),
	})
}
