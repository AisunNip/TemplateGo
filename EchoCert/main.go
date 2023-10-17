package main

import (
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type (
	User struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func main() {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/show", show) //http://localhost:8080/show?team=ITID&member=Aisun

	e.Validator = &CustomValidator{validator: validator.New()}
	e.POST("/users", func(c echo.Context) (err error) {
		u := new(User)
		if err = c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(u); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, u)
	})
	// 	curl -X POST http://localhost:1323/users \
	// 	-H 'Content-Type: application/json' \
	// 	-d '{"name":"Joe","email":"joe@invalid-domain"}'
	//   {"message":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"}

	// if err := e.Start(":8080"); err != http.ErrServerClosed {
	// 	log.Fatal(err)
	// }

	if err := e.StartTLS(":443", "Cert\\server.pem", "Cert\\private.pem"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func show(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
