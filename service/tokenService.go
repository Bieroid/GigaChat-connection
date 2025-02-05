package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type TokenService struct {
	authKey   string
	token     TokenResponse
	client    *http.Client
	tokenURL  string
	tokenBody []byte
	mu        sync.Mutex
}

func NewTokenService(authKey, tokenURL string, client *http.Client) *TokenService {
	return &TokenService{
		authKey:   authKey,
		client:    client,
		tokenURL:  tokenURL,
		tokenBody: []byte(`scope=GIGACHAT_API_PERS`),
		token:     TokenResponse{AccessToken: "", ExpiresAt: 2},
	}
}

func (t *TokenService) CreateNewToken() (err error) {
	req, err := http.NewRequest("POST", t.tokenURL, bytes.NewReader(t.tokenBody))
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Basic "+t.authKey)
	req.Header.Set("RqUID", uuid.New().String())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("неверный статус ответа")
	}

	if err = json.NewDecoder(resp.Body).Decode(&t.token); err != nil {
		return
	}

	return
}

func (t *TokenService) GetAccessToken() (err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.token.ExpiresAt > time.Now().Unix() {
		return
	}

	err = t.CreateNewToken()
	if err != nil {
		return errors.New("ошибка при получении нового токена")
	}

	return
}
