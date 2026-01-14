package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/guidiguidi/RateMonitorBC/internal/bestchange"
    "github.com/guidiguidi/RateMonitorBC/config"
    "github.com/guidiguidi/RateMonitorBC/internal/httpapi"
)

func main() {
    // Загрузка конфига
    if err := config.Load(); err != nil {
        log.Fatal("Config load failed:", err)
    }

    // Инициализация зависимостей
    bc := bestchange.NewClient(config.Cfg.BestChange.APIKey)
    h := httpapi.NewHandler(bc)

    // Роутер
    r := gin.Default()
    r.GET("/health", health)
    r.POST("/api/best-exchange", h.GetBestExchange)

    // Сервер
    srv := &http.Server{
        Addr:         ":" + config.Cfg.Server.Port,
        Handler:      r,
        ReadTimeout:  time.Duration(config.Cfg.Server.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(config.Cfg.Server.WriteTimeout) * time.Second,
    }

    // Graceful shutdown
    go func() {
        log.Printf("Server starting on :%s", config.Cfg.Server.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    log.Println("Server stopped")
}

func health(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  "ok",
        "version": "0.1.0",
        "config": gin.H{
            "port":        config.Cfg.Server.Port,
            "readTimeout": config.Cfg.Server.ReadTimeout,
        },
    })
}
