// internal/utils/hash.go
package utils

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "fmt"

    "golang.org/x/crypto/bcrypt"
)

const (
    defaultCost = 12
    saltLength  = 16
)

// HashPassword creates a bcrypt hash of a password
func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    return string(hash), nil
}

// ComparePassword compares a password with a hash
func ComparePassword(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateToken generates a random token
func GenerateToken(length int) (string, error) {
    b := make([]byte, length)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(b), nil
}

// GenerateHash generates a SHA-256 hash of the input
func GenerateHash(input string) string {
    hash := sha256.New()
    hash.Write([]byte(input))
    return fmt.Sprintf("%x", hash.Sum(nil))
}

// GenerateSalt generates a random salt
func GenerateSalt() (string, error) {
    salt := make([]byte, saltLength)
    _, err := rand.Read(salt)
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(salt), nil
}