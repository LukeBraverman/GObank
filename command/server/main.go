package main

import (
    "github.com/LukeBraverman/GObank/internal/account"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/LukeBraverman/GObank/internal/ledger"

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
        panic(err) // or log.Fatal(err)
    }
    ledgerSvc := ledger.NewService(ledgerRepo)
    // --- Account wiring ---
    repo := account.NewRepository()
    svc := account.NewService(repo, ledgerSvc)
    handler := account.NewHandler(svc)

    handler.RegisterRoutes(r)
    // ----------------------

    // svc.OpenAccount("alice","Alice");


    r.Run(":8081")
}
