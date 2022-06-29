package listing

import (
	"os"
	"testing"

	"github.com/jrpespinas/kumot-coding-challenge/pkg/domain"
	"github.com/jrpespinas/kumot-coding-challenge/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define Mock Repository for our Github API
type MockGithubRepository struct {
	mock.Mock
}

// Define mock ShowDetails
func (m *MockGithubRepository) GetDetails(usernames string) (*domain.User, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(*domain.User), args.Error(1)
}

// Define Mock Repository for our Redis Cache
type MockRedisRepository struct {
	mock.Mock
}

// Define mock Get method
func (m *MockRedisRepository) Get(key string) (*domain.User, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(*domain.User), args.Error(1)
}

// Define mock Set method
func (m *MockRedisRepository) Set(key string, user *domain.User) error {
	args := m.Called()
	return args.Error(0)
}

func TestValidateEmptyList(t *testing.T) {
	testService := NewService(nil, nil, logging.NewLogger(os.Getenv("LOG_LEVEL")))
	err := testService.Validate(nil)

	// Assert Nil
	assert.NotNil(t, err)

	// Assert Error message
	expected := "please input between 1 to 10 usernames only"
	actual := err.Error()
	assert.Equal(t, expected, actual)

}

func TestValidateMoreThanTenUsers(t *testing.T) {
	testService := NewService(nil, nil, logging.NewLogger(os.Getenv("LOG_LEVEL")))

	// mock usernames
	usernames := []string{"user1", "user2", "user3", "user4", "user5", "user6", "user7", "user8", "user9", "user10", "user11"}
	err := testService.Validate(usernames)

	// Assert Nil
	assert.NotNil(t, err)

	// Assert Error message
	expected := "please input between 1 to 10 usernames only"
	actual := err.Error()
	assert.Equal(t, expected, actual)
}

func TestShowDetailsFromGithub(t *testing.T) {
	mockGithub := new(MockGithubRepository)
	mockRedis := new(MockRedisRepository)

	// Define mock username
	mockUsername := []string{"jrpespinas"}

	name := "Jan Rodolf Espinas"
	login := "jrpespinas"
	company := "@demandscience"
	followers := 8
	publicRepos := 19

	// Define mock Output
	mockOutput := domain.User{
		Name:        &name,
		Login:       &login,
		Company:     &company,
		Followers:   &followers,
		PublicRepos: &publicRepos,
	}

	// Setup expectation:
	// Mock missed Cache
	mockGithub.On("GetDetails").Return(&mockOutput, nil)
	mockRedis.On("Get").Return(&domain.User{}, nil)
	mockRedis.On("Set").Return(nil)

	// Initialize Listing Service
	testService := NewService(mockGithub, mockRedis, logging.NewLogger(os.Getenv("LOG_LEVEL")))

	// Mock Get details
	result, _ := testService.ShowDetails(mockUsername)

	// Mock Behavior Assertion
	mockGithub.AssertExpectations(t)
	mockRedis.AssertExpectations(t)

	// Data Assertion
	assert.NotNil(t, result[0].Name)
	assert.Equal(t, result[0].Name, &name)
	assert.Equal(t, result[0].Login, &login)
	assert.Equal(t, result[0].Company, &company)
	assert.Equal(t, result[0].Followers, &followers)
	assert.Equal(t, result[0].PublicRepos, &publicRepos)
}

func TestShowDetailsFromRedis(t *testing.T) {
	mockGithub := new(MockGithubRepository)
	mockRedis := new(MockRedisRepository)

	// Define mock username
	mockUsername := []string{"jrpespinas"}

	name := "Jan Rodolf Espinas"
	login := "jrpespinas"
	company := "@demandscience"
	followers := 8
	publicRepos := 19

	// Define mock Output
	mockOutput := domain.User{
		Name:        &name,
		Login:       &login,
		Company:     &company,
		Followers:   &followers,
		PublicRepos: &publicRepos,
	}

	// Setup expectation
	// Mock Cache Hit
	mockRedis.On("Get").Return(&mockOutput, nil)

	// Initialize Listing Service
	testService := NewService(mockGithub, mockRedis, logging.NewLogger(os.Getenv("LOG_LEVEL")))

	// Mock Get details
	result, _ := testService.ShowDetails(mockUsername)

	// Mock Behavior Assertion
	mockGithub.AssertExpectations(t)
	mockRedis.AssertExpectations(t)

	// Data Assertion
	assert.NotNil(t, result[0].Name)
	assert.Equal(t, result[0].Name, &name)
	assert.Equal(t, result[0].Login, &login)
	assert.Equal(t, result[0].Company, &company)
	assert.Equal(t, result[0].Followers, &followers)
	assert.Equal(t, result[0].PublicRepos, &publicRepos)
}
