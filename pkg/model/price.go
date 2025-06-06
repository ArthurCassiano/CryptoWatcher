package model

// Price representa o preço de uma criptomoeda em um dado momento.
type Price struct {
    Symbol    string  `json:"symbol"`    // Ex.: "BTCUSDT"
    Price     float64 `json:"price"`     // Último preço
    Timestamp int64   `json:"timestamp"` // Unix timestamp em segundos
}
