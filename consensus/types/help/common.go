package help

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/crypto/sha3"
	"github.com/kardiachain/go-kardia/lib/rlp"
)

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func Hash256Byte(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	hashx := sha256.New()
	hashx.Write(data)
	return hashx.Sum(nil), nil
}

func PanicSanity(v interface{}) {
	panic(fmt.Sprintf("Panicked on a Sanity Check: %v", v))
}
