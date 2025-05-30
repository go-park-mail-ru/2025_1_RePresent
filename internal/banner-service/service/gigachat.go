package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GigaChatService struct {
	logger          *zap.SugaredLogger
	authKey         string
	token           string
	tokenExpiration time.Time
	baseURL         string
	httpClient      *http.Client
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
	TokenType   string `json:"token_type"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func NewGigaChatService(logger *zap.SugaredLogger, authKey string) *GigaChatService {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	return &GigaChatService{
		logger:     logger,
		authKey:    authKey,
		baseURL:    "https://ngw.devices.sberbank.ru:9443",
		httpClient: client,
	}
}

func (g *GigaChatService) GetToken() (string, error) {
	if g.token != "" && g.tokenExpiration.After(time.Now()) {
		return g.token, nil
	}
	g.logger.Debugw("Getting new token from GigaChat API")

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "GIGACHAT_API_PERS")

	req, err := http.NewRequest("POST", g.baseURL+"/api/v2/oauth", strings.NewReader(data.Encode()))
	if err != nil {
		g.logger.Errorw("Failed to create auth request", "error", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", uuid.New().String())
	authHeader := g.authKey
	if !strings.HasPrefix(strings.TrimSpace(authHeader), "Basic ") {
		authHeader = "Basic " + authHeader
	}
	req.Header.Set("Authorization", authHeader)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		g.logger.Errorw("Auth request failed", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get token: %s (%s)", resp.Status, string(body))
	}

	var trsp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&trsp); err != nil {
		return "", err
	}
	g.tokenExpiration = time.Unix(trsp.ExpiresAt, 0)
	g.token = trsp.AccessToken
	return g.token, nil
}

func (g *GigaChatService) GenerateDescription(bannerTitle, bannerContent string) (string, error) {
	token, err := g.GetToken()
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(
		"Создай краткое привлекательное описание для баннера '%s': %s",
		bannerTitle, bannerContent,
	)

	reqBody, _ := json.Marshal(ChatRequest{
		Model:    "GigaChat",
		Messages: []ChatMessage{{Role: "user", Content: prompt}},
	})

	req, err := http.NewRequest("POST", g.baseURL+"/api/v2/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("chat completion failed: %s (%s)", resp.Status, string(body))
	}

	var chResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chResp); err != nil {
		return "", err
	}
	if len(chResp.Choices) == 0 {
		return "", fmt.Errorf("no completion returned")
	}
	return chResp.Choices[0].Message.Content, nil
}
