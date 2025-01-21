// internal/config/odoo.go
type OdooConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	APIKey   string
	URL      string
}

// Get Odoo database URL
func (c *OdooConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		c.Host, c.Port, c.Username, c.Password, c.Database)
}