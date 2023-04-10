package endpoints

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"example/models"
	"example/pkg"
	"example/services"
)

type Subscription struct {
	subscriptionService services.Subscription
	userService         services.User
}
type createSubscriptionRequest struct {
	UserID   int             `json:"user_id"`
	PlanType models.PlanType `json:"plan_type"`
}

func NewSubscription(s services.Subscription, userService services.User) *Subscription {
	return &Subscription{
		subscriptionService: s,
		userService:         userService,
	}
}

func (s *Subscription) Register(g *echo.Group) {
	g.POST("", s.Create)
	g.GET("", s.Find)
	g.GET("/:id", s.GetByID)
	g.GET("/users/:user_id", s.FindByUser)
	g.GET("/users/:user_id/active", s.FindActive)
}

func (s *Subscription) Create(c echo.Context) error {
	var req createSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request: %s", err.Error()))
	}
	if req.PlanType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "plan_type is required")
	}
	if _, err := s.userService.GetByID(pkg.PrimaryKey(req.UserID)); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("user not found: %s", err.Error()))
	}
	subscription := &models.Subscription{UserID: pkg.PrimaryKey(req.UserID), PlanType: req.PlanType}
	if err := s.subscriptionService.Create(subscription); err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, services.ErrInvalidPlanType) {
			statusCode = http.StatusBadRequest
		}
		return echo.NewHTTPError(statusCode, err)
	}
	return c.JSON(http.StatusOK, subscription)
}

func (s *Subscription) GetByID(c echo.Context) error {
	subscriptionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid subscription id: %s", err.Error()))
	}
	subscription, err := s.subscriptionService.GetByID(pkg.PrimaryKey(subscriptionID))
	if err == nil {
		return c.JSON(http.StatusOK, subscription)
	}
	if errors.Is(err, pkg.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err)
}

func (s *Subscription) FindByUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid user id: %s", err.Error()))
	}
	subscriptions, err := s.subscriptionService.GetByUserID(pkg.PrimaryKey(userID))
	if err == nil {
		return c.JSON(http.StatusOK, subscriptions)
	}
	if errors.Is(err, pkg.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusInternalServerError, err)
}

func (s *Subscription) FindActive(c echo.Context) error {
	userID, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid user id: %s", err.Error()))
	}
	subscriptions, err := s.subscriptionService.GetActiveForUser(pkg.PrimaryKey(userID))
	if err == nil {
		return c.JSON(http.StatusOK, subscriptions)
	}
	if errors.Is(err, pkg.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusInternalServerError, err)
}

func (s *Subscription) Find(c echo.Context) error {
	subscriptions, err := s.subscriptionService.Find()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, subscriptions)
}
