package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jeecg/jimureport-go/internal/config"
	"github.com/jeecg/jimureport-go/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Claims struct {
	UserID   string   `json:"userId"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	TenantID string   `json:"tenantId"`
	jwt.RegisteredClaims
}

var jwtCfg *config.JWTConfig

func InitJWT(cfg *config.JWTConfig) {
	jwtCfg = cfg
}

func GenerateToken(user *models.User) (string, error) {
	if jwtCfg == nil {
		return "", errors.New("JWT config not initialized")
	}

	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Roles:    []string{user.Role},
		TenantID: user.TenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtCfg.Issuer,
			Subject:   user.ID,
			Audience:  []string{jwtCfg.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtCfg.Secret))

	return tokenString, err
}

func ValidateToken(tokenString string) (*Claims, error) {
	if jwtCfg == nil {
		return nil, errors.New("JWT config not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtCfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserByCredentials(db *gorm.DB, username, password string) (*models.User, error) {
	var user models.User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	if !CheckPassword(password, user.Password) {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
