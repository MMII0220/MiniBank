package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY")) // преобразуем сразу в []byte

func (s *Service) Register(req domain.ReqRegister, role domain.Role) (domain.User, error) {
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

	err = s.repo.CreateUser(&user)
	if err != nil {
		return domain.User{}, err
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
			return user, err
		}

		// после создания account создаём карту
		card, err := s.CreateCardForAccount(account.ID, user.FullName)
		if err != nil {
			return domain.User{}, err
		}

		fmt.Printf("Создана карта: %s для пользователя %s\n", card.CardNumber, user.FullName)
	}

	// var reqLimit domain.Limit
	// Создаем стандартный лимит для нового пользователя
	err = s.repo.CreateDailyLimitForUser(user.ID, 1000.0) // 1000 TJS дневной лимит
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *Service) Login(req domain.ReqLogin) (string, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
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
func (s *Service) ParseToken(tokenStr string) (domain.User, error) {
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
