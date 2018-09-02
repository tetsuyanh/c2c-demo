package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/tetsuyanh/c2c-demo/repository"
	"github.com/tetsuyanh/c2c-demo/service"
)

const (
	headerSessionToken = "X-C2c-Session-Token"
	sessionUserId      = "sessioUserId"
	authedUserId       = "authedUserId"
	selectOption       = "selectOption"
)

var (
	userSrv service.UserService
	itemSrv service.ItemService
	dealSrv service.DealService
)

func Setup() error {
	userSrv = service.GetUserService()
	itemSrv = service.GetItemService()
	dealSrv = service.GetDealService()
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
		v1.Use(setSessionUserId)
		v1.POST("/auths/publish", handlerPostAuthPublish)
		v1.POST("/auths/login", handlerPostAuthLogin)

		// after here, require session of authenticated user
		v1.Use(setAuthenticatedUserId)
		v1.POST("/items", handlerPostItem)
		v1.GET("/items", setSelectOption, handlerGetItems)
		v1.GET("/item/:id", handlerGetItem)
		v1.PUT("/items/:id", handlerPutItem)
		v1.DELETE("/items/:id", handlerDeleteItem)
		v1.GET("/deals/seller", setSelectOption, handlerGetDealAsSeller)
		v1.GET("/deals/buyer", setSelectOption, handlerGetDealAsBuyer)
		v1.POST("/deals", handlerPostDeal)
	}
}

func setSessionUserId(c *gin.Context) {
	u, err := userSrv.GetUser(c.GetHeader(headerSessionToken))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	c.Set(sessionUserId, u.Id)
}

// expect called after setSessionUserId, use sessionUserId
func setAuthenticatedUserId(c *gin.Context) {
	a, err := userSrv.GetAuthentication(c.GetString(sessionUserId))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	if !a.Enabled {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("not enabled"))
	}
	c.Set(authedUserId, a.UserId)
}

func setSelectOption(c *gin.Context) {
	opt := repository.DefaultOption()

	// keys
	opt.SetUserId(c.GetString(authedUserId))

	// conditions
	limit := c.Query("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		opt.SetLimit(l)
	}
	offset := c.Query("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		opt.SetOffset(o)
	}

	c.Set(selectOption, opt)
}
