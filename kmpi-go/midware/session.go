package midware

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"kmpi-go/log"

	"net/http"
	"time"
)

func SessionMidware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIp := c.ClientIP()

		session := sessions.Default(c)
		loginToken := session.Get("loginToken")
		if loginToken == nil {
			session.Set("loginToken", "login"+time.Now().String())
			err := session.Save()
			if err != nil {
				fmt.Println(err.Error())

			}
		} else {
			log.Info("loginToken:", loginToken)
		}
		if clientIp != "192.168.0.255" {
			// 验证通过，会继续访问下一个中间件
			c.Next()
		} else {
			// 验证不通过，不再调用后续的函数处理
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"message": "访问未授权"})
			// return可省略, 只要前面执行Abort()就可以让后面的handler函数不再执行
			return
		}
	}
}
