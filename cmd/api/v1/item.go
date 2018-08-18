package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tetsuyanh/c2c-demo/model"
)

func handlerPostItem(c *gin.Context) {
	req := &model.Item{}
	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	i, err := itemSrv.CreateItem(c.GetString(requestUserID), req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, i)
}

func handlerPutItem(c *gin.Context) {
	req := &model.Item{}
	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	i, err := itemSrv.UpdateItem(c.Param("id"), req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, i)
}

func handlerDeleteItem(c *gin.Context) {
	if err := itemSrv.DeleteItem(c.Param("id")); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusOK)
}
