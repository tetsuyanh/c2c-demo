package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/tetsuyanh/c2c-demo/repository"
	"github.com/tetsuyanh/c2c-demo/service"
)

const (
	headerSessionToken = "X-C2c-Session-Token"
	requestUserID      = "requestUserID"
	requesSelectOption = "requestSelectOption"
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
		v1.Use(setAuthenticatedUserID)
		v1.POST("/auths/publish", handlerPostAuthPublish)
		v1.POST("/auths/login", handlerPostAuthLogin)
		v1.POST("/items", handlerPostItem)
		v1.PUT("/items/:id", handlerPutItem)
		v1.DELETE("/items/:id", handlerDeleteItem)
		v1.GET("/deals/seller", setSelectOption, handlerGetDealAsSeller)
		v1.GET("/deals/buyer", setSelectOption, handlerGetDealAsBuyer)
		v1.POST("/deals", handlerPostDeal)
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

func setSelectOption(c *gin.Context) {
	opt := repository.DefaultOption()

	// keys
	opt.SetUserId(c.GetString(requestUserID))

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

	c.Set(requesSelectOption, opt)
}
