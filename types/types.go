package types

import (
	"math/big"
)

// Хранит номер блока в 16-ричной и 10-ричной системах.
type BlockNumber struct {
	Hex string `json:"result"`
	Int *big.Int
}

// Представляет транзакцию с адресами отправителя и получателя.
type Transaction struct {
	From string `json:"from"`
	To   string `json:"to"`
}
