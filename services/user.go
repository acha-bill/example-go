package services

import (
	"example/models"
	"example/pkg"
	"example/repo"
)

// User is the interface that all user services must implement
type User interface {
	// Create creates a new user
	Create(*models.User) error
	// GetByID returns a user by its ID
	GetByID(key pkg.PrimaryKey) (*models.User, error)
	// GetByUsername returns a user by its username
	GetByUsername(username string) ([]*models.User, error)
}

type user struct {
	r repo.User
}

func (u *user) Create(m *models.User) error {
	return u.r.Create(m)
}

func (u *user) GetByID(key pkg.PrimaryKey) (*models.User, error) {
	return u.r.GetByID(key)
}

func (u *user) GetByUsername(username string) ([]*models.User, error) {
	return u.r.GetByUsername(username)
}

// NewUser returns a new user service
func NewUser(r repo.User) User {
	return &user{r}
}

var _ User = &user{}
