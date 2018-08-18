package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tetsuyanh/c2c-demo/model"
)

func handlerPostSession(c *gin.Context) {
	s, err := userSrv.CreateUserSession()
	if err != nil {
		log.Printf("userSrv.CreateNewUserSession: %v\n", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, s)
}

func handlerGetAuthEnable(c *gin.Context) {
	token, _ := c.GetQuery("token")
	if err := userSrv.EnableAuth(token); err != nil {
		log.Printf("userSrv.EnableAuth: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

func handlerPostAuthPublish(c *gin.Context) {
	req := model.Authentication{}
	if err := c.BindJSON(&req); err != nil {
		log.Printf("c.BindJSON: %v\n", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	a, err := userSrv.CreateAuth(c.GetString(requestUserID), *req.EMail, *req.Password)
	if err != nil {
		log.Printf("userSrv.CreateAuth: %v\n", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// TODO: should hide token in response and tell it by sending email for production
	c.JSON(http.StatusOK, a)
}

func handlerPostAuthLogin(c *gin.Context) {
	req := model.Authentication{}
	if err := c.BindJSON(&req); err != nil {
		log.Printf("c.BindJSON: %v\n", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	a, errFind := userSrv.GetAuth(*req.EMail, *req.Password)
	if errFind != nil {
		log.Printf("userSrv.errFind: %v\n", errFind)
		c.AbortWithError(http.StatusUnauthorized, errFind)
		return
	}
	// TODO: merge user data
	// previsou user is a.UserID
	// current user is c.GetString(AuthenticatedUserID)
	s, errCreate := userSrv.CreateSession(*a.UserId)
	if errCreate != nil {
		log.Printf("userSrv.errCreate: %v\n", errCreate)
		c.AbortWithError(http.StatusInternalServerError, errCreate)
		return
	}
	c.JSON(http.StatusOK, s)
}
