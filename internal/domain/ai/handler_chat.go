package ai

import (
	"atlas_food/internal/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	service ChatService
}

func NewChatHandler(service ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User tidak terautentikasi")
		return
	}

	res, err := h.service.Chat(userID.(string), req)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, res)
}

func (h *ChatHandler) GetSessions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User tidak terautentikasi")
		return
	}

	res, err := h.service.GetSessions(userID.(string))
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"sessions": res})
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User tidak terautentikasi")
		return
	}
	sessionID := c.Param("session_id")

	res, err := h.service.GetMessages(userID.(string), sessionID)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"messages": res})
}

func (h *ChatHandler) DeleteSession(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User tidak terautentikasi")
		return
	}
	sessionID := c.Param("session_id")

	err := h.service.DeleteSession(userID.(string), sessionID)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Sesi berhasil dihapus"})
}

// SetupChatRoutes mendaftarkan endpoint chat ke router
func (h *ChatHandler) SetupChatRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	chat := router.Group("/ai/chat", authMiddleware)
	{
		chat.POST("", h.Chat)
		chat.GET("/sessions", h.GetSessions)
		chat.GET("/sessions/:session_id", h.GetMessages)
		chat.DELETE("/sessions/:session_id", h.DeleteSession)
	}
}
