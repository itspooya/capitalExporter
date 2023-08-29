package account

import (
	"bytes"
	"capitalExporter/config"
	"capitalExporter/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	nonDemoBaseUrl = "https://api-capital.backend-capital.com/"
	demoBaseUrl    = "https://demo-api-capital.backend-capital.com/"
)

var baseUrl string
var email string
var password string
var APIKEY string
var demo bool
var debug bool
var sugar = logger.InitLogger(debug)

func init() {
	// Set the base URL based on the demo flag
	setBaseUrl()
}

func setBaseUrl() {
	if demo {
		baseUrl = demoBaseUrl
	} else {
		baseUrl = nonDemoBaseUrl
	}
}

type SessionResponse struct {
	CST            string
	XSecurityToken string
	ExpiryTime     time.Time
}

type sessionJson struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type BalanceInfo struct {
	Balance    float64 `json:"balance"`
	Deposit    float64 `json:"deposit"`
	ProfitLoss float64 `json:"profitLoss"`
	Available  float64 `json:"available"`
}

type Account struct {
	AccountID   string      `json:"accountId"`
	AccountName string      `json:"accountName"`
	Status      string      `json:"status"`
	AccountType string      `json:"accountType"`
	Preferred   bool        `json:"preferred"`
	Balance     BalanceInfo `json:"balance"`
	Currency    string      `json:"currency"`
}

type DetailsResponse struct {
	Accounts []Account `json:"accounts"`
}

func GenerateSession(cfg config.Configuration) (SessionResponse, error) {
	email = cfg.Email
	password = cfg.Password
	APIKEY = cfg.APIKey
	demo = cfg.Demo
	url := baseUrl + "api/v1/session"
	payload := sessionJson{
		Identifier: email,
		Password:   password,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		sugar.Errorw("error marshaling payload", "error", err)
		return SessionResponse{}, fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		sugar.Errorw("error creating request", "error", err)
		return SessionResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-CAP-API-KEY", APIKEY)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second} // Add a timeout
	resp, err := client.Do(req)
	if err != nil {
		sugar.Errorw("error sending request", "error", err)
		return SessionResponse{}, fmt.Errorf("error sending request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			sugar.Fatalf("Error closing response body: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		sugar.Errorw("unexpected status code", "status", resp.StatusCode)
		return SessionResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	session := SessionResponse{
		CST:            resp.Header.Get("CST"),
		XSecurityToken: resp.Header.Get("X-SECURITY-TOKEN"),
		ExpiryTime:     time.Now().Add(10 * time.Minute),
	}

	return session, nil
}

func GetDetails(CST, XSecurityToken string) ([]Account, error) {
	url := baseUrl + "api/v1/accounts"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		sugar.Errorw("error creating request", "error", err)
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("CST", CST)
	req.Header.Set("X-SECURITY-TOKEN", XSecurityToken)

	client := &http.Client{Timeout: 10 * time.Second} // Add a timeout
	resp, err := client.Do(req)
	if err != nil {
		sugar.Errorw("error sending request", "error", err)
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			sugar.Fatalf("Error closing response body: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		sugar.Errorw("unexpected status code", "status", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		sugar.Errorw("error reading response body", "error", err)
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var details DetailsResponse
	if err := json.Unmarshal(body, &details); err != nil {
		sugar.Errorw("error unmarshalling JSON", "error", err)
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return details.Accounts, nil
}

func RefreshTokenIfNeeded(sessionTokens SessionResponse) (SessionResponse, error) {
	// Check if the tokens are close to expiring (e.g., within 1 minute of expiry).
	if time.Now().Add(1 * time.Minute).After(sessionTokens.ExpiryTime) {
		newTokens, err := GenerateSession(config.Configuration{
			Email:    email,
			Password: password,
		},
		)
		if err != nil {
			sugar.Errorw("error generating new session", "error", err)
			return SessionResponse{}, err
		}
		return newTokens, nil
	}
	return sessionTokens, nil
}
