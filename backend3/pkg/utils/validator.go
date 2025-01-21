// internal/utils/validator.go
package utils

import (
    "regexp"
    "unicode"
)

var (
    emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

type ValidationError struct {
    Field   string
    Message string
}

func ValidateEmail(email string) *ValidationError {
    if !emailRegex.MatchString(email) {
        return &ValidationError{
            Field:   "email",
            Message: "invalid email format",
        }
    }
    return nil
}

func ValidatePassword(password string) *ValidationError {
    if len(password) < 8 {
        return &ValidationError{
            Field:   "password",
            Message: "password must be at least 8 characters",
        }
    }

    var hasUpper, hasLower, hasNumber, hasSpecial bool
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }

    if !(hasUpper && hasLower && hasNumber && hasSpecial) {
        return &ValidationError{
            Field:   "password",
            Message: "password must contain at least one uppercase letter, one lowercase letter, one number, and one special character",
        }
    }

    return nil
}

func ValidatePhone(phone string) *ValidationError {
    if !phoneRegex.MatchString(phone) {
        return &ValidationError{
            Field:   "phone",
            Message: "invalid phone number format",
        }
    }
    return nil
}