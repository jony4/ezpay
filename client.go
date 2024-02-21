package ezpay

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nats-io/nuid"
	log "github.com/sirupsen/logrus"
)

var (
	TokenAlgorithm = "HS256"
	TokenFields    = []string{"exp", "nonce"}
	TokenExpiresIn = 30 * 60

	BaseURL            = "https://revenue.ezboti.com/api/v1/server"
	DefaultContentType = "text/plain"
)

var (
	ErrFailedResponse = errors.New("failed response")
)

type Client struct {
	ProjectID     string
	ProjectSecret string
	PaywallID     string
	BaseURL       string
}

// NewClient creates a new ezpay client
// 每个 Client 实例对应一个项目，一个项目对应一个 paywall。
// 不同的项目可以实例化不同的 Client 进行调用。
func NewClient(cfg *Config) (*Client, error) {
	return &Client{
		ProjectID:     cfg.ProjectID,
		ProjectSecret: cfg.ProjectSecret,
		PaywallID:     cfg.PaywallID,
		BaseURL:       BaseURL,
	}, nil
}

func (c *Client) decodeToken(token string) (jwt.MapClaims, error) {
	var claims jwt.MapClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		hmac, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if hmac != jwt.SigningMethodHS256 {
			log.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.ProjectSecret), nil
	})
	if err != nil {
		log.Errorf("ParseWithClaims Error, err=%+v", err)
		return nil, err
	}

	for _, field := range TokenFields {
		if _, ok := claims[field]; !ok {
			log.Errorf("token did not contain %s field", field)
			return nil, fmt.Errorf("token did not contain %s field", field)
		}
	}

	return claims, nil
}

func (c *Client) encodeToken(payload map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"exp":   time.Now().Unix() + int64(TokenExpiresIn),
		"nonce": nuid.Next(),
	}
	for k, v := range payload {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["project_id"] = c.ProjectID

	signed, err := token.SignedString([]byte(c.ProjectSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (c *Client) Call(ctx context.Context, api string, params interface{}) (interface{}, error) {
	payload := map[string]interface{}{
		"method": api,
		"params": params,
	}
	content, _ := c.encodeToken(payload)

	url := c.BaseURL + "/" + api
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", DefaultContentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Do Error, err=%+v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("ReadAll Error, err=%+v", err)
		return nil, err
	}

	respBody, err := c.decodeToken(string(body))
	if err != nil {
		return nil, err
	}

	if err := respBody.Valid(); err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		return respBody["result"], ErrFailedResponse
	}

	return respBody["result"], nil
}
