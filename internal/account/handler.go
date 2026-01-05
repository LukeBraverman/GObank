package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/accounts", h.OpenAccount)
	r.GET("/accounts/:accountNumber", h.GetAccount)
	r.POST("/accounts/transfer", h.Transfer)
	r.GET("/accounts/:accountNumber/ledger", h.GetLedgerEntries)
	r.POST("/accounts/deposit", h.Deposit)
	r.POST("/accounts/withdraw", h.Withdraw)

}

/*
   ===== Request DTOs =====
*/

// OpenAccountRequest represents the JSON payload for opening an account
type OpenAccountRequest struct {
	AccountNumber string `json:"accountNumber"`
	Name          string `json:"name"`
}

type TransferRequest struct {
	IdempotencyKey string  `json:"idempotencyKey"`
	From           string  `json:"from"`
	To             string  `json:"to"`
	Amount         float64 `json:"amount"`
}

type DepositRequest struct {
	IdempotencyKey string  `json:"idempotencyKey"`
	AccountNumber  string  `json:"accountNumber"`
	Amount         float64 `json:"amount"`
}

type WithdrawRequest struct {
	IdempotencyKey string  `json:"idempotencyKey"`
	AccountNumber  string  `json:"accountNumber"`
	Amount         float64 `json:"amount"`
}

/*
   ===== Handlers =====
*/

// OpenAccount creates a new account from request data
func (h *Handler) OpenAccount(c *gin.Context) {
	var req OpenAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	acc, err := h.service.OpenAccount(
		req.AccountNumber,
		req.Name,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func (h *Handler) Transfer(c *gin.Context) {
	var req TransferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.Transfer(
		req.IdempotencyKey,
		req.From,
		req.To,
		req.Amount,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "transfer complete"})
}

func (h *Handler) GetLedgerEntries(c *gin.Context) {
	accountNumber := c.Param("accountNumber")

	entries, err := h.service.GetLedgerEntries(accountNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func (h *Handler) Withdraw(c *gin.Context) {
	var req WithdrawRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.Withdraw(
		req.IdempotencyKey,
		req.AccountNumber,
		req.Amount,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *Handler) Deposit(c *gin.Context) {
	var req DepositRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.Deposit(
		req.IdempotencyKey,
		req.AccountNumber,
		req.Amount,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}
