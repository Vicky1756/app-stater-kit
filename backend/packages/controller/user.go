package controller

import (
	"app-starter-kit/backend/packages/db"
	"app-starter-kit/backend/packages/models"
	"app-starter-kit/backend/packages/utils"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context, dbConn *sql.DB) error {
	req := new(models.LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	var dbUser models.User
	// Using the query logic you established
	err := dbConn.QueryRow(db.LoginQuery, req.Email).Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.Password)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Database error"})
	}

	// Create the JWT Claims
	claims := &jwt.StandardClaims{
		Subject:   fmt.Sprintf("%d", dbUser.ID),
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	}

	// Sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not generate token"})
	}

	// Set JWT as a Cookie
	cookie := new(http.Cookie)
	cookie.Name = "access_token"
	cookie.Value = t
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.HttpOnly = true // Prevents JavaScript access (XSS protection)
	cookie.Secure = false  // Set to true only if using HTTPS
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteLaxMode // CSRF protection

	c.SetCookie(cookie)

	// Return user info.
	return c.JSON(http.StatusOK, models.UserResponseWithTokens{
		ID:          dbUser.ID,
		Name:        dbUser.Name,
		Email:       dbUser.Email,
		AccessToken: t,
	})
}
func CreateUser(c echo.Context, dbConn *sql.DB) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	//  User Validation
	errs := utils.ValidateUser(*user)
	if len(errs) > 0 {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"success": false,
			"errors":  errs,
		})
	}

	// Check existence
	if utils.UserExists(user, dbConn) {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Email already exists"})
	}

	// Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error hashing password"})
	}
	user.Password = string(hashed)

	// Timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Use Scan to put the Postgres-generated ID into the user struct
	err = dbConn.QueryRow(db.Registerquery,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Database error: " + err.Error()})
	}

	// Success Response
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}
