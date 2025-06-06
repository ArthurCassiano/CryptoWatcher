package main

import (
    "log"

   "github.com/ArthurCassiano/CryptoWatcher/internal/server"
    "github.com/ArthurCassiano/CryptoWatcher/internal/binance"
    "github.com/ArthurCassiano/CryptoWatcher/internal/cache"
)

func main() {
    // 1) Cria o cache de preços
    priceCache := cache.NewPriceCache()

    // 2) Defina os símbolos que serão monitorados
    symbols := []string{"BTCUSDT", "ETHUSDT", "BNBUSDT"}

    // 3) Inicia goroutine que conecta ao WebSocket e atualiza o cache
    go binance.StartWebSocket(priceCache, symbols)

    // 4) Inicia o servidor HTTP (bloqueante). Após isso, roda até interrupção.
    server.StartHTTPServer(priceCache, "8080")

    // Obs.: Se quiser usar porta diferente, basta trocar “8080” por outra.
    // Por exemplo, "8000" ou usar variável de ambiente.
    log.Println("Encerrando aplicação")
}
