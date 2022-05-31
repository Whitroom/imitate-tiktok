package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"gitee.com/Whitroom/imitate-tiktok/middlewares"
	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func hashEncode(str string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("failed to hash:%w", err)
	}
	return string(hash)
}

func comparePasswords(sourcePwd, hashPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(sourcePwd))
	return err == nil
}

var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required" form:"username"`
	Password string `json:"password" binding:"required" form:"password"`
}

func Register(ctx *gin.Context) {
	var user RegisterRequest
	if err := ctx.ShouldBindQuery(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "绑定失败",
		})
		return
	}

	if crud.GetUserByName(sql.DB, user.Username) == nil {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 2,
			StatusMsg:  "存在用户姓名",
		})
		return
	}

	newUser := &models.User{
		Name:     user.Username,
		Password: hashEncode(user.Password),
		Content:  "",
	}

	newUser = crud.CreateUser(sql.DB, newUser)

	token, err := middlewares.Sign(newUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			StatusCode: 3,
			StatusMsg:  "token创建失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "用户创建成功",
		},
		UserId: int64(newUser.ID),
		Token:  token,
	})
}

func Login(ctx *gin.Context) {
	var user RegisterRequest
	if err := ctx.ShouldBindQuery(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "绑定失败",
			},
		})
		return
	}
	existedUser := crud.GetUserByName(sql.DB, user.Username)

	if existedUser == nil {
		ctx.JSON(http.StatusNotFound, Response{
			StatusCode: 1,
			StatusMsg:  "找不到用户",
		})
		return
	}

	pwdMatch := comparePasswords(user.Password, existedUser.Password)
	if !pwdMatch {
		ctx.JSON(http.StatusUnauthorized, UserLoginResponse{
			Response: Response{
				StatusCode: 2,
				StatusMsg:  "用户名或密码错误",
			},
		})
		return
	}

	token, err := middlewares.Sign(existedUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			StatusCode: 3,
			StatusMsg:  "token创建失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   int64(existedUser.ID),
		Token:    token,
	})
}

// 查询用户信息接口函数。
func UserInfo(ctx *gin.Context) {

	user_, _ := ctx.Get("User")
	user, _ := user_.(*models.User)

	toUserID_, err := strconv.ParseUint(ctx.Query("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 2,
			StatusMsg:  "user_id错误" + err.Error(),
		})
		return
	}

	toUserID := uint(toUserID_)

	toUser, err := crud.GetUserByID(sql.DB, toUserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			StatusCode: 3,
			StatusMsg:  "找不到用户",
		})
		return
	}

	responseUser := UserPointerModelChange(toUser)
	if user != nil {
		responseUser.IsFollow = crud.IsUserFollow(sql.DB, user.ID, toUserID)
	} else {
		responseUser.IsFollow = false
	}

	ctx.JSON(http.StatusOK, UserResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "已找到用户",
		},
		User: responseUser,
	})

}
