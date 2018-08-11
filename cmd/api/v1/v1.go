package v1

import (
	"github.com/gin-gonic/gin"
)

const (
	// HeaderSessionToken is request header
	HeaderSessionToken = "X-C2C-Session-Token"
	// QueryToken is request query parameter
	QueryToken = "token"
	// ContextAuthorizedUserID is authorized userId
	ContextAuthorizedUserID = "ContextAuthorizedUserID"
)

// Router is routing v1
func Router(e *gin.Engine) {
	v1 := e.Group("/v1")
	{
		// public

		// after here, require to be authorized
		v1.Use(middlewareAuthorized)
	}
}

func middlewareAuthorized(c *gin.Context) {
}
