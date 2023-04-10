package main

import (
	"github.com/labstack/echo/v4"

	"example/endpoints"
	"example/pkg"
	"example/repo"
	"example/services"
)

func main() {
	e := echo.New()

	db := pkg.NewDB()
	userRepo, err := repo.NewUser(db)
	if err != nil {
		e.Logger.Fatalf("failed to create user repo: %s", err.Error())
	}
	userService := services.NewUser(userRepo)
	userEndpoint := endpoints.NewUser(userService)
	userEndpoint.Register(e.Group("/users"))

	e.Logger.Fatal(e.Start(":8080"))
}
