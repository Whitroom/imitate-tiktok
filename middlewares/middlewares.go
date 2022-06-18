package middlewares

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"gitee.com/Whitroom/imitate-tiktok/common/response"
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
		db := sql.GetSession()
		token := ctx.Query("token")
		userID, err := Parse(ctx, token)
		if err != nil {
			ctx.Abort()
			return
		}
		user, err := crud.GetUserByID(db, userID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, response.Response{
				StatusCode: response.NOTFOUND,
				StatusMsg:  "找不到相应用户",
			})
			ctx.Abort()
			return
		}
		ctx.Set("User", user)
		ctx.Next()
	}
}
