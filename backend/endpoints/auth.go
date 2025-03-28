package endpoints

import (
	"log"
	"net/http"

	"github.com/DanielStefanK/stream-bingo/auth"
	"github.com/DanielStefanK/stream-bingo/config"
	"github.com/DanielStefanK/stream-bingo/db"
	"github.com/DanielStefanK/stream-bingo/models"
	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	//parts request body
	var req LoginRequest = LoginRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		parseError(ctx)
		return
	}

	user, err := auth.AuthenticateLocalUser(db.GetDB(), req.Email, req.Password)
	if err != nil {
		log.Println("Error authenticating user")
		ctx.JSON(401, NewErrorResponse(ErrAuthenticatingUser, "could not authenticate user", nil))
		return
	}

	token, err := auth.GenerateJWT(user)
	if err != nil {
		log.Println("Error generating token")
		ctx.JSON(500, NewErrorResponse(ErrInternal, "could not authenticate user", nil))
		return
	}
	ctx.JSON(200, newSuccessResponse(LoginResponse{
		Token: token,
		User: User{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			Provider:  user.Provider,
			Admin:     user.Admin,
		},
	}))
}

func Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.BindJSON(&req); err != nil {
		parseError(ctx)
		return
	}

	user, err := auth.CreateLocalUser(db.GetDB(), req.Name, req.Email, req.Password)
	if err != nil {
		log.Println("Error creating user", err)
		ctx.JSON(500, NewErrorResponse(ErrInternal, "could not create user", nil))
		return
	}

	token, err := auth.GenerateJWT(user)
	if err != nil {
		ctx.JSON(500, NewErrorResponse(ErrInternal, "could not authenticate user", nil))
		return
	}
	ctx.JSON(200, newSuccessResponse(LoginResponse{
		Token: token,
		User: User{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			Provider:  user.Provider,
			Admin:     user.Admin,
		},
	}))
}
func OAuthRedirect(ctx *gin.Context) {
	config.ReloadConfig()
	providerName := ctx.Param("provider")
	log.Printf("OAuth redirect request for provider: %s", providerName)

	// Prüfen, ob der Anbieter existiert
	provider, exists := auth.ProviderFromConfig(providerName)
	if exists != nil {
		log.Printf("Invalid OAuth provider: %s", providerName)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	// Generiere die OAuth-URL
	url := provider.AuthCodeURL("random-state")
	log.Printf("Redirecting user to %s OAuth URL: %s", providerName, url)

	// Weiterleitung zur OAuth-Login-Seite
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}
func OAuthCallback(ctx *gin.Context) {
	providerName := ctx.Param("provider")
	code := ctx.Query("code")
	log.Printf("OAuth callback received for provider: %s", providerName)

	// validate

	code, user, err := auth.OAuthCallback(providerName, code, db.GetDB())
	if err != nil {
		var rs *Response
		if err.Error() == auth.ErrInvalidProvider {
			rs = NewErrorResponse(ErrAuthenticatingUser, "invalid provider", nil)
		} else if err.Error() == auth.ErrMissingCode {
			rs = NewErrorResponse(ErrMissingValue, "missing code", map[string]interface{}{"fieldname": "code"})
		} else if err.Error() == auth.ErrTokenExchangeFailed {
			rs = NewErrorResponse(ErrAuthenticatingUser, "could not perform code exchange", nil)
		} else if err.Error() == auth.ErrFailedToGetUser {
			rs = NewErrorResponse(ErrAuthenticatingUser, "could not fetch data from user provider", nil)
		} else if err.Error() == auth.ErrCreationFailed {
			rs = NewErrorResponse(ErrInternal, "failed to create user", nil)
		} else if err.Error() == auth.ErrJWTSigningFailed {
			rs = NewErrorResponse(ErrInternal, "failed to create Token", nil)
		} else {
			rs = NewErrorResponse(ErrInternal, "internal error", nil)
		}

		ctx.JSON(http.StatusInternalServerError, rs)
		return
	}

	resp := newSuccessResponse(LoginResponse{
		Token: code,
		User: User{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			Provider:  user.Provider,
			Admin:     user.Admin,
		},
	})
	ctx.JSON(http.StatusOK, resp)
}

func Me(ctx *gin.Context) {
	// get user from token
	userI, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(401, NewErrorResponse(ErrAuthorizingUser, "unauthorized service", nil))
		return
	}

	// get user
	user := userI.(*models.User)

	ctx.JSON(http.StatusOK, newSuccessResponse(User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		Provider:  user.Provider,
		Admin:     user.Admin,
	}))

}

type User struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar"`
	Provider  string `json:"provider"`
	Admin     bool   `json:"admin"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterRequest struct {
	LoginRequest
	Name string `json:"name"`
}
