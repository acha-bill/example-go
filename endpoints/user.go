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

// User is the endpoint for users
type User struct {
	userService services.User
}

type createUserRequest struct {
	Username string `json:"username"`
}

// NewUser returns a new user endpoint
func NewUser(s services.User) *User {
	return &User{userService: s}
}

// Register registers the user endpoint
func (u *User) Register(g *echo.Group) {
	g.GET("/:id", u.GetByID)
	g.GET("", u.Find)
	g.POST("", u.Create)
}

func (u *User) Create(c echo.Context) error {
	var req createUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request: %s", err.Error()))
	}
	if req.Username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username is required")
	}
	user := &models.User{Username: req.Username}
	if err := u.userService.Create(user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

func (u *User) GetByID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid user id: %s", err.Error()))
	}
	user, err := u.userService.GetByID(pkg.PrimaryKey(userID))
	if err == nil {
		return c.JSON(http.StatusOK, user)
	}
	if errors.Is(err, pkg.ErrNotFound) {
		return c.NoContent(http.StatusNotFound)
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

func (u *User) Find(c echo.Context) error {
	username := c.Param("username")
	var users []*models.User
	var err error
	if username != "" {
		users, err = u.userService.GetByUsername(username)
	} else {
		users, err = u.userService.FindAll()
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}
