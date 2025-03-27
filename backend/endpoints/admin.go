package endpoints

import (
	"strconv"

	"github.com/DanielStefanK/stream-bingo/db"
	"github.com/DanielStefanK/stream-bingo/models"
	"github.com/gin-gonic/gin"
)

func DeactiveUser(ctx *gin.Context) {
	//parse request
	var req UserStateChangeRequest
	if err := ctx.BindJSON(&req); err != nil {
		parseError(ctx)
		return
	}

	// user id from param from string to number
	userIdStr := ctx.Param("userId")
	userId, err := strconv.Atoi(userIdStr)

	if err != nil || userId == 0 {
		ctx.JSON(400, NewErrorResponse(ErrMissingValue, "missing user id", map[string]interface{}{"fieldName": "userId"}))
		return
	}

	//deactive user
	user := &models.User{}
	db.GetDB().First(user, userId)
	// check if user exists
	if user.ID == 0 {
		ctx.JSON(404, NewErrorResponse(ErrResourceNotFound, "user not found", map[string]interface{}{"identifier": userId, "type": "user", "identifierName": "id"}))
		return
	}
	user.Active = req.Active
	db.GetDB().Save(user)

	//return success
	ctx.JSON(200, newSuccessResponse(createUserDtoWithAllFields(user)))
}

type UserStateChangeRequest struct {
	Active bool `json:"active"`
}

func GetUsers(ctx *gin.Context) {
	page, limit, sortBy := parsePaginationQueryParams(ctx)
	usersDb := []models.User{}
	var total int64 = 0

	db.GetDB().Model(&models.User{}).Count(&total)
	db.GetDB().Unscoped().Order(sortBy).Offset((page - 1) * limit).Limit(limit).Find(&usersDb)

	//map users to dto
	users := make([]UserWithAllFields, len(usersDb))
	for i, user := range usersDb {
		users[i] = createUserDtoWithAllFields(&user)
	}

	//return success
	ctx.JSON(200, newSuccessResponse(PaginationResponse{
		Total: total,
		Data:  users,
	}))
}

func DeleteUser(ctx *gin.Context) {

	// user Id frpm param
	userId := ctx.Param("userId")

	//delete user
	user := &models.User{}
	db.GetDB().First(user, userId)
	// check if user exists
	if user.ID == 0 {
		ctx.JSON(404, NewErrorResponse(ErrResourceNotFound, "user not found", map[string]interface{}{"identifier": userId, "type": "user", "identifierName": "id"}))
		return
	}
	db.GetDB().Delete(user)

	//return success
	ctx.JSON(200, newSuccessResponse(createUserDtoWithAllFields(user)))
}

func createUserDtoWithAllFields(user *models.User) UserWithAllFields {
	return UserWithAllFields{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		Provider:  user.Provider,
		Admin:     user.Admin,
		Active:    user.Active,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		DeleteAt: func() *string {
			if user.DeletedAt.Valid {
				str := user.DeletedAt.Time.String()
				return &str
			}
			return nil
		}(),
		ProviderID: user.ProviderID,
	}
}

type UserWithAllFields struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	AvatarURL  string  `json:"avatar"`
	Provider   string  `json:"provider"`
	Admin      bool    `json:"admin"`
	Active     bool    `json:"active"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
	DeleteAt   *string `json:"deleteAt"`
	ProviderID string  `json:"providerId"`
}
