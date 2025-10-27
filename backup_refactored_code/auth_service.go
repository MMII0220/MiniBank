package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY")) // преобразуем сразу в []byte

func Register(req domain.ReqRegister, role domain.Role) (domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		FullName: req.FullName,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	err = repository.CreateUser(&user)
	if err != nil {
		return domain.User{}, err
	}

	// Создаём дефолтные счета для пользователя
	currencies := []string{"TJS", "USD", "EUR"}
	err = CreateAccountsForUser(user.ID, currencies)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create accounts: %w", err)
	}

	return user, nil
}

func Login(req domain.ReqLogin) (string, error) {
	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Создаем JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret) // jwtSecret теперь []byte
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Проверка JWT и роль
func ParseToken(tokenStr string) (domain.User, error) {
	var user domain.User

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil // тоже []byte
	})
	if err != nil || !token.Valid {
		return user, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		user.ID = int(claims["user_id"].(float64))
		user.Role = domain.Role(claims["role"].(string))
		return user, nil
	}

	return user, errors.New("invalid token claims")
}
