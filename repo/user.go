package repo

import (
	"fmt"

	"example/models"
	"example/pkg"
)

const usersTable = "users"

// User is the interface that all user repositories must implement
type User interface {
	// Create creates a new user
	Create(*models.User) error
	// GetByID returns a user by its ID
	GetByID(key pkg.PrimaryKey) (*models.User, error)
	// GetByUsername returns a user by its username
	GetByUsername(username string) ([]*models.User, error)
	// FindAll returns all users
	FindAll() ([]*models.User, error)
}

type user struct {
	db pkg.DB
}

func (u *user) Create(m *models.User) error {
	table, err := u.db.Table(usersTable)
	if err != nil {
		return fmt.Errorf("error getting table: %w", err)
	}
	if err = table.Insert(m); err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	return nil
}

func (u *user) GetByID(key pkg.PrimaryKey) (*models.User, error) {
	table, err := u.db.Table(usersTable)
	if err != nil {
		return nil, fmt.Errorf("error getting table: %w", err)
	}
	model, err := table.Get(key)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return model.(*models.User), nil
}

func (u *user) GetByUsername(username string) ([]*models.User, error) {
	table, err := u.db.Table(usersTable)
	if err != nil {
		return nil, fmt.Errorf("error getting table: %w", err)
	}
	ms, err := table.Find(func(model pkg.Model) bool {
		user := model.(*models.User)
		return user.Username == username
	})
	if err != nil {
		return nil, fmt.Errorf("error finding users: %w", err)
	}
	users := make([]*models.User, len(ms))
	for i, model := range ms {
		users[i] = model.(*models.User)
	}
	return users, nil
}

func (u *user) FindAll() ([]*models.User, error) {
	table, err := u.db.Table(usersTable)
	if err != nil {
		return nil, fmt.Errorf("error getting table: %w", err)
	}
	ms, err := table.Find(func(model pkg.Model) bool {
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("error finding users: %w", err)
	}
	users := make([]*models.User, len(ms))
	for i, model := range ms {
		users[i] = model.(*models.User)
	}
	return users, nil
}

// NewUser returns a new user repository
func NewUser(db pkg.DB) (User, error) {
	err := db.AddTable(usersTable)
	if err != nil {
		return nil, fmt.Errorf("error adding table: %w", err)
	}
	return &user{db}, nil
}

var _ User = &user{}
