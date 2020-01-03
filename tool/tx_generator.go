/*
 *  Copyright 2018 KardiaChain
 *  This file is part of the go-kardia library.
 *
 *  The go-kardia library is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Lesser General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  The go-kardia library is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 *  GNU Lesser General Public License for more details.
 *
 *  You should have received a copy of the GNU Lesser General Public License
 *  along with the go-kardia library. If not, see <http://www.gnu.org/licenses/>.
 */

package tool

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/kardiachain/go-kardia/kai/base"
	"github.com/kardiachain/go-kardia/kvm"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/kardiachain/go-kardia/kai/state"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/crypto"
	"github.com/kardiachain/go-kardia/lib/log"
	"github.com/kardiachain/go-kardia/mainchain/tx_pool"
	"github.com/kardiachain/go-kardia/types"
)

type Account struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

const (
	DefaultGasLimit             = kvm.TxGas // currently we don't care about tx fee and cost.
	DefaultFaucetAcc            = "0x2BB7316884C7568F2C6A6aDf2908667C0d241A66"
	DefaultFaucetPrivAcc        = "4561f7d91a4f95ef0a72550fa423febaad3594f91611f9a2b10a7af4d3deb9ed"
	DefaultGenRandomWithStateTx = 1
	DefaultGenRandomTx          = 2
	// constants related to account to call smc
	KardiaAccountToCallSmc = "0xBA30505351c17F4c818d94a990eDeD95e166474b"
	KardiaPrivKeyToCallSmc = "ae1a52546294bed6e734185775dbc84009de00bdf51b709471e2415c31ceeed7"
)


// ======================= Genesis Const =======================

var InitValue = big.NewInt(int64(math.Pow10(10))) // Update Genesis Account Values
var InitValueInCell = InitValue.Mul(InitValue, big.NewInt(int64(math.Pow10(18))))

// GenesisAccounts are used to initialized accounts in genesis block
var GenesisAccounts = map[string]*big.Int{
	// TODO(kiendn): These addresses are same of node address. Change to another set.
	"0xc1fe56E3F58D3244F606306611a5d10c8333f1f6": InitValueInCell,
	"0x7cefC13B6E2aedEeDFB7Cb6c32457240746BAEe5": InitValueInCell,
	"0xfF3dac4f04dDbD24dE5D6039F90596F0a8bb08fd": InitValueInCell,
	"0x071E8F5ddddd9f2D4B4Bdf8Fc970DFe8d9871c28": InitValueInCell,
	"0x94FD535AAB6C01302147Be7819D07817647f7B63": InitValueInCell,
	"0xa8073C95521a6Db54f4b5ca31a04773B093e9274": InitValueInCell,
	"0xe94517a4f6f45e80CbAaFfBb0b845F4c0FDD7547": InitValueInCell,
	"0xBA30505351c17F4c818d94a990eDeD95e166474b": InitValueInCell,
	"0x212a83C0D7Db5C526303f873D9CeaA32382b55D0": InitValueInCell,
	"0x8dB7cF1823fcfa6e9E2063F983b3B96A48EEd5a4": InitValueInCell,
	"0x66BAB3F68Ff0822B7bA568a58A5CB619C4825Ce5": InitValueInCell,
	"0x88e1B4289b639C3b7b97899Be32627DCd3e81b7e": InitValueInCell,
	"0xCE61E95666737E46B2453717Fe1ba0d9A85B9d3E": InitValueInCell,
	"0x1A5193E85ffa06fde42b2A2A6da7535BA510aE8C": InitValueInCell,
	"0xb19BC4477ff32EC13872a2A827782DeA8b6E92C0": InitValueInCell,
	"0x0fFFA18f6c90ce3f02691dc5eC954495EA483046": InitValueInCell,
	"0x8C10639F908FED884a04C5A49A2735AB726DDaB4": InitValueInCell,
	"0x2BB7316884C7568F2C6A6aDf2908667C0d241A66": InitValueInCell,

	// TODO(namdoh): Re-enable after parsing node index fixed in main.go
	//"0x36BE7365e6037bD0FDa455DC4d197B07A2002547": 100000000,
}

