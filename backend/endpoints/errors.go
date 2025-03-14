package endpoints

import (
	"log"

	"github.com/gin-gonic/gin"
)

const (
	ErrInternal           = "InternalError"
	ErrInvalidRequestBody = "InvalidRequestBody"
	ErrAuthenticatingUser = "AuthenticatingUser"
	ErrAuthorizingUser    = "AuthorizingUser"
	ErrMissingValue       = "MissingValue"
)

func parseError(ctx *gin.Context) {
	log.Println("Error parsing request body")
	ctx.JSON(400, newErrorResponse(ErrInvalidRequestBody, "could not parse request", nil))
}
