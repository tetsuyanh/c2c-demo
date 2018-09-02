package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tetsuyanh/c2c-demo/model"
	"github.com/tetsuyanh/c2c-demo/repository"
)

func handlerGetDealAsSeller(c *gin.Context) {
	obj, _ := c.Get(selectOption)
	ds, err := dealSrv.GetDealAsSeller(obj.(*repository.Option))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, ds)
}

func handlerGetDealAsBuyer(c *gin.Context) {
	obj, _ := c.Get(selectOption)
	ds, err := dealSrv.GetDealAsBuyer(obj.(*repository.Option))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, ds)
}

func handlerPostDeal(c *gin.Context) {
	req := model.Deal{}
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	d, err := dealSrv.Establish(req.ItemId, c.GetString(authedUserId))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, d)
}
