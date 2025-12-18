package account

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type Handler struct {
    service *Service
}

// Constructor
func NewHandler(service *Service) *Handler {
    return &Handler{
        service: service,
    }
}

// RegisterRoutes wires all account-related routes
func (h *Handler) RegisterRoutes(r *gin.Engine) {
    r.POST("/accounts", h.OpenAccount)
    r.GET("/accounts/:accountNumber", h.GetAccount)
}

/*
   ===== Request DTOs =====
*/

// OpenAccountRequest represents the JSON payload for opening an account
type OpenAccountRequest struct {
    AccountNumber  string  `json:"accountNumber"`
    Name           string  `json:"name"`
    InitialBalance float64 `json:"initialBalance"`
}

/*
   ===== Handlers =====
*/

// OpenAccount creates a new account from request data
func (h *Handler) OpenAccount(c *gin.Context) {
    var req OpenAccountRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid request body",
        })
        return
    }

    acc, err := h.service.OpenAccount(
        req.AccountNumber,
        req.Name,
        req.InitialBalance,
    )
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, acc)
}

// GetAccount fetches an account by account number
func (h *Handler) GetAccount(c *gin.Context) {
    accountNumber := c.Param("accountNumber")

    acc, err := h.service.GetAccount(accountNumber)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, acc)
}
