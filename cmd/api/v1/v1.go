package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tetsuyanh/c2c-demo/service"
)

const (
	headerSessionToken = "X-C2c-Session-Token"
	requestUserID      = "requestUserID"
)

var (
	userSrv service.UserService
	itemSrv service.ItemService
)

func Setup() error {
	userSrv = service.GetUserService()
	itemSrv = service.GetItemService()
	return nil
}

func Router(e *gin.Engine) {
	v1 := e.Group("/v1")
	{
		// public
		v1.POST("/sessions", handlerPostSession)
		v1.GET("/auths/enable", handlerGetAuthEnable)
		// v1.GET("/items", handlerGetItems)
		// v1.Get("/items/:id", handlerGetItemOne)

		// after here, require session to identify user
		v1.Use(setAuthenticatedUserID)
		v1.POST("/auths/publish", handlerPostAuthPublish)
		v1.POST("/auths/login", handlerPostAuthLogin)
		v1.POST("/items", handlerPostItem)
		v1.PUT("/items/:id", handlerPutItem)
		v1.DELETE("/items/:id", handlerDeleteItem)
	}
}

func setAuthenticatedUserID(c *gin.Context) {
	u, err := userSrv.GetUser(c.GetHeader(headerSessionToken))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	c.Set(requestUserID, u.ID)
}
