package main

import (
	"log"
	"net/http"

	"github.com/LukeBraverman/GObank/internal/account"
	"github.com/LukeBraverman/GObank/internal/ledger"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// --- Ledger creation --
	//ledgerRepo := ledger.NewRepository()
	ledgerRepo, err := ledger.NewSQLiteRepository("ledger.db")
	if err != nil {
		log.Fatalf("failed to initialize ledger repository: %v", err)
	}
	ledgerSvc := ledger.NewService(ledgerRepo)
	// --- Account wiring ---
	accountRepo, err := account.NewSQLiteRepository("accounts.db")
	if err != nil {
		log.Fatalf("failed to initialize account repository: %v", err)
	}
	svc := account.NewService(accountRepo, ledgerSvc)
	handler := account.NewHandler(svc)

	handler.RegisterRoutes(r)
	// ----------------------

	// svc.OpenAccount("alice","Alice");

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
