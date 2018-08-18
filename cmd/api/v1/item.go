package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tetsuyanh/c2c-demo/model"
)

func handlerPostItem(c *gin.Context) {
	req := &model.Item{}
	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	uId := c.GetString(requestUserID)
	i := model.DefaultItem()
	i.UserId = &uId
	i.Label = req.Label
	i.Description = req.Description
	i.Price = req.Price
	if err := itemSrv.CreateItem(i); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, i)
}

func handlerPutItem(c *gin.Context) {
	i, err := itemSrv.GetItem(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err := c.BindJSON(i); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	t := time.Now()
	i.UpdatedAt = &t
	// TODO: validation
	if err := itemSrv.UpdateItem(i); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, i)
}

func handlerDeleteItem(c *gin.Context) {
	i, err := itemSrv.GetItem(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err := itemSrv.DeleteItem(i); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
