package auth

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ValidateJWT überprüft die Gültigkeit des Tokens und gibt die Claims zurück
func ValidateJWT(tokenString string) (*jwt.MapClaims, error) {
	log.Println("Validating JWT token")

	// Token parsen und validieren
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Sicherstellen, dass der Algorithmus übereinstimmt
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Unexpected signing method:", token.Header["alg"])
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Println("JWT parsing error:", err)
		return nil, err
	}

	// Claims extrahieren und Gültigkeit prüfen
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Ablaufdatum prüfen
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				log.Println("JWT token expired")
				return nil, errors.New(ErrJWTExpired)
			}
		}
		log.Println("JWT token is valid")
		return &claims, nil
	}

	log.Println("Invalid JWT token")
	return nil, errors.New(ErrJWTInvalid)
}

// GetUserFromToken extrahiert Benutzerinformationen aus einem validierten Token
func GetUserFromToken(tokenString string) (uint, error) {
	log.Println("Extracting user from JWT token")

	claims, err := ValidateJWT(tokenString)
	if err != nil {
		log.Println("Token validation failed:", err)
		return 0, err
	}

	// User ID extrahieren
	userIDFloat, ok := (*claims)["user_id"].(float64)
	if !ok {
		log.Println("user_id missing or invalid in token")
		return 0, errors.New(ErrJWTInvalid)
	}
	userID := uint(userIDFloat)

	log.Printf("Extracted user from token: ID=%d", userID)
	return userID, nil
}
