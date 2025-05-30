package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

type GigaChatService struct {
	logger          *zap.SugaredLogger
	authKey         string
	token           string
	tokenExpiration time.Time
	baseURL         string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
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
	return &GigaChatService{
		logger:  logger,
		authKey: authKey,
		baseURL: "https://ngw.devices.sberbank.ru:9443",
	}
}

func (g *GigaChatService) GetToken() (string, error) {
	if g.token != "" && g.tokenExpiration.After(time.Now()) {
		return g.token, nil
	}

	g.logger.Debugw("Getting new token from GigaChat API")

	data := url.Values{}
	data.Set("scope", "GIGACHAT_API_PERS")

	req, err := http.NewRequest("POST", g.baseURL+"/api/v2/oauth", strings.NewReader(data.Encode()))
	if err != nil {
		g.logger.Errorw("Failed to create auth request", "error", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", fmt.Sprintf("%d", time.Now().UnixNano()))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", g.authKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		g.logger.Errorw("Failed to execute auth request", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		g.logger.Errorw("Failed to get token", "status", resp.Status, "body", string(body))
		return "", fmt.Errorf("failed to get token: %s", resp.Status)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		g.logger.Errorw("Failed to decode token response", "error", err)
		return "", err
	}

	g.token = tokenResp.AccessToken
	g.tokenExpiration = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	g.logger.Debugw("Successfully got token from GigaChat API", "expires_in", tokenResp.ExpiresIn)

	return g.token, nil
}

func (g *GigaChatService) GenerateDescription(bannerTitle string, bannerContent string) (string, error) {
	token, err := g.GetToken()
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf("Создай краткое и привлекательное описание для баннера с названием '%s' и содержанием '%s'. Описание должно быть не более 200 символов и мотивировать пользователя кликнуть по баннеру.", bannerTitle, bannerContent)

	chatRequest := ChatRequest{
		Model: "GigaChat",
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(chatRequest)
	if err != nil {
		g.logger.Errorw("Failed to marshal chat request", "error", err)
		return "", err
	}

	req, err := http.NewRequest("POST", g.baseURL+"/api/v2/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		g.logger.Errorw("Failed to create completion request", "error", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		g.logger.Errorw("Failed to execute completion request", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		g.logger.Errorw("Failed to get completion", "status", resp.Status, "body", string(body))
		return "", fmt.Errorf("failed to get completion: %s", resp.Status)
	}

	var chatResponse ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		g.logger.Errorw("Failed to decode completion response", "error", err)
		return "", err
	}

	if len(chatResponse.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	return chatResponse.Choices[0].Message.Content, nil
}
