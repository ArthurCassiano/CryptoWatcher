package server

import (
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/ArthurCassiano/CryptoWatcher/internal/cache"
)

func StartHTTPServer(pc *cache.PriceCache, port string) {
    router := gin.Default()

    // Healthcheck (útil para monitoramento)
    router.GET("/health", func(c *gin.Context) {
        c.String(http.StatusOK, "OK")
    })

    // Endpoint principal: /prices
    // Retorna o JSON com todos os preços armazenados no cache
    router.GET("/prices", func(c *gin.Context) {
        prices := pc.GetAll()
        c.JSON(http.StatusOK, prices)
    })

    srv := &http.Server{
        Addr:           ":" + port,
        Handler:        router,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20, // 1 MiB
    }

    log.Printf("[HTTP] Servidor rodando em http://localhost:%s", port)
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("[HTTP] ListenAndServe falhou: %v", err)
    }
}
