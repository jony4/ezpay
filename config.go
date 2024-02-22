package ezpay

type Config struct {
	ProjectID     string `json:"project_id"`
	ProjectSecret string `json:"project_secret"`
	PaywallID     string `json:"paywall_id"`
}
