// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package math provides integer math utilities.
package common

import (
	"fmt"
	"math/big"
)

// Various big integer limit values.
var (
	tt255     = BigPow(2, 255)
	tt256     = BigPow(2, 256)
	tt256m1   = new(big.Int).Sub(tt256, big.NewInt(1))
	tt63      = BigPow(2, 63)
	MaxBig256 = new(big.Int).Set(tt256m1)
	MaxBig63  = new(big.Int).Sub(tt63, big.NewInt(1))
)

const (
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes = wordBits / 8
)

// HexOrDecimal256 marshals big.Int as hex or decimal.
type HexOrDecimal256 big.Int

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *HexOrDecimal256) UnmarshalText(input []byte) error {
	bigint, ok := ParseBig256(string(input))
	if !ok {
		return fmt.Errorf("invalid hex or decimal integer %q", input)
	}
	*i = HexOrDecimal256(*bigint)
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (i *HexOrDecimal256) MarshalText() ([]byte, error) {
	if i == nil {
		return []byte("0x0"), nil
	}
	return []byte(fmt.Sprintf("%#x", (*big.Int)(i))), nil
}

// ParseBig256 parses s as a 256 bit integer in decimal or hexadecimal syntax.
// Leading zeros are accepted. The empty string parses as zero.
func ParseBig256(s string) (*big.Int, bool) {
	if s == "" {
		return new(big.Int), true
	}
	var bigint *big.Int
	var ok bool
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		bigint, ok = new(big.Int).SetString(s[2:], 16)
	} else {
		bigint, ok = new(big.Int).SetString(s, 10)
	}
	if ok && bigint.BitLen() > 256 {
		bigint, ok = nil, false
	}
	return bigint, ok
}

// MustParseBig256 parses s as a 256 bit big integer and panics if the string is invalid.
func MustParseBig256(s string) *big.Int {
	v, ok := ParseBig256(s)
	if !ok {
		panic("invalid 256 bit integer: " + s)
	}
	return v
}

// BigPow returns a ** b as a big integer.
func BigPow(a, b int64) *big.Int {
	r := big.NewInt(a)
	return r.Exp(r, big.NewInt(b), nil)
}

// BigMax returns the larger of x or y.
func BigMax(x, y *big.Int) *big.Int {
	if x.Cmp(y) < 0 {
		return y
	}
	return x
}

// BigMin returns the smaller of x or y.
func BigMin(x, y *big.Int) *big.Int {
	if x.Cmp(y) > 0 {
		return y
	}
	return x
}

// FirstBitSet returns the index of the first 1 bit in v, counting from LSB.
func FirstBitSet(v *big.Int) int {
	for i := 0; i < v.BitLen(); i++ {
		if v.Bit(i) > 0 {
			return i
		}
	}
	return v.BitLen()
}

// PaddedBigBytes encodes a big integer as a big-endian byte slice. The length
// of the slice is at least n bytes.
func PaddedBigBytes(bigint *big.Int, n int) []byte {
	if bigint.BitLen()/8 >= n {
		return bigint.Bytes()
	}
	ret := make([]byte, n)
	ReadBits(bigint, ret)
	return ret
}

// bigEndianByteAt returns the byte at position n,
// in Big-Endian encoding
// So n==0 returns the least significant byte
func bigEndianByteAt(bigint *big.Int, n int) byte {
	words := bigint.Bits()
	// Check word-bucket the byte will reside in
	i := n / wordBytes
	if i >= len(words) {
		return byte(0)
	}
	word := words[i]
	// Offset of the byte
	shift := 8 * uint(n%wordBytes)

	return byte(word >> shift)
}

// Byte returns the byte at position n,
// with the supplied padlength in Little-Endian encoding.
// n==0 returns the MSB
// Example: bigint '5', padlength 32, n=31 => 5
func Byte(bigint *big.Int, padlength, n int) byte {
	if n >= padlength {
		return byte(0)
	}
	return bigEndianByteAt(bigint, padlength-1-n)
}

// ReadBits encodes the absolute value of bigint as big-endian bytes. Callers must ensure
// that buf has enough space. If buf is too short the result will be incomplete.
func ReadBits(bigint *big.Int, buf []byte) {
	i := len(buf)
	for _, d := range bigint.Bits() {
		for j := 0; j < wordBytes && i > 0; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}
}

// U256 encodes as a 256 bit two's complement number. This operation is destructive.
func U256(x *big.Int) *big.Int {
	return x.And(x, tt256m1)
}

