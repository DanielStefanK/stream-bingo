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
	ErrResourceNotFound   = "ResourceNotFound"
)

func parseError(ctx *gin.Context) {
	log.Println("Error parsing request body")
	ctx.JSON(400, NewErrorResponse(ErrInvalidRequestBody, "could not parse request", nil))
}
