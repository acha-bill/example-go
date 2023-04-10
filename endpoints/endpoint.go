package endpoints

import "github.com/labstack/echo/v4"

type Endpoint interface {
	Register(g *echo.Group)
}
