package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tetsuyanh/c2c-demo/model"
	"github.com/tetsuyanh/c2c-demo/repository"
)

func handlerPostItem(c *gin.Context) {
	req := &model.Item{}
	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	i, err := itemSrv.CreateItem(c.GetString(authedUserId), req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, i)
}

func handlerGetPublicItems(c *gin.Context) {
	obj, _ := c.Get(selectOption)
	is, err := itemSrv.GetPublicItems(obj.(*repository.Option))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, is)
}

func handlerGetPublicItem(c *gin.Context) {
	i, err := itemSrv.GetPublicItem(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, i)
}

func handlerGetMyItems(c *gin.Context) {
	obj, _ := c.Get(selectOption)
	is, err := itemSrv.GetMyItems(obj.(*repository.Option))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, is)
}

func handlerGetMyItem(c *gin.Context) {
	i, err := itemSrv.GetMyItem(c.Param("id"), c.GetString(authedUserId))
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
	i, err := itemSrv.UpdateItem(c.Param("id"), c.GetString(authedUserId), req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, i)
}

func handlerDeleteItem(c *gin.Context) {
	if err := itemSrv.DeleteItem(c.Param("id"), c.GetString(authedUserId)); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusOK)
}
