package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tetsuyanh/c2c-demo/model"
)

func handlerPostItem(c *gin.Context) {
	i := model.DefaultItem()
	if err := c.BindJSON(i); err != nil {
		b, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Println("req: ", string(b))
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	uId := c.GetString(requestUserID)
	t := time.Now()
	i.UserID = &uId
	i.CreatedAt = &t
	i.UpdatedAt = &t
	// TODO: validation
	userID := c.GetString(requestUserID)
	i.UserID = &userID
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
