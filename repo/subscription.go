package repo

import (
	"fmt"

	"example/models"
	"example/pkg"
)

// Subscription is a repository for subscriptions.
type Subscription interface {
	// Create creates a new subscription
	Create(*models.Subscription) error
	// GetByID returns a subscription by its ID
	GetByID(key pkg.PrimaryKey) (*models.Subscription, error)
	// GetBy returns a subscription by a filter function
	GetBy(filter func(*models.Subscription) bool) ([]*models.Subscription, error)
}

const subscriptionsTable = "subscription"

type subscription struct {
	db pkg.DB
}

func (s *subscription) Create(m *models.Subscription) error {
	table, err := s.db.Table(subscriptionsTable)
	if err != nil {
		return fmt.Errorf("error getting table: %w", err)
	}
	if err = table.Insert(m); err != nil {
		return fmt.Errorf("error inserting subscription: %w", err)
	}
	return nil
}

func (s *subscription) GetByID(key pkg.PrimaryKey) (*models.Subscription, error) {
	table, err := s.db.Table(subscriptionsTable)
	if err != nil {
		return nil, fmt.Errorf("error getting table: %w", err)
	}
	model, err := table.Get(key)
	if err != nil {
		return nil, fmt.Errorf("error getting subscription: %w", err)
	}
	return model.(*models.Subscription), nil
}

func (s *subscription) GetBy(filter func(*models.Subscription) bool) ([]*models.Subscription, error) {
	table, err := s.db.Table(subscriptionsTable)
	if err != nil {
		return nil, fmt.Errorf("error getting table: %w", err)
	}
	ms, err := table.Find(func(model pkg.Model) bool {
		subscription := model.(*models.Subscription)
		return filter(subscription)
	})
	if err != nil {
		return nil, fmt.Errorf("error finding subscriptions: %w", err)
	}
	var subscriptions []*models.Subscription
	for _, m := range ms {
		subscriptions = append(subscriptions, m.(*models.Subscription))
	}
	return subscriptions, nil
}

func NewSubscription(db pkg.DB) (Subscription, error) {
	if err := db.AddTable(subscriptionsTable); err != nil {
		return nil, fmt.Errorf("error adding table: %w", err)
	}
	return &subscription{db: db}, nil
}

var _ Subscription = (*subscription)(nil)
