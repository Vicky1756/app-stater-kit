package utils

import (
	"app-starter-kit/backend/packages/db"
	"app-starter-kit/backend/packages/models"
	"database/sql"
	"regexp"
)

// Pre-compile regex for efficiency (Conscientious Developer tip)
var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func ValidateUser(user models.User) []string {
	var errs []string

	if !emailRegex.MatchString(user.Email) {
		errs = append(errs, "Invalid email address")
	}
	// Using 8 characters is a better security standard for a dev
	if len(user.Password) < 4 {
		errs = append(errs, "Password must be at least 4 characters")
	}
	if len(user.Name) < 1 {
		errs = append(errs, "Please enter your full name")
	}

	return errs
}

func ValidatePasswordReset(resetPassword models.ResetPasswordRequest) (bool, string) {
	if len(resetPassword.Password) < 4 {
		return false, "Password must be at least 4 characters"
	}
	if resetPassword.Password != resetPassword.ConfirmPassword {
		return false, "Passwords do not match"
	}
	return true, ""
}

// Move this to your models package or fix the 'model' vs 'models' typo
func UserExists(user *models.User, dbConn *sql.DB) bool {
	var count int
	err := dbConn.QueryRow(db.CheckUserExistsQuery, user.Email).Scan(&count)

	if err != nil {
		return false
	}
	return count > 0
}
