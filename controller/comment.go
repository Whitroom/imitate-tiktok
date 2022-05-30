package controller

import (
	"net/http"
	"strconv"

	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

// // CommentAction no practical effect, just check if token is valid
// func CommentAction(c *gin.Context) {
// 	token := c.Query("token")

// 	if _, exist := usersLoginInfo[token]; exist {
// 		c.JSON(http.StatusOK, Response{StatusCode: 0})
// 	} else {
// 		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
// 	}
// }

// // CommentList all videos have same demo comment list
// func CommentList(c *gin.Context) {
// 	c.JSON(http.StatusOK, CommentListResponse{
// 		Response:    Response{StatusCode: 0},
// 		CommentList: DemoComments,
// 	})
// }

type commentlistactionrequest struct {
	UserID     uint `form:"user_id" binding:"required"`
	VideoID    uint `form:"video_id" binding:"required"`
	ActionType uint `form:"action_type" binding:"required,min=1,max=2"`
}

func CommentsModelChange(Comments []models.Comment) []Comment {
	var comments []Comment
	for _, comment := range Comments {
		comments = append(comments, CommentModelChange(comment))
	}
	return comments
}

func CommentAction(c *gin.Context) {
	var request commentlistactionrequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "bind failed",
		})
		return
	}
	//通过videoid查询视频信息
	var modelsvideo models.Video
	if err := sql.DB.Where("id=?", request.VideoID).First(&modelsvideo).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	videos := VideoModelChange(&modelsvideo)
	if request.ActionType == 1 {
		if videos.CommentCount == 0 {
			c.JSON(http.StatusBadRequest, Response{
				StatusCode: 1,
				StatusMsg:  "the commenttext is empty",
			})
			return
		} else {
			//添加评论记录
		}
	} else {
		var modelscomment models.Comment
		_ = sql.DB.First(&modelscomment, request.UserID).Error
		comment := CommentModelChange(modelscomment)
		if comment.Content == "" {
			c.JSON(http.StatusNotFound, Response{
				StatusCode: 1,
				StatusMsg:  "comments doesn't exist",
			})
			return
		} else {
			//删除评论
		}
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "success",
	})
}

func CommentList(c *gin.Context) {
	userid := c.Query("user_id")
	if userid == "" {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exixt",
		})
		return
	}
	user_id, _ := strconv.ParseInt(userid, 0, 0)
	var commentlist []models.Comment
	_ = sql.DB.Find(&commentlist, "user_id=?", user_id).Error
	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: 0,
		},
		CommentList: CommentsModelChange(commentlist),
	})
}
