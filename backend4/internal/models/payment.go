package models

type PaymentData struct {
	SessionData string                 `json:"sessionData"`
	ClientKey   string                 `json:"clientKey"`
	Config      map[string]interface{} `json:"config"`
}
