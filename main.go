package main

import (
	"connection/comand"
	"connection/service"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка при загрузке .env файла")
		return
	}
	
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	authKey := os.Getenv("AUTH_KEY")
	fmt.Println(authKey)
	tokenURL := "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	tokenService := service.NewTokenService(authKey, tokenURL, client)

	requestURL := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	chatService := service.NewGigaChatService(tokenService, requestURL, client)

	chatHandler := handler.NewHandler(chatService)

	e.GET("/lite", chatHandler.HandleGigaChatRequest)

	e.Start(":8008")
}
