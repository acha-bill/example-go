package models

import (
	"time"

	"example/pkg"
)

type PlanType string

const (
	PlanTypeFree    PlanType = "free"
	PlanTypeBasic   PlanType = "basic"
	PlanTypePremium PlanType = "premium"
)

type Plan struct {
	Type     PlanType
	Price    float32
	Duration time.Duration
}

type Subscription struct {
	ID        pkg.PrimaryKey
	UserID    pkg.PrimaryKey
	PlanType  PlanType
	CreatedAt time.Time
}

func (s *Subscription) GetID() pkg.PrimaryKey {
	return s.ID
}
func (s *Subscription) SetID(id pkg.PrimaryKey) {
	s.ID = id
}

var _ pkg.Model = (*Subscription)(nil)
