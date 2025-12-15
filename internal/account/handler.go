package account

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type Handler struct {
    service *service
}

// Constructor
func NewHandler(service *service) *Handler {
    return &Handler{
        service: service,
    }
}

// Route registration
func (h *Handler) RegisterRoutes(r *gin.Engine) {
    r.GET("/accounts/test", h.GetTestAccount)
}

// HTTP handler method
func (h *Handler) GetTestAccount(c *gin.Context) {
    account := h.service.GetTestAccount()
    c.JSON(http.StatusOK, account)
}
