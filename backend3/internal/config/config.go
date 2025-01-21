// internal/config/config.go
package config

import (
    "fmt"
    "time"

    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Odoo     OdooConfig
    Adyen    AdyenConfig
    JWT      JWTConfig
    CORS     CORSConfig
    AWS      AWSConfig
}

func New() (*Config, error) {
    viper.SetConfigName(".env")
    viper.SetConfigType("env")
    viper.AddConfigPath(".")
    
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, fmt.Errorf("failed to read config file: %w", err)
        }
    }

    config := &Config{}
    
    // Server Configuration
    config.Server = ServerConfig{
        Port:         viper.GetString("PORT"),
        Environment:  viper.GetString("ENV"),
        LogLevel:     viper.GetString("LOG_LEVEL"),
        FrontendURL:  viper.GetString("FRONTEND_URL"),
        RateLimit: RateLimitConfig{
            Enabled:  viper.GetBool("RATE_LIMIT_ENABLED"),
            Requests: viper.GetInt("RATE_LIMIT_REQUESTS"),
            Duration: viper.GetDuration("RATE_LIMIT_DURATION"),
        },
    }

    // Database Configuration
    config.Database = DatabaseConfig{
        Host:              viper.GetString("ODOO_DB_HOST"),
        Port:              viper.GetString("ODOO_DB_PORT"),
        Name:              viper.GetString("ODOO_DB_NAME"),
        User:              viper.GetString("ODOO_DB_USER"),
        Password:          viper.GetString("ODOO_DB_PASSWORD"),
        SSLMode:           viper.GetString("DB_SSL_MODE"),
        MaxConnections:    viper.GetInt("DB_MAX_CONNECTIONS"),
        MaxIdleConnections: viper.GetInt("DB_MAX_IDLE_CONNECTIONS"),
    }

    // Redis Configuration
    config.Redis = RedisConfig{
        Host:     viper.GetString("REDIS_HOST"),
        Port:     viper.GetString("REDIS_PORT"),
        Password: viper.GetString("REDIS_PASSWORD"),
        DB:       viper.GetInt("REDIS_DB"),
    }

    // Odoo Configuration
    config.Odoo = OdooConfig{
        URL:      viper.GetString("ODOO_URL"),
        Database: viper.GetString("ODOO_DB_NAME"),
        Username: viper.GetString("ODOO_USERNAME"),
        Password: viper.GetString("ODOO_PASSWORD"),
        APIKey:   viper.GetString("ODOO_API_KEY"),
    }

    // Adyen Configuration
    config.Adyen = AdyenConfig{
        APIKey:           viper.GetString("ADYEN_API_KEY"),
        MerchantAccount:  viper.GetString("ADYEN_MERCHANT_ACCOUNT"),
        Environment:      viper.GetString("ADYEN_ENVIRONMENT"),
        WebhookHMACKey:   viper.GetString("ADYEN_WEBHOOK_HMAC_KEY"),
        ClientKey:        viper.GetString("ADYEN_CLIENT_KEY"),
    }

    // JWT Configuration
    config.JWT = JWTConfig{
        Secret:           viper.GetString("JWT_SECRET"),
        ExpirationHours:  viper.GetInt("JWT_EXPIRATION_HOURS"),
    }

    // CORS Configuration
    config.CORS = CORSConfig{
        AllowedOrigins:  viper.GetStringSlice("ALLOWED_ORIGINS"),
        AllowedMethods:  viper.GetStringSlice("ALLOWED_METHODS"),
        AllowedHeaders:  viper.GetStringSlice("ALLOWED_HEADERS"),
    }

    // AWS Configuration
    config.AWS = AWSConfig{
        Region:          viper.GetString("AWS_REGION"),
        AccessKeyID:     viper.GetString("AWS_ACCESS_KEY_ID"),
        SecretAccessKey: viper.GetString("AWS_SECRET_ACCESS_KEY"),
        BucketName:      viper.GetString("AWS_BUCKET_NAME"),
    }

    return config, nil
}