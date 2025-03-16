package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/DanielStefanK/stream-bingo/config"
	"github.com/DanielStefanK/stream-bingo/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var jwtSecret = []byte("secret")

var conf = config.GetConfig()

// CreateLocalUser creates a new user with email & password authentication
func CreateLocalUser(db *gorm.DB, name, email, password string) (*models.User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Println("Error hashing password:")
		return nil, err
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Provider: "local",
	}

	if err := db.Create(&user).Error; err != nil {
		log.Println("Error creating user")
		return nil, err
	}

	return &user, nil
}

// AuthenticateLocalUser verifies login credentials
func AuthenticateLocalUser(db *gorm.DB, email, password string) (*models.User, error) {
	var user models.User

	if err := db.Where("email = ? AND provider = 'local'", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("User not found")
		} else {
			log.Println("Error finding user")
		}

		return nil, err
	}

	if err := CheckPassword(user.Password, password); err != nil {
		log.Println("Invalid password")
		return nil, err
	}

	return &user, nil
}

func CreateOrGetUser(db *gorm.DB, provider, providerID, name, email, avatarURL string) (*models.User, error) {
	var user models.User
	if err := db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			user = models.User{
				Name:       name,
				Email:      email,
				Provider:   provider,
				ProviderID: providerID,
				AvatarURL:  avatarURL,
			}
			if err := db.Create(&user).Error; err != nil {
				log.Println("Error creating user", err)
				return nil, err
			}
		} else {
			log.Println("Error finding user", err)
			return nil, err
		}
	}
	return &user, nil
}

// GenerateJWT creates a JWT token for authenticated users
func GenerateJWT(user *models.User) (string, error) {
	log.Printf("Generating JWT for user ID: %d, email: %s", user.ID, user.Email)

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("Error generating JWT for user %d: %v", user.ID, err)
		return "", err
	}

	log.Printf("JWT generated successfully for user ID: %d", user.ID)
	return signedToken, nil
}

// OAuthLoginRedirect redirects user to the OAuth provider
func OAuthLoginRedirect(providerName string) (string, error) {
	log.Printf("OAuth login request for provider: %s", providerName)

	provider, exists := ProviderFromConfig(providerName)
	if exists != nil {
		log.Printf("Invalid OAuth provider: %s", providerName)
		return "", errors.New(ErrInvalidProvider)
	}

	url := provider.AuthCodeURL("random-state")
	log.Printf("Redirecting user to %s OAuth URL: %s", providerName, url)
	return url, nil
}

// OAuthCallback handles OAuth callbacks for all providers
func OAuthCallback(providerName string, code string, db *gorm.DB) (string, *models.User, error) {
	log.Printf("OAuth callback received for provider: %s", providerName)

	provider, exists := ProviderFromConfig(providerName)
	if exists != nil {
		log.Printf("Invalid OAuth provider: %s", providerName)
		return "", nil, errors.New(ErrInvalidProvider)
	}

	if code == "" {
		log.Printf("Missing authorization code for provider: %s", providerName)
		return "", nil, errors.New(ErrMissingCode)
	}

	log.Printf("Exchanging code for token with provider: %s", providerName)
	token, err := provider.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Error exchanging code for token (provider: %s): %v", providerName, err)
		return "", nil, errors.New(ErrTokenExchangeFailed)
	}

	log.Printf("Fetching user info from provider: %s", providerName)
	userInfo, err := fetchOAuthUser(providerName, token.AccessToken)
	if err != nil {
		log.Printf("Error fetching user info (provider: %s): %v", providerName, err)
		return "", nil, errors.New(ErrFailedToGetUser)
	}

	log.Printf("User info retrieved")

	user, err := CreateOrGetOAuthUser(db, providerName, userInfo.ID, userInfo.Name, userInfo.Email, userInfo.AvatarURL)
	if err != nil {
		log.Printf("Error creating or retrieving user (provider: %s): %v", providerName, err)
		return "", nil, errors.New(ErrCreationFailed)
	}

	log.Printf("User authenticated successfully: ID %d, email: %s", user.ID, user.Email)

	jwtToken, err := GenerateJWT(user)
	if err != nil {
		log.Printf("Error generating JWT for user ID %d: %v", user.ID, err)
		return "", nil, errors.New(ErrJWTSigningFailed)
	}

	log.Printf("OAuth login successful, returning JWT for user ID %d", user.ID)
	return jwtToken, user, nil
}

// fetchOAuthUser fetches user info from an OAuth provider
func fetchOAuthUser(providerName, accessToken string) (UserInfo, error) {
	info, _ := GetUserInfoFromProvider(providerName, accessToken)
	return info, nil
}

func CreateOrGetOAuthUser(db *gorm.DB, provider, providerID, name, email, avatarURL string) (*models.User, error) {
	var user models.User

	log.Printf("Creating or getting user: %s, %s", provider, providerID)

	if err := db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			user = models.User{
				Name:       name,
				Email:      email,
				Provider:   provider,
				ProviderID: providerID,
				AvatarURL:  avatarURL,
			}
			if err := db.Create(&user).Error; err != nil {
				log.Println("Error creating user")
				db.Rollback()
			}
		} else {
			log.Println("Error accessing user db")
			return nil, err
		}
	}

	return &user, nil
}

// CheckPassword verifies a hashed password
func CheckPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// HashPassword securely hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func ProviderFromConfig(providerName string) (oauth2.Config, error) {
	config.ReloadConfig()
	conf = config.GetConfig()
	provider, exists := conf.OAuth.Providers[providerName]
	if !exists {
		log.Printf("Invalid OAuth provider: %s", providerName)
		return oauth2.Config{}, errors.New(ErrInvalidProvider)
	}

	providerConfig := oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.Endpoint.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  provider.Endpoint.AuthURL,
			TokenURL: provider.Endpoint.TokenURL,
		},
		Scopes: provider.Scopes,
	}

	return providerConfig, nil
}
