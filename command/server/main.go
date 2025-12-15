package main

import (
    "github.com/LukeBraverman/GObank/internal/account"
    "net/http"

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

    // --- Account wiring ---
    repo := account.NewRepository()
    svc := account.NewService(repo)
    handler := account.NewHandler(svc)

    handler.RegisterRoutes(r)
    // ----------------------

    r.Run(":8081")
}
