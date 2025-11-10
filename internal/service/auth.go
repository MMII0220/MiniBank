package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

func (s *Service) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", s.translateError(err)
	}
	return hex.EncodeToString(bytes), nil
}

func (s *Service) createAccessToken(userID int, role domain.Role) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"type":    "access",
		"exp":     time.Now().Add(time.Minute * 15).Unix(), // Короткий срок жизни - 15 минут
	})
	return token.SignedString(jwtSecret)
}

func (s *Service) createRefreshToken(userID int, role domain.Role) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 дней
	})
	return token.SignedString(jwtSecret)
}

func (s *Service) Register(req domain.ReqRegister, role domain.Role) (domain.User, error) {
	log := logger.GetLogger()
	log.Info().
		Str("email", req.Email).
		Str("role", string(role)).
		Msg("Starting user registration")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to hash password")
		return domain.User{}, s.translateError(err)
	}

	user := domain.User{
		FullName: req.FullName,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	err = s.repo.CreateUser(&user)
	if err != nil {
		return domain.User{}, s.translateError(err)
	}

	currencies := []string{"TJS", "USD", "EUR"}
	for _, currency := range currencies {
		account := domain.Account{
			UserID:   user.ID,
			Currency: currency,
			Balance:  "0",
			Blocked:  false,
		}

		err := s.repo.CreateAccount(&account)
		if err != nil {
			return user, s.translateError(err)
		}

		// после создания account создаём карту
		card, err := s.CreateCardForAccount(account.ID, user.FullName)
		if err != nil {
			return domain.User{}, s.translateError(err)
		}

		fmt.Printf("Создана карта: %s для пользователя %s\n", card.CardNumber, user.FullName)
	}

	// var reqLimit domain.Limit
	// Создаем стандартный лимит для нового пользователя
	err = s.repo.CreateDailyLimitForUser(user.ID, 1000.0) // 1000 TJS дневной лимит
	if err != nil {
		return user, s.translateError(err)
	}

	return user, nil
}

func (s *Service) Login(req domain.ReqLogin) (domain.TokenResponse, error) {
	var response domain.TokenResponse

	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return response, s.translateError(err) // Теперь вернет ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return response, s.translateError(err)
	}

	// Создаем access токен (15 минут)
	accessToken, err := s.createAccessToken(user.ID, user.Role)
	if err != nil {
		return response, s.translateError(err)
	}

	// Создаем refresh токен (7 дней)
	refreshToken, err := s.createRefreshToken(user.ID, user.Role)
	if err != nil {
		return response, s.translateError(err)
	}

	response = domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * 60, // 15 минут в секундах
		TokenType:    "Bearer",
	}

	return response, nil
}

// Обновление токенов через refresh токен
func (s *Service) RefreshToken(req domain.ReqRefreshToken) (domain.TokenResponse, error) {
	var response domain.TokenResponse

	// Парсим refresh токен
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return response, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return response, errors.New("invalid refresh token claims")
	}

	// Проверяем что это именно refresh токен
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return response, errors.New("invalid token type")
	}

	// Получаем данные пользователя
	userID := int(claims["user_id"].(float64))
	role := domain.Role(claims["role"].(string))

	// Проверяем срок действия
	exp := int64(claims["exp"].(float64))
	if time.Now().Unix() > exp {
		return response, errors.New("refresh token expired")
	}

	// Создаем новые токены
	newAccessToken, err := s.createAccessToken(userID, role)
	if err != nil {
		return response, s.translateError(err)
	}

	newRefreshToken, err := s.createRefreshToken(userID, role)
	if err != nil {
		return response, s.translateError(err)
	}

	response = domain.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    15 * 60, // 15 минут в секундах
		TokenType:    "Bearer",
	}

	return response, nil
}

// Проверка access токена
func (s *Service) ParseToken(tokenStr string) (domain.User, error) {
	var user domain.User

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return user, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return user, errors.New("invalid token claims")
	}

	// Проверяем что это access токен
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return user, errors.New("invalid token type, access token required")
	}

	user.ID = int(claims["user_id"].(float64))
	user.Role = domain.Role(claims["role"].(string))

	return user, nil
}
