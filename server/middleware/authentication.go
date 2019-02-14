package middleware

import (
	"go-eladmin/common"
	"go-eladmin/core/jwt"
	"go-eladmin/model"
	"go-eladmin/model/model_user"
	"go-eladmin/server/sync_map"
	"net/http"

	"github.com/gin-gonic/gin"
)

//JWTAuthMiddleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		aRes := model.NewResponse()
		//header拿token
		token := c.GetHeader(common.KeyAuthorization)
		if token == "" {
			//cookie拿token
			token, _ = c.Cookie(common.KeyAuthorization)
			if token == "" {
				aRes.SetErrorInfo(http.StatusUnauthorized, "请求未携带token，无权限访问")
				c.JSON(http.StatusUnauthorized, aRes)
				c.Abort()
				return
			}
		}
		claims, err := jwt.ParseToken(token)
		if err != nil {
			sync_map.RemoveKey(token)
			aRes.SetErrorInfo(http.StatusUnauthorized, err.Error())
			c.JSON(http.StatusUnauthorized, aRes)
			c.Abort()
			return
		}
		if !sync_map.Value(token) {
			_, err = model_user.FindLoginByToken(token)
			if err != nil {
				aRes.SetErrorInfo(http.StatusUnauthorized, "token无效，无权限访问")
				c.JSON(http.StatusUnauthorized, aRes)
				c.Abort()
				return
			}
			sync_map.SetKeyValue(token)
		}

		c.Set(common.KeyContextUserId, claims.UserId)
		c.Next()
	}
}
