package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"gitee.com/Whitroom/imitate-tiktok/common/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// 签发token, 仅传入userID, token令牌中仅包含id和创建时间
func Sign(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"nbf": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(Secret))
	return tokenString, err
}

// 解析token, 仅传入token字符串, 仅判断是否转换成功与令牌是否超时
func Parse(ctx *gin.Context, tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("加密方式错误: %v", t.Header["alg"])
		}
		return []byte(Secret), nil
	})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			StatusCode: response.BADREQUEST,
			StatusMsg:  "token获取错误, 请重新登陆获取",
		})
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			StatusCode: response.BADREQUEST,
			StatusMsg:  "token解析错误, 请重新登陆获取",
		})
		return 0, err
	}
	userID := uint(claims["id"].(float64))
	createTime, _ := claims["nbf"].(int64)

	if time.Now().Unix()-createTime > 2*int64(time.Hour) {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			StatusCode: response.BADREQUEST,
			StatusMsg:  "token超时, 请重新登陆获取",
		})
		return 0, err
	}

	return userID, nil
}
