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

type CommentResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

func CommentAction(ctx *gin.Context) {
	db := sql.GetDB()

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
		if len(request.CommentText) == 0 {
			ctx.JSON(http.StatusBadRequest, Response{
				StatusCode: 1,
				StatusMsg:  "评论文本为空",
			})
			return
		}
		comment := crud.CreateComment(&db, &models.Comment{
			UserID:  user.ID,
			VideoID: request.VideoID,
			Content: request.CommentText,
		})
		ctx.JSON(http.StatusOK, CommentResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "添加评论成功",
			},
			Comment: CommentModelChange(&db, *comment),
		})
	} else {
		if request.CommentID == 0 {
			ctx.JSON(http.StatusBadRequest, Response{
				StatusCode: 1,
				StatusMsg:  "删除失败",
			})
			return
		}
		crud.DeleteComment(&db, request.CommentID)
		ctx.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "评论删除成功",
		})
	}

}

func CommentList(ctx *gin.Context) {
	db := sql.GetDB()

	videoID := QueryIDAndValid(ctx, "video_id")
	if videoID == 0 {
		return
	}

	ctx.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: CommentsModelChange(&db, crud.GetComments(&db, videoID)),
	})
}
