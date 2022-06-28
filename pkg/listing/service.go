package listing

import (
	"errors"

	"github.com/jrpespinas/kumot-coding-challenge/pkg/domain"
	"github.com/jrpespinas/kumot-coding-challenge/pkg/repository/api"
	"github.com/jrpespinas/kumot-coding-challenge/pkg/repository/cache"
	"github.com/rs/zerolog"
)

// Define a service interface for the business logic
type Service interface {
	// Check if the request of is appropriate
	Validate(usernames []string) error

	// Return the list of details by username
	ShowDetails(usernames []string) ([]domain.User, error)
}

// Define service struct that implements the service interface
type service struct {
	repository api.Repository
	caching    cache.Repository
	logger     zerolog.Logger
}

// Define a constructor to inject the Repository to this service
func NewService(repository api.Repository, caching cache.Repository, logger zerolog.Logger) Service {
	return &service{
		repository: repository,
		caching:    caching,
		logger:     logger,
	}
}

// Validate if given array length is in between 1 and 10
func (s *service) Validate(usernames []string) error {
	if len(usernames) > 10 || len(usernames) < 1 {
		return errors.New("please input between 1 to 10 usernames only")
	}
	return nil
}

// Return an error and a list of user details to the controller layer
func (s *service) ShowDetails(usernames []string) ([]domain.User, error) {
	var users []domain.User
	for _, user := range usernames {
		// Caching
		// 1. Get details from cache (if stored previously)
		s.logger.Info().Str("layer", "service").Msg("checking cache for existing user details")
		details, err := s.caching.Get(user)
		if err != nil {
			s.logger.Error().Str("layer", "service").Msg(err.Error())
			return users, errors.New("service unavailable")
		}

		// 2. return details from repository
		if details.Name == nil {
			s.logger.Info().Str("layer", "service").Msg("pulling from user details from api")
			details, _ = s.repository.GetDetails(user)
			if details.Name == nil {
				s.logger.Info().Str("layer", "service").Msg("user does not exist")
				continue
			}

			// 3. set details
			s.logger.Info().Str("layer", "service").Msg("store user details in cache")
			err = s.caching.Set(user, details)
			if err != nil {
				s.logger.Error().Str("layer", "service").Msg(err.Error())
				return users, errors.New("service unavailable")
			}
		}

		s.logger.Info().Str("layer", "service").Msg("appending record to slice")
		users = append(users, *details)
	}

	// Check if array is empty
	if len(users) == 0 {
		s.logger.Info().Str("layer", "service").Msg("no users returned")
		return users, errors.New("no users returned")
	}

	// Sort the user details alphabetically by `name`
	s.logger.Info().Str("layer", "service").Msg("sorting user details")
	SortUsers(users, "name")

	return users, nil
}
