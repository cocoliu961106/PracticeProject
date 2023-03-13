package middlewares

import (
	"PracticeProject/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthAdminCheck
// 验证用户是不是普通用户
func AuthUserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的token并解析，判断isAdmin字段是否为1
		auth := c.GetHeader("Authorization")
		userClaim, err := helper.AnalyseToken(auth)
		if err != nil {
			c.Abort() // 阻止后续中间件执行
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}
		if userClaim == nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized User",
			})
			return
		}
		c.Set("user", userClaim)
		c.Next() // 马上执行下一个中间件，挂起c.NEXT()后的代码，待后面的中间件执行完后再执行
	}
}
