package utils

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

var allowedProfileImages = map[string]bool{
	"BRO-1-A": true, "BRO-1-B": true, "BRO-1-C": true,
	"BRO-2-A": true, "BRO-2-B": true, "BRO-2-C": true,
	"BRO-3-A": true, "BRO-3-B": true, "BRO-3-C": true,
	"BRO-4-A": true, "BRO-4-B": true, "BRO-4-C": true,
	"BRO-5-A": true, "BRO-5-B": true, "BRO-5-C": true,
	"BRO-6-A": true, "BRO-6-B": true, "BRO-6-C": true,
	"BRO-7-A": true, "BRO-7-B": true, "BRO-7-C": true,
	"BRO-8-A": true, "BRO-8-B": true, "BRO-8-C": true,
}

func ValidateProfileImage(profileImage string) error {
	if !allowedProfileImages[profileImage] {
		return errors.New("invalid profile image")
	}
	return nil
}

func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name is required")
	}
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	return nil
}
