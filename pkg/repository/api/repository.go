package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jrpespinas/kumot-coding-challenge/pkg/domain"
)

// Define data access interface that will be used by the service
type Repository interface {
	GetDetails(username string) (*domain.User, error)
}

// Define a repository structure that implements the interface
type repository struct {
	url string
}

// Define a constructor to inject the github api
func NewRepository(url string) Repository {
	return &repository{
		url: url,
	}
}

func (r *repository) GetDetails(username string) (*domain.User, error) {
	var userData *domain.User

	// Get the user details from the url
	url := fmt.Sprintf("%v/%v", r.url, username)
	res, err := http.Get(url)
	if err != nil {
		return userData, errors.New("unable to process your request")
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&userData); err != nil {
		return userData, errors.New("user does not exist")
	}
	return userData, nil
}
