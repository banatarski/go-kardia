package types

import "math/big"

type MatchOrderInput struct {
	FromType    string
	ToType      string
	FromAddress string
	Receiver    string
	Txid        string
	Amount      *big.Int
	Timestamp   *big.Int
}
