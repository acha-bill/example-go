package services

import (
	"errors"
	"time"

	"example/models"
	"example/pkg"
	"example/repo"
)

var (
	ErrNoActiveSubscription = errors.New("no active subscription")
)

const planDuration = 30 * 24 * time.Hour

var plans = map[models.PlanType]models.Plan{
	models.PlanTypeFree: {
		Type:     models.PlanTypeFree,
		Price:    0,
		Duration: 0,
	},
	models.PlanTypeBasic: {
		Type:     models.PlanTypeBasic,
		Price:    9.99,
		Duration: planDuration,
	},
	models.PlanTypePremium: {
		Type:     models.PlanTypePremium,
		Price:    19.99,
		Duration: planDuration,
	},
}

type subscription struct {
	r repo.Subscription
}

// Subscription is the interface that all subscription services must implement
type Subscription interface {
	// Create creates a new subscription
	Create(*models.Subscription) error
	// GetByID returns a subscription by its ID
	GetByID(key pkg.PrimaryKey) (*models.Subscription, error)
	// GetByUserID returns a subscription by its user ID
	GetByUserID(key pkg.PrimaryKey) ([]*models.Subscription, error)
	// GetByPlanType returns a subscription by its plan type
	GetByPlanType(planType models.PlanType) ([]*models.Subscription, error)
	// GetActiveForUser returns the active subscription for a user
	GetActiveForUser(key pkg.PrimaryKey) (*models.Subscription, error)
}

// NewSubscription returns a new Subscription service
func NewSubscription(r repo.Subscription) Subscription {
	return &subscription{r}
}

func (s *subscription) Create(m *models.Subscription) error {
	m.CreatedAt = time.Now()
	return s.r.Create(m)
}

func (s *subscription) GetByID(id pkg.PrimaryKey) (*models.Subscription, error) {
	return s.r.GetByID(id)
}

func (s *subscription) GetByUserID(id pkg.PrimaryKey) ([]*models.Subscription, error) {
	return s.r.GetBy(func(m *models.Subscription) bool {
		return m.UserID == id
	})
}

func (s *subscription) GetByPlanType(planType models.PlanType) ([]*models.Subscription, error) {
	return s.r.GetBy(func(m *models.Subscription) bool {
		return m.PlanType == planType
	})
}

func (s *subscription) GetActiveForUser(key pkg.PrimaryKey) (*models.Subscription, error) {
	subs, err := s.r.GetBy(func(m *models.Subscription) bool {
		return m.UserID == key && m.CreatedAt.Add(plans[m.PlanType].Duration).After(time.Now())
	})
	if err != nil {
		return nil, err
	}
	if len(subs) == 0 {
		return nil, ErrNoActiveSubscription
	}
	return subs[0], nil
}
