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
	clientID        string
	token           string
	tokenExpiration time.Time
	baseURL         string
	chatURL         string
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

func NewGigaChatService(logger *zap.SugaredLogger, authKey, clientID string) *GigaChatService {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}
	return &GigaChatService{
		logger:     logger,
		authKey:    authKey,
		clientID:   clientID,
		baseURL:    "https://ngw.devices.sberbank.ru:9443",
		chatURL:    "https://gigachat.devices.sberbank.ru/api/v1",
		httpClient: client,
	}
}

func (g *GigaChatService) GetToken() (string, error) {
	if g.token != "" && g.tokenExpiration.After(time.Now()) {
		return g.token, nil
	}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "GIGACHAT_API_PERS")

	req, err := http.NewRequest("POST", g.baseURL+"/api/v2/oauth", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", uuid.New().String())
	auth := strings.TrimSpace(g.authKey)
	if !strings.HasPrefix(auth, "Basic ") {
		auth = "Basic " + auth
	}
	req.Header.Set("Authorization", auth)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get token: %s (%s)", resp.Status, string(body))
	}

	var tr TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}
	g.token = tr.AccessToken
	g.tokenExpiration = time.Unix(tr.ExpiresAt, 0)
	return g.token, nil
}

func (g *GigaChatService) Chat(messages []ChatMessage) (string, error) {
	token, err := g.GetToken()
	if err != nil {
		return "", err
	}
	reqBody := ChatRequest{Model: "GigaChat", Messages: messages}
	data, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", g.chatURL+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("x-client-id", g.clientID)
	req.Header.Set("x-request-id", uuid.New().String())
	req.Header.Set("x-session-id", uuid.New().String())

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("chat completion failed: %s (%s)", resp.Status, string(body))
	}

	var ch ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&ch); err != nil {
		return "", err
	}
	if len(ch.Choices) == 0 {
		return "", fmt.Errorf("no chat choices returned")
	}
	return ch.Choices[0].Message.Content, nil
}

func (g *GigaChatService) GenerateDescription(title, _ string) (string, error) {
	prompt := fmt.Sprintf("Создай краткое привлекательное описание для баннера '%s'.", title)
	return g.Chat([]ChatMessage{{Role: "user", Content: prompt}})
}
