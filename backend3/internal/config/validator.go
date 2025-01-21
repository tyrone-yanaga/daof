// internal/config/validator.go
package config

import (
    "fmt"
    "strings"
)

func (c *Config) Validate() error {
    var errors []string

    // Validate Server Config
    if c.Server.Port == "" {
        errors = append(errors, "server port is required")
    }

    if c.Server.Environment == "" {
        errors = append(errors, "environment is required")
    }

    // Validate Database Config
    if c.Database.Host == "" {
        errors = append(errors, "database host is required")
    }

    if c.Database.Name == "" {
        errors = append(errors, "database name is required")
    }

    if c.Database.User == "" {
        errors = append(errors, "database user is required")
    }

    // Validate Redis Config if enabled
    if c.Redis.Host == "" {
        errors = append(errors, "redis host is required")
    }

    // Validate Odoo Config
    if c.Odoo.URL == "" {
        errors = append(errors, "odoo URL is required")
    }

    if c.Odoo.APIKey == "" {
        errors = append(errors, "odoo API key is required")
    }

    // Validate JWT Config
    if c.JWT.Secret == "" {
        errors = append(errors, "JWT secret is required")
    }

    if c.JWT.ExpirationHours <= 0 {
        errors = append(errors, "JWT expiration hours must be positive")
    }

    // Validate Adyen Config if in production
    if c.Server.Environment == "production" {
        if c.Adyen.APIKey == "" {
            errors = append(errors, "Adyen API key is required in production")
        }
        if c.Adyen.MerchantAccount == "" {
            errors = append(errors, "Adyen merchant account is required in production")
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("configuration validation failed: %s", strings.Join(errors, "; "))
    }

    return nil
}