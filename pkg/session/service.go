package session

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strconv"

	"github.com/jrpespinas/kumot-coding-challenge/pkg/repository/cache"
)

// Define a Service interface for our session
type Service interface {
	GenerateToken() (string, error)
	VerifyToken(tokenString string) error
}

// Define session structure
type service struct {
	sessionStore cache.Repository
}

// Define constructor and inject cache repository
func NewService(sessionStore cache.Repository) Service {
	return &service{
		sessionStore: sessionStore,
	}
}

// Generate a simple random alphanumeric string
func (s *service) GenerateToken() (string, error) {
	length, _ := strconv.Atoi(os.Getenv("TOKEN_LENGTH"))

	// Generate random token
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("service unavailable")
	}
	token := hex.EncodeToString(b)

	// Cache session
	s.sessionStore.SetSession(token, 1)
	return token, nil
}

// Verify generated token string
func (s *service) VerifyToken(tokenString string) error {
	val, err := s.sessionStore.GetSession(tokenString)
	if val != 1 {
		return errors.New("session does not exist")
	}

	if err != nil {
		return errors.New("service unavailable")
	}
	return nil
}