// S256 interprets x as a two's complement number.
// x must not exceed 256 bits (the result is undefined if it does) and is not modified.
//
//   S256(0)        = 0
//   S256(1)        = 1
//   S256(2**255)   = -2**255
//   S256(2**256-1) = -1
func S256(x *big.Int) *big.Int {
	if x.Cmp(tt255) < 0 {
		return x
	}
	return new(big.Int).Sub(x, tt256)
}

// Exp implements exponentiation by squaring.
// Exp returns a newly-allocated big integer and does not change
// base or exponent. The result is truncated to 256 bits.
//
// Courtesy @karalabe and @chfast
func Exp(base, exponent *big.Int) *big.Int {
	result := big.NewInt(1)

	for _, word := range exponent.Bits() {
		for i := 0; i < wordBits; i++ {
			if word&1 == 1 {
				U256(result.Mul(result, base))
			}
			U256(base.Mul(base, base))
			word >>= 1
		}
	}
	return result
}

type BigInt struct {
	Value uint64 `json:"value"`
	Pos   bool   `json:"pos"`
}

func NewBigInt64(x int64) *BigInt {
	if x < 0 {
		return &BigInt{Value: uint64(-x), Pos: false}
	}

	return &BigInt{Value: uint64(x), Pos: true}
}

func NewBigUint64(x uint64) *BigInt {
	return NewBigInt64(int64(x))
}

func NewBigInt32(x int) *BigInt {
	return NewBigInt64(int64(x))
}

// IsGreaterThan returns true if x is greater than y
func (x *BigInt) IsGreaterThan(y *BigInt) bool {
	return x.Int64() > y.Int64()
}

// IsGreaterThanInt returns true if x is greater than y
func (x *BigInt) IsGreaterThanInt(y int) bool {
	return x.Int64() > int64(y)
}

// IsGreaterThanInt returns true if x is greater than y
func (x *BigInt) IsGreaterThanInt64(y int64) bool {
	return x.Int64() > y
}

// IsGreaterOrEqualThanInt returns true if x is greater than or equals to y
func (x *BigInt) IsGreaterThanOrEqualToInt(y int) bool {
	return x.Int64() >= int64(y)
}

// IsGreaterOrEqualThanInt64 returns true if x is greater than or equals to y
func (x *BigInt) IsGreaterThanOrEqualToInt64(y int64) bool {
	return x.Int64() >= y
}

// IsLessThan returns true if x is less than y
func (x *BigInt) IsLessThan(y *BigInt) bool {
	return x.Int64() < y.Int64()
}

// IsLessThan returns true if x is less than y
func (x *BigInt) IsLessThanInt(y int) bool {
	return x.Int32() < y
}

// IsLessThan returns true if x is less than y
func (x *BigInt) IsLessThanInt64(y int64) bool {
	return x.Int64() < y
}

// IsLessThan returns true if x is less than y
func (x *BigInt) IsLessThanOrEquals(y *BigInt) bool {
	return x.Int64() <= y.Int64()
}

// IsLessThan returns true if x is less than y
func (x *BigInt) IsLessThanOrEqualsUint64(y uint64) bool {
	return x.Uint64() <= y
}

// Equals returns true if x equals to y
func (x *BigInt) Equals(y *BigInt) bool {
	return x.Int64() == y.Int64()
}

// Equals returns true if x equals to y
func (x *BigInt) EqualsInt(y int) bool {
	return x.Int32() == y
}

// Equals returns true if x equals to y
func (x *BigInt) EqualsInt64(y int64) bool {
	return x.Int64() == y
}

// Equals returns true if x equals to y
func (x *BigInt) EqualsUint64(y uint64) bool {
	return x.Uint64() == y
}

// Equals returns true if x equals to y
func (x *BigInt) Add(y int64) *BigInt {
	return NewBigInt64(x.Int64() + y)
}

// Equals returns true if x equals to y
func (x *BigInt) AddUint64(y uint64) *BigInt {
	return x.Add(int64(y))
}

func (x *BigInt) Int32() int {
	return int(x.Int64())
}

func (x *BigInt) Int64() int64 {
	if x.Pos {
		return int64(x.Value)
	}
	return int64(x.Value) * -1
}

func (x *BigInt) Uint64() uint64 {
	return uint64(x.Int64())
}

func (x *BigInt) Copy() *BigInt {
	cpy := *x
	return &cpy
}

func (x *BigInt) String() string {
	return fmt.Sprintf("%v", x.Int64())
}
