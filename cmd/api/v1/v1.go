package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tetsuyanh/c2c-demo/service"
)

const (
	HeaderSessionToken  = "X-C2c-Session-Token"
	AuthenticatedUserID = "AuthenticatedUserID"
)

var (
	userSrv service.UserService
)

func Setup() error {
	userSrv = service.GetUserService()
	return nil
}

func Router(e *gin.Engine) {
	v1 := e.Group("/v1")
	{
		// public
		v1.POST("/sessions", handlerPostSession)
		v1.GET("/auths/enable", handlerGetAuthEnable)

		// after here, require session to identify user
		v1.Use(setAuthenticatedUserID)
		v1.POST("/auths/publish", handlerPostAuthPublish)
		v1.POST("/auths/login", handlerPostAuthLogin)
	}
}

func setAuthenticatedUserID(c *gin.Context) {
	u, err := userSrv.GetUser(c.GetHeader(HeaderSessionToken))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	c.Set(AuthenticatedUserID, u.ID)
}
