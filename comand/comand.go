package handler

import (
	"connection/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	chatService *service.GigaChatService
}

func NewHandler(chatService *service.GigaChatService) *Handler {
	return &Handler{chatService: chatService}
}

func (h *Handler) HandleGigaChatRequest(c echo.Context) error {
	requestStr := c.QueryParam("answer")
	rquid := uuid.New().String()

	response, err := h.chatService.SendRequestAndGetResponse(requestStr, rquid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка обработки запроса"})
	}

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("X-Request-ID", rquid)
	return c.JSON(http.StatusOK, map[string]string{"message": response})
}
