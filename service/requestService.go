package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type GigaChatService struct {
	client           *http.Client
	tokenFromService *TokenService
	requestURL       string
	method           string
}

func NewGigaChatService(token *TokenService, requestURL string, client *http.Client) *GigaChatService {
	return &GigaChatService{
		client:           client,
		tokenFromService: token,
		requestURL:       requestURL,
		method:           "POST",
	}
}

type message struct {
	Content string `json:"content"`
}

type choice struct {
	Message message `json:"message"`
}

type gigaChatResponse struct {
	Choices []choice `json:"choices"`
}

func createRequestBody(userMessage string) []byte {
	requestBody := map[string]interface{}{
		"model":    "GigaChat",
		"messages": []map[string]string{{"role": "user", "content": userMessage}},
		"stream":   false,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil
	}

	return body
}

func (g *GigaChatService) SendRequestAndGetResponse(userMessage, reqUid string) (responseFromGigaChat string, err error) {
	err = g.tokenFromService.GetAccessToken()
	if err != nil {
		return
	}

	body := createRequestBody(userMessage)
	if body == nil {
		return "", errors.New("ошибка при формировании тела запроса")
	}

	req, err := http.NewRequest(g.method, g.requestURL, bytes.NewReader(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", reqUid)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.tokenFromService.token.AccessToken)

	res, err := g.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var chatResponse gigaChatResponse
	if err = json.Unmarshal(responseBody, &chatResponse); err != nil {
		return
	}

	if len(chatResponse.Choices) > 0 {
		return chatResponse.Choices[0].Message.Content, nil
	}

	return "", errors.New("нет ответа от GigaChat")
}