//  GenesisAddrKeys maps genesis account addresses to private keys.
var GenesisAddrKeys = map[string]string{
	"0xc1fe56E3F58D3244F606306611a5d10c8333f1f6": "8843ebcb1021b00ae9a644db6617f9c6d870e5fd53624cefe374c1d2d710fd06",
	"0x7cefC13B6E2aedEeDFB7Cb6c32457240746BAEe5": "77cfc693f7861a6e1ea817c593c04fbc9b63d4d3146c5753c008cfc67cffca79",
	"0xfF3dac4f04dDbD24dE5D6039F90596F0a8bb08fd": "98de1df1e242afb02bd5dc01fbcacddcc9a4d41df95a66f629139560ca6e4dbb",
	"0x071E8F5ddddd9f2D4B4Bdf8Fc970DFe8d9871c28": "32f5c0aef7f9172044a472478421c63fd8492640ff2d0eaab9562389db3a8efe",
	"0x94FD535AAB6C01302147Be7819D07817647f7B63": "68b53a92d846baafdc782cb9cad65d77020c8d747eca7b621370b52b18c91f9a",
	"0xa8073C95521a6Db54f4b5ca31a04773B093e9274": "049de018e08c3bcd59c1a21f0cf7de8f17fe51f8ce7d9c2120d17b1f0251b265",
	"0xe94517a4f6f45e80CbAaFfBb0b845F4c0FDD7547": "9fdd56a3c2a536dc8f981d935f0f3f2ea04e125547fdfffa37e157ce86ff1007",
	"0xBA30505351c17F4c818d94a990eDeD95e166474b": "ae1a52546294bed6e734185775dbc84009de00bdf51b709471e2415c31ceeed7",
	"0x212a83C0D7Db5C526303f873D9CeaA32382b55D0": "b34bd81838a4a335fb3403d0bf616eca1eb9a4b4716c7dda7c617503cfeaab67",
	"0x8dB7cF1823fcfa6e9E2063F983b3B96A48EEd5a4": "0cf7ae0332a891044659ace49a0732fa07c2872b4aef479945501f385a23e689",
	"0x66BAB3F68Ff0822B7bA568a58A5CB619C4825Ce5": "2003be66077b0873f5bedb32a596530eb8a0c908c32dda7771f169ee137c1f82",
	"0x88e1B4289b639C3b7b97899Be32627DCd3e81b7e": "9dce5ec0b40e363e898f296c01345c12a0961f1cccad098964c73ed85fef5850",
	"0xCE61E95666737E46B2453717Fe1ba0d9A85B9d3E": "f0b2f6f24b70481a51712639badf0e5587545080dc53e0664770adb9881823fb",
	"0x1A5193E85ffa06fde42b2A2A6da7535BA510aE8C": "83731e17afb0da61c0026eaf780364eee367c50a5225ece89a63ad94a4a1f088",
	"0xb19BC4477ff32EC13872a2A827782DeA8b6E92C0": "fc09d3f004b1ee430fee60568aa29748e277e76f1f372eea9d2b9ff1e27bdfdb",
	"0x0fFFA18f6c90ce3f02691dc5eC954495EA483046": "5605dd5f4db003c396956b4b80c093c472ccef4021181aa910125d7c57324152",
	"0x8C10639F908FED884a04C5A49A2735AB726DDaB4": "9813a1dffe303131d1fe80b6fe872206267abd8ff84a52c907b0d32df582b1eb",
	"0x2BB7316884C7568F2C6A6aDf2908667C0d241A66": "4561f7d91a4f95ef0a72550fa423febaad3594f91611f9a2b10a7af4d3deb9ed",
	// TODO(namdoh): Re-enable after parsing node index fixed in main.go
	//"e049a09c992c882bc2deb780323a247c6ee0951f8b4c5c1dd0fc2fc22ce6493d": "0x36BE7365e6037bD0FDa455DC4d197B07A2002547",
}

var (
	defaultAmount   = big.NewInt(10)
	defaultGasPrice = big.NewInt(1)
)

type GeneratorTool struct {
	nonceMap map[string]uint64 // Map of nonce counter for each address
	accounts []Account
	mu       sync.Mutex
}

func NewGeneratorTool(accounts []Account) *GeneratorTool {
	genTool := new(GeneratorTool)
	genTool.nonceMap = make(map[string]uint64, 0)
	genTool.accounts = accounts
	return genTool
}

