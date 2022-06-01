package controller

import (
	"net/http"

	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

func CommentAction(ctx *gin.Context) {
	var request struct {
		VideoID     uint   `form:"video_id" binding:"required"`
		ActionType  uint   `form:"action_type" binding:"required,min=1,max=2"`
		CommentText string `form:"comment_text" binding:"omitempty"`
		CommentID   uint   `form:"comment_id" binding:"omitempty"`
	}

	if !BindAndValid(ctx, &request) {
		return
	}
	user := GetUserFromCtx(ctx)
	if request.ActionType == 1 {
		comment := crud.CreateComment(&models.Comment{
			UserID:  user.ID,
			VideoID: request.VideoID,
			Content: request.CommentText,
		})
		ctx.JSON(http.StatusOK, CommentResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "添加评论成功",
			},
			Comment: CommentModelChange(*comment),
		})
	} else {
		crud.DeleteComment(request.CommentID)
		ctx.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "评论删除成功",
		})
	}

}

func CommentList(ctx *gin.Context) {
	videoID := QueryIDAndValid(ctx, "video_id")
	if videoID == 0 {
		return
	}

	ctx.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: CommentsModelChange(crud.GetComments(videoID)),
	})
}
