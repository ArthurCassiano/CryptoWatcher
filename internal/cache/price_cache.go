package cache

import (
    "sync"
    "time"

    "github.com/ArthurCassiano/CryptoWatcher/pkg/model"
)

// Defina aqui os símbolos que sua aplicação vai monitorar.
// Isso deve bater exatamente com o que você passará para o WebSocket mais adiante.
var symbols = []string{"BTCUSDT", "ETHUSDT", "BNBUSDT"}

type PriceCache struct {
    mu   sync.RWMutex
    data map[string]model.Price
}

// NewPriceCache inicializa o cache com cada símbolo em price=0 e timestamp atual.
func NewPriceCache() *PriceCache {
    pc := &PriceCache{
        data: make(map[string]model.Price),
    }
    for _, sym := range symbols {
        pc.data[sym] = model.Price{
            Symbol:    sym,
            Price:     0,
            Timestamp: time.Now().Unix(),
        }
    }
    return pc
}

// Update atualiza o cache para um símbolo específico.
// Usa Lock() para escrita, garantindo thread-safety.
func (pc *PriceCache) Update(p model.Price) {
    pc.mu.Lock()
    defer pc.mu.Unlock()
    pc.data[p.Symbol] = p
}

// GetAll retorna uma cópia slice de todos os preços atuais.
// Usa RLock() para permitir leituras concorrentes.
func (pc *PriceCache) GetAll() []model.Price {
    pc.mu.RLock()
    defer pc.mu.RUnlock()
    out := make([]model.Price, 0, len(pc.data))
    for _, v := range pc.data {
        out = append(out, v)
    }
    return out
}

// Get retorna o preço para um símbolo específico, ou false se não existir.
func (pc *PriceCache) Get(symbol string) (model.Price, bool) {
    pc.mu.RLock()
    defer pc.mu.RUnlock()
    p, ok := pc.data[symbol]
    return p, ok
}