// GenerateTx generate an array of transfer transactions within genesis accounts.
// numTx: number of transactions to send, default to 10.
func (genTool *GeneratorTool) GenerateTx(numTx int) []*types.Transaction {
	if numTx <= 0 || len(genTool.accounts) == 0 {
		return nil
	}
	result := make([]*types.Transaction, numTx)
	genTool.mu.Lock()

	signer := types.HomesteadSigner{}

	for i := 0; i < numTx; i++ {
		senderKey, toAddr := randomTxAddresses(genTool.accounts)
		senderPublicKey := crypto.PubkeyToAddress(senderKey.PublicKey)
		senderAddrS := senderPublicKey.String()
		nonce := genTool.nonceMap[senderAddrS]
		amount := big.NewInt(int64(RandomInt(10, 20)))
		amount = amount.Mul(amount, big.NewInt(int64(math.Pow10(18))))
		tx, err := types.SignTx(signer, types.NewTransaction(
			nonce,
			toAddr,
			amount,
			1000,
			big.NewInt(1),
			nil,
		), senderKey)
		if err != nil {
			panic(fmt.Sprintf("Fail to sign generated tx: %v", err))
		}
		result[i] = tx
		nonce += 1
		genTool.nonceMap[senderAddrS] = nonce
	}
	genTool.mu.Unlock()
	return result
}

func (genTool *GeneratorTool) GenerateRandomTx(numTx int) types.Transactions {
	if numTx <= 0 || len(genTool.accounts) == 0 {
		return nil
	}
	result := make(types.Transactions, numTx)
	genTool.mu.Lock()

	signer := types.HomesteadSigner{}
	for i := 0; i < numTx; i++ {
		senderKey, toAddr := randomTxAddresses(genTool.accounts)
		senderPublicKey := crypto.PubkeyToAddress(senderKey.PublicKey)
		amount := big.NewInt(int64(RandomInt(10, 20)))
		amount = amount.Mul(amount, big.NewInt(int64(math.Pow10(18))))
		senderAddrS := senderPublicKey.String()

		if _, ok := genTool.nonceMap[senderAddrS]; !ok {
			genTool.nonceMap[senderAddrS] = 1
		}
		nonce := genTool.nonceMap[senderAddrS]
		tx, err := types.SignTx(signer, types.NewTransaction(
			nonce,
			toAddr,
			amount,
			DefaultGasLimit,
			defaultGasPrice,
			nil,
		), senderKey)
		if err != nil {
			panic(fmt.Sprintf("Fail to sign generated tx: %v", err))
		}
		result[i] = tx
		nonce += 1
		genTool.nonceMap[senderAddrS] = nonce
	}
	genTool.mu.Unlock()
	return result
}

func (genTool *GeneratorTool) GenerateRandomTxWithState(numTx int, stateDb *state.StateDB) []*types.Transaction {
	if numTx <= 0 || len(genTool.accounts) == 0 {
		return nil
	}
	result := make([]*types.Transaction, numTx)
	genTool.mu.Lock()

	signer := types.HomesteadSigner{}

	for i := 0; i < numTx; i++ {
		senderKey, toAddr := randomTxAddresses(genTool.accounts)
		senderPublicKey := crypto.PubkeyToAddress(senderKey.PublicKey)
		nonce := stateDb.GetNonce(senderPublicKey)
		amount := big.NewInt(int64(RandomInt(10, 20)))
		amount = amount.Mul(amount, big.NewInt(int64(math.Pow10(18))))
		senderAddrS := senderPublicKey.String()

		//get nonce from sender mapping
		nonceMap := genTool.GetNonce(senderAddrS)
		if nonce < nonceMap { // check nonce from statedb and nonceMap
			nonce = nonceMap
		}

		tx, err := types.SignTx(signer, types.NewTransaction(
			nonce,
			toAddr,
			amount,
			DefaultGasLimit,
			defaultGasPrice,
			nil,
		), senderKey)
		if err != nil {
			panic(fmt.Sprintf("Fail to sign generated tx: %v", err))
		}
		result[i] = tx
		nonce += 1
		genTool.nonceMap[senderAddrS] = nonce
	}
	genTool.mu.Unlock()
	return result
}

