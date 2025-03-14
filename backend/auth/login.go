package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/DanielStefanK/stream-bingo/config"
	"github.com/DanielStefanK/stream-bingo/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

	provider, exists := conf.OAuth.Providers[providerName]
	if !exists {
		log.Printf("Invalid OAuth provider: %s", providerName)
		return "", errors.New(ErrInvalidProvider)
	}

	url := provider.Config.AuthCodeURL("random-state")
	log.Printf("Redirecting user to %s OAuth URL: %s", providerName, url)
	return url, nil
}

// OAuthCallback handles OAuth callbacks for all providers
func OAuthCallback(providerName string, code string, db *gorm.DB) (string, *models.User, error) {
	log.Printf("OAuth callback received for provider: %s", providerName)

	provider, exists := conf.OAuth.Providers[providerName]
	if !exists {
		log.Printf("Invalid OAuth provider: %s", providerName)
		return "", nil, errors.New(ErrInvalidProvider)
	}

	if code == "" {
		log.Printf("Missing authorization code for provider: %s", providerName)
		return "", nil, errors.New(ErrMissingCode)
	}

	log.Printf("Exchanging code for token with provider: %s", providerName)
	token, err := provider.Config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Error exchanging code for token (provider: %s): %v", providerName, err)
		return "", nil, errors.New(ErrTokenExchangeFailed)
	}

	log.Printf("Fetching user info from provider: %s", providerName)
	userInfo, err := fetchOAuthUser(provider.UserURL, token.AccessToken)
	if err != nil {
		log.Printf("Error fetching user info (provider: %s): %v", providerName, err)
		return "", nil, errors.New(ErrFailedToGetUser)
	}

	log.Printf("User info retrieved: %+v", userInfo)

	user, err := CreateOrGetOAuthUser(db, providerName, userInfo["id"], userInfo["name"], userInfo["email"], userInfo["avatar_url"])
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
func fetchOAuthUser(userURL, accessToken string) (map[string]string, error) {
	log.Printf("Fetching user info from URL: %s", userURL)

	req, _ := http.NewRequest("GET", userURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to %s: %v", userURL, err)
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]string
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		log.Printf("Error decoding user info JSON from %s: %v", userURL, err)
		return nil, err
	}

	log.Printf("Successfully fetched user info from %s", userURL)
	return userInfo, nil
}

func CreateOrGetOAuthUser(db *gorm.DB, provider, providerID, name, email, avatarURL string) (*models.User, error) {
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
				log.Println("Error creating user")
				return nil, err
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
