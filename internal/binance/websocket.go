package binance

import (
    "encoding/json"
    "fmt"
    "log"
    "strings"
    "time"

    "github.com/gorilla/websocket"
	"github.com/ArthurCassiano/CryptoWatcher/internal/cache"
	"github.com/ArthurCassiano/CryptoWatcher/pkg/model"
)

// URL base do WebSocket da Binance
const baseURL = "wss://fstream.binance.com"

// StartWebSocket conecta ao WebSocket da Binance para receber preços em real time.
// - pc: cache que irá armazenar os valores.
// - symbols: lista de símbolos (ex.: []{"BTCUSDT","ETHUSDT"}).
func StartWebSocket(pc *cache.PriceCache, symbols []string) {
    // 1) Monta a rota do stream: streams=btcusdt@ticker/ethusdt@ticker/ ...
    var streams []string
    for _, sym := range symbols {
        streams = append(streams, strings.ToLower(sym)+"@ticker")
    }
    url := fmt.Sprintf("%s/stream?streams=%s", baseURL, strings.Join(streams, "/"))

    for {
        // 2) Conecta ao WebSocket
        ws, _, err := websocket.DefaultDialer.Dial(url, nil)
        if err != nil {
            log.Printf("[WebSocket] Falha ao conectar: %v. Tentando reconectar em 5 segundos...", err)
            time.Sleep(5 * time.Second)
            continue
        }
        log.Println("[WebSocket] Conectado à Binance:", url)

        // 3) Entra no loop de leitura de mensagens
        for {
            _, msg, err := ws.ReadMessage()
            if err != nil {
                log.Printf("[WebSocket] Erro na leitura: %v. Reconectando em 3 segundos...", err)
                ws.Close()
                time.Sleep(3 * time.Second)
                break // sai do loop interno e volta a tentar reconectar externamente
            }

            // 4) O JSON vem no formato:
            // {
            //   "stream": "btcusdt@ticker",
            //   "data": {
            //       "s": "BTCUSDT",
            //       "c": "69250.12",
            //       // ... outros campos (volume, variação, etc.) ...
            //   }
            // }
            var payload struct {
                Stream string          `json:"stream"`
                Data   json.RawMessage `json:"data"`
            }
            if err := json.Unmarshal(msg, &payload); err != nil {
                log.Printf("[WebSocket] Erro ao parsear payload externo: %v", err)
                continue
            }

            // 5) Dentro de Data, queremos apenas "s" e "c"
            var ticker struct {
                Symbol string `json:"s"` // ex.: "BTCUSDT"
                Price  string `json:"p"` // último preço, como string
            }
            log.Printf("[WebSocket] JSON recebido em payload.Data: %s", string(payload.Data))
            if err := json.Unmarshal(payload.Data, &ticker); err != nil {
                log.Printf("[WebSocket] Erro ao parsear dados do ticker: %v", err)
                continue
            }

            // 6) Converte string para float64
            var priceFloat float64
            if _, err := fmt.Sscan(ticker.Price, &priceFloat); err != nil {
                log.Printf("[WebSocket] Erro ao converter preço: %v", err)
                continue
            }

            // 7) Atualiza o cache com novo valor e timestamp
            pc.Update(model.Price{
                Symbol:    ticker.Symbol,
                Price:     priceFloat,
                Timestamp: time.Now().Unix(),
            })
        }
    }
}