func (genTool *GeneratorTool) GenerateRandomTxWithAddressState(numTx int, txPool base.TxPool) types.Transactions {
	if numTx <= 0 || len(genTool.accounts) == 0 {
		return nil
	}
	result := make(types.Transactions, numTx)
	genTool.mu.Lock()

	signer := types.HomesteadSigner{}

	for i := 0; i < numTx; i++ {
		senderKey, toAddr := randomTxAddresses(genTool.accounts)
		senderPublicKey := crypto.PubkeyToAddress(senderKey.PublicKey)
		nonce := txPool.Nonce(senderPublicKey)
		amount := big.NewInt(int64(RandomInt(10, 20)))
		amount = amount.Mul(amount, big.NewInt(int64(math.Pow10(18))))
		senderAddrS := senderPublicKey.String()

		//get nonce from sender mapping
		nonceMap := genTool.GetNonce(senderAddrS)
		if nonce < nonceMap { // check nonce from statedb and nonceMap
			nonce = nonceMap
		}

		tx, err := types.SignTx(signer, types.NewTransaction(
			nonce,
			toAddr,
			amount,
			DefaultGasLimit,
			defaultGasPrice,
			nil,
		), senderKey)
		if err != nil {
			panic(fmt.Sprintf("Fail to sign generated tx: %v", err))
		}
		result[i] = tx
		nonce += 1
		genTool.nonceMap[senderAddrS] = nonce
	}
	genTool.mu.Unlock()
	return result
}

func (genTool *GeneratorTool) GetNonce(address string) uint64 {
	return genTool.nonceMap[address]
}

// GenerateSmcCall generates tx which call a smart contract's method
// if isIncrement is true, nonce + 1 to prevent duplicate nonce if generateSmcCall is called twice.
func GenerateSmcCall(senderKey *ecdsa.PrivateKey, address common.Address, input []byte, txPool *tx_pool.TxPool, isIncrement bool) *types.Transaction {
	senderAddress := crypto.PubkeyToAddress(senderKey.PublicKey)
	nonce := txPool.Nonce(senderAddress)
	if isIncrement {
		nonce++
	}

	signer := types.HomesteadSigner{}

	tx, err := types.SignTx(signer, types.NewTransaction(
		nonce,
		address,
		big.NewInt(0),
		5000000,
		big.NewInt(1),
		input,
	), senderKey)
	if err != nil {
		panic(fmt.Sprintf("Fail to generate smc call: %v", err))
	}
	log.Error("GenerateSmcCall", "nonce", tx.Nonce(), "tx", tx.Hash().Hex())
	return tx
}

func randomTxAddresses(accounts []Account) (senderKey *ecdsa.PrivateKey, toAddr common.Address) {
	for {
		senderKey = randomGenesisPrivateKey(accounts)
		toAddr = randomGenesisAddress()
		privateKeyBytes := crypto.FromECDSA(senderKey)
		privateKeyHex := hexutil.Encode(privateKeyBytes)[2:]
		if senderKey != nil && crypto.PubkeyToAddress(senderKey.PublicKey) != toAddr && privateKeyHex != KardiaPrivKeyToCallSmc && privateKeyHex != DefaultFaucetPrivAcc {
			// skip senderAddr = toAddr && senderAddr that call smc
			break
		}
	}
	return senderKey, toAddr
}

func randomGenesisAddress() common.Address {
	size := len(GenesisAddrKeys)
	randomI := rand.Intn(size)
	index := 0
	for addrS := range GenesisAddrKeys {
		if index == randomI {
			return common.HexToAddress(addrS)
		}
		index++
	}
	panic("impossible failure")
}

func randomAddress() common.Address {
	rand.Seed(time.Now().UTC().UnixNano())
	address := make([]byte, 20)
	rand.Read(address)
	return common.BytesToAddress(address)
}

func randomGenesisPrivateKey(accounts []Account) *ecdsa.PrivateKey {
	size := len(accounts)
	randomI := rand.Intn(size)
	index := 0
	for _, account := range accounts {
		if index == randomI {
			pkByte, _ := hex.DecodeString(account.PrivateKey)
			return crypto.ToECDSAUnsafe(pkByte)
		}
		index++
	}
	panic("impossible failure")
}

func GetRandomGenesisAccount() common.Address {
	size := len(GenesisAccounts)
	randomI := rand.Intn(size)
	index := 0
	for addrS := range GenesisAccounts {
		if index == randomI {
			return common.HexToAddress(addrS)
		}
		index++
	}
	panic("impossible failure")
}

func GetAccounts(genesisAccounts map[string]string) []Account {
	accounts := make([]Account, 0)
	for key, value := range genesisAccounts {
		accounts = append(accounts, Account{
			Address:    key,
			PrivateKey: value,
		})
	}

	return accounts
}

func RandomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	n := min + rand.Intn(max-min+1)
	return n
}
