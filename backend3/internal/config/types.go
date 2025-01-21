// internal/config/types.go
package config

import (
    "fmt"
    "time"
)

type ServerConfig struct {
    Port         string
    Environment  string
    LogLevel     string
    FrontendURL  string
    RateLimit    RateLimitConfig
}

type RateLimitConfig struct {
    Enabled  bool
    Requests int
    Duration time.Duration
}

type DatabaseConfig struct {
    Host              string
    Port              string
    Name              string
    User              string
    Password          string
    SSLMode           string
    MaxConnections    int
    MaxIdleConnections int
}

func (c DatabaseConfig) DSN() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
    )
}

type RedisConfig struct {
    Host     string
    Port     string
    Password string
    DB       int
}

func (c RedisConfig) Address() string {
    return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type OdooConfig struct {
    URL      string
    Database string
    Username string
    Password string
    APIKey   string
}

type AdyenConfig struct {
    APIKey           string
    MerchantAccount  string
    Environment      string
    WebhookHMACKey   string
    ClientKey        string
}

type JWTConfig struct {
    Secret           string
    ExpirationHours  int
}

type CORSConfig struct {
    AllowedOrigins []string
    AllowedMethods []string
    AllowedHeaders []string
}

type AWSConfig struct {
    Region          string
    AccessKeyID     string
    SecretAccessKey string
    BucketName      string
}