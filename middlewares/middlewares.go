package middlewares

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/crud"
	"github.com/gin-gonic/gin"
)

var Secret []byte

func InitSecret() {
	conf, _ := os.Open("./confs/secret.json")
	defer conf.Close()
	value, _ := ioutil.ReadAll(conf)
	json.Unmarshal([]byte(value), &map[string][]byte{"secret": Secret})
}

// 验证用户中间件, 若没有token会返回400, 验证失败会返回401, 找不到用户会返回404, 响应code为1, 2, 3
func AuthUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		if token == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"StatusCode": 1,
				"StatusMsg":  "没有相应的token, 请重新登陆获取",
			})
			ctx.Abort()
			return
		}
		UserID, err := Parse(ctx, token)
		if err != nil {
			ctx.Abort()
			return
		}
		User, err := crud.GetUserByID(sql.DB, UserID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"StatusCode": 3,
				"StatusMsg":  "token解析错误, 请重新登陆获取",
			})
			ctx.Abort()
			return
		}
		ctx.Set("User", User)
		ctx.Next()
	}
}
