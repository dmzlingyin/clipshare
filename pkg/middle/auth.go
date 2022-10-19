package middle

import "github.com/gin-gonic/gin"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func Login(username, password string) error {
	return nil
}

func Register(username, password string) error {
	return nil
}
