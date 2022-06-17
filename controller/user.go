package controller

import (
	"fmt"
	"net/http"

	"gitee.com/Whitroom/imitate-tiktok/common"
	"gitee.com/Whitroom/imitate-tiktok/common/response"
	"gitee.com/Whitroom/imitate-tiktok/middlewares"
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required" form:"username"`
	Password string `json:"password" binding:"required" form:"password"`
}

func hashEncode(str string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("failed to hash:%w", err)
	}
	return string(hash)
}

func comparePasswords(sourcePwd, hashPwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(sourcePwd)) == nil
}

func Register(ctx *gin.Context) {
	db := sql.GetSession()

	var request RegisterRequest
	if !common.BindAndValid(ctx, &request) {
		return
	}
	if crud.GetUserByName(db, request.Username) != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 2,
			StatusMsg:  "存在用户姓名",
		})
		return
	}
	newUser := crud.CreateUser(db, &models.User{
		Name:     request.Username,
		Password: hashEncode(request.Password),
		Content:  "",
	})

	token, err := middlewares.Sign(newUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			StatusCode: 3,
			StatusMsg:  "token创建失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.UserLoginResponse{
		Response: response.Response{
			StatusCode: 0,
			StatusMsg:  "用户创建成功",
		},
		UserID: int64(newUser.ID),
		Token:  token,
	})
}

func Login(ctx *gin.Context) {
	db := sql.GetSession()

	var request RegisterRequest
	if !common.BindAndValid(ctx, &request) {
		return
	}
	existedUser := crud.GetUserByName(db, request.Username)

	if existedUser == nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			StatusCode: 2,
			StatusMsg:  "找不到用户",
		})
		return
	}

	if !comparePasswords(request.Password, existedUser.Password) {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			StatusCode: 3,
			StatusMsg:  "用户名或密码错误",
		})
		return
	}

	token, err := middlewares.Sign(existedUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			StatusCode: 4,
			StatusMsg:  "token创建失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.UserLoginResponse{
		Response: response.Response{StatusCode: 0},
		UserID:   int64(existedUser.ID),
		Token:    token,
	})
}

// 查询用户信息接口函数。
func UserInfo(ctx *gin.Context) {
	db := sql.GetSession()

	var user *models.User

	token := ctx.Query("token")
	if token != "" {
		userID, err := middlewares.Parse(ctx, token)
		if err != nil {
			return
		}
		user, err = crud.GetUserByID(db, userID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, response.Response{
				StatusCode: 3,
				StatusMsg:  "找不到用户",
			})
			return
		}
	}

	toUserID := common.QueryIDAndValid(ctx, "user_id")
	if toUserID == 0 {
		return
	}

	toUser, err := crud.GetUserByID(db, toUserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			StatusCode: 3,
			StatusMsg:  "找不到用户",
		})
		return
	}

	responseUser := common.UserModelChange(db, *toUser)
	if user != nil {
		responseUser.IsFollow = crud.IsUserFollow(db, user.ID, toUserID)
	} else {
		responseUser.IsFollow = false
	}

	ctx.JSON(http.StatusOK, response.UserResponse{
		Response: response.Response{
			StatusCode: 0,
			StatusMsg:  "已找到用户",
		},
		User: responseUser,
	})

}
