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

package kai

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/params"
	"github.com/kardiachain/go-kardia/blockchain"
	"github.com/kardiachain/go-kardia/blockchain/chaindb"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/log"
	"github.com/kardiachain/go-kardia/lib/rlp"
	"github.com/kardiachain/go-kardia/state"
	"github.com/kardiachain/go-kardia/tool"
	"github.com/kardiachain/go-kardia/types"
	"github.com/kardiachain/go-kardia/vm"
)

const (
	defaultGasPrice = 50 * params.Shannon
	defaultTimeOutForStaticCall = 5
)

// BlockJSON represents Block in JSON format
type BlockJSON struct {
	Hash           string               `json:"hash"`
	Height         uint64               `json:"height"`
	LastBlock      string               `json:"lastBlock"`
	CommitHash     string               `json:"commitHash"`
	Time           *big.Int             `json:"time"`
	NumTxs         uint64               `json:"num_txs"`
	GasLimit       uint64               `json:"gasLimit"`
	GasUsed        uint64               `json:"gasUsed"`
	Validator      string               `json:"validator"`
	TxHash         string               `json:"data_hash"`    // transactions
	Root           string               `json:"stateRoot"`    // state root
	ReceiptHash    string               `json:"receiptsRoot"` // receipt root
	Bloom          int64                `json:"logsBloom"`
	ValidatorsHash string               `json:"validators_hash"` // validators for the current block
	ConsensusHash  string               `json:"consensus_hash"`
	Txs            []*PublicTransaction `json:"txs"`
}

// PublicKaiAPI provides APIs to access Kai full node-related
// information.
type PublicKaiAPI struct {
	kaiService *Kardia
}

// NewPublicKaiAPI creates a new Kai protocol API for full nodes.
func NewPublicKaiAPI(kaiService *Kardia) *PublicKaiAPI {
	return &PublicKaiAPI{kaiService}
}

// NewBlockJSON creates a new Block JSON data from Block
func NewBlockJSON(block types.Block) *BlockJSON {
	txs := block.Transactions()
	transactions := make([]*PublicTransaction, 0, len(txs))
	for index, tx := range txs {
		json := NewPublicTransaction(tx, block.Hash(), block.Height(), uint64(index))
		transactions = append(transactions, json)
	}

	return &BlockJSON{
		Hash:           block.Hash().Hex(),
		Height:         block.Height(),
		LastBlock:      block.Header().LastBlockID.String(),
		Txs:            transactions,
		CommitHash:     block.LastCommitHash().Hex(),
		Time:           block.Header().Time,
		NumTxs:         block.Header().NumTxs,
		GasLimit:       block.Header().GasLimit,
		GasUsed:        block.Header().GasUsed,
		Validator:      block.Header().Coinbase.Hex(),
		TxHash:         block.Header().TxHash.Hex(),
		Root:           block.Header().Root.Hex(),
		ReceiptHash:    block.Header().ReceiptHash.Hex(),
		Bloom:          block.Header().Bloom.Big().Int64(),
		ValidatorsHash: block.Header().ValidatorsHash.Hex(),
		ConsensusHash:  block.Header().ConsensusHash.Hex(),
	}
}

// BlockNumber returns current block number
func (s *PublicKaiAPI) BlockNumber() uint64 {
	return s.kaiService.blockchain.CurrentBlock().Height()
}

// GetBlockByHash returns block by block hash
func (s *PublicKaiAPI) GetBlockByHash(blockHash string) *BlockJSON {
	if blockHash[0:2] == "0x" {
		blockHash = blockHash[2:]
	}
	block := s.kaiService.blockchain.GetBlockByHash(common.HexToHash(blockHash))
	return NewBlockJSON(*block)
}

// GetBlockByNumber returns block by block number
func (s *PublicKaiAPI) GetBlockByNumber(blockNumber uint64) *BlockJSON {
	block := s.kaiService.blockchain.GetBlockByHeight(blockNumber)
	if block == nil {
		return nil
	}
	return NewBlockJSON(*block)
}

// Validator returns node's validator, nil if current node is not a validator
func (s *PublicKaiAPI) Validator() map[string]interface{} {
	if val := s.kaiService.csManager.Validator(); val != nil {
		return map[string]interface{}{
			"address":     val.Address.Hex(),
			"votingPower": val.VotingPower,
		}
	}
	return nil
}

// Validators returns a list of validator
func (s *PublicKaiAPI) Validators() []map[string]interface{} {
	if vals := s.kaiService.csManager.Validators(); vals != nil && len(vals) > 0 {
		results := make([]map[string]interface{}, len(vals))
		for i, val := range vals {
			results[i] = map[string]interface{}{
				"address":     val.Address.Hex(),
				"votingPower": val.VotingPower,
			}
		}
		return results
	}
	return nil
}

type PublicTransaction struct {
	BlockHash        string        `json:"blockHash"`
	BlockNumber      common.Uint64 `json:"blockNumber"`
	From             string        `json:"from"`
	Gas              common.Uint64 `json:"gas"`
	GasPrice         common.Uint64 `json:"gasPrice"`
	Hash             string        `json:"hash"`
	Input            string        `json:"input"`
	Nonce            common.Uint64 `json:"nonce"`
	To               string        `json:"to"`
	TransactionIndex uint          `json:"transactionIndex"`
	Value            common.Uint64 `json:"value"`
}

type Log struct {
	Address     string   `json:"address"`
	Topics      []string `json:"topics"`
	Data        string   `json:"data"`
	BlockHeight uint64   `json:"blockHeight"`
	TxHash      string   `json:"transactionHash"`
	TxIndex     uint     `json:"transactionIndex"`
	BlockHash   string   `json:"blockHash"`
	Index       uint     `json:"logIndex"`
	Removed     bool     `json:"removed"`
}

// NewPublicTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func NewPublicTransaction(tx *types.Transaction, blockHash common.Hash, blockNumber uint64, index uint64) *PublicTransaction {
	from, _ := types.Sender(tx)

	result := &PublicTransaction{
		From:     from.Hex(),
		Gas:      common.Uint64(tx.Gas()),
		GasPrice: common.Uint64(tx.GasPrice().Int64()),
		Hash:     tx.Hash().Hex(),
		Input:    common.Encode(tx.Data()),
		Nonce:    common.Uint64(tx.Nonce()),
		Value:    common.Uint64(tx.Value().Int64()),
	}
	if tx.To() != nil {
		result.To = tx.To().Hex()
	} else {
		result.To = "0x"
	}
	if blockHash != (common.Hash{}) {
		result.BlockHash = blockHash.Hex()
		result.BlockNumber = common.Uint64(blockNumber)
		result.TransactionIndex = uint(index)
	}
	return result
}

// PublicTransactionAPI provides public apis relate to transactions
type PublicTransactionAPI struct {
	s *Kardia
}

// NewPublicTransactionAPI is a constructor of PublicTransactionAPI
func NewPublicTransactionAPI(service *Kardia) *PublicTransactionAPI {
	return &PublicTransactionAPI{service}
}

// SendRawTransaction decode encoded data into tx and then add tx into pool
func (a *PublicTransactionAPI) SendRawTransaction(ctx context.Context, txs string) (string, error) {
	log.Info("SendRawTransaction", "data", txs)
	tx := new(types.Transaction)
	encodedTx := common.FromHex(txs)
	if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
		return common.Hash{}.Hex(), err
	}
	return tx.Hash().Hex(), a.s.TxPool().AddRemote(tx)
}

// SendRawTransactions decode encoded data into tx and then add tx into pool
func (a *PublicTransactionAPI) SendRawTransactions(ctx context.Context, txs ...string) ([]string, error) {
	log.Info("SendRawTransaction", "data", txs)
	if len(txs) > 0 {
		results := make([]string, 0)
		for _, tx := range txs {
			if hash, err := a.SendRawTransaction(ctx, tx); err != nil {
				return results, err
			} else {
				results = append(results, hash)
			}

		}
	}
	return nil, nil
}

// KardiaCall execute a contract method call only against
// state on the local node. No tx is generated and submitted
// onto the blockchain
func (s *PublicKaiAPI) KardiaCall(ctx context.Context, call types.CallArgsJSON, blockNumber uint64) (string, error) {
	args := types.NewArgs(call)
	result, _, _, err := s.doCall(ctx, args, blockNumber, vm.Config{}, defaultTimeOutForStaticCall * time.Second)
	return common.Encode(result), err
}

// PendingTransactions returns pending transactions
func (a *PublicTransactionAPI) PendingTransactions() ([]*PublicTransaction, error) {
	pending, err := a.s.TxPool().Pending()
	if err != nil {
		return nil, err
	}

	transactions := make([]*PublicTransaction, 0, len(pending))

	// loop through pending txs
	for _, txs := range pending {
		for _, tx := range txs {
			jsonData := NewPublicTransaction(tx, common.Hash{}, 0, 0)
			transactions = append(transactions, jsonData)
		}
	}

	return transactions, nil
}

// GetTransaction gets transaction by transaction hash
func (a *PublicTransactionAPI) GetTransaction(hash string) *PublicTransaction {
	txHash := common.HexToHash(hash)
	tx, blockHash, height, index := chaindb.ReadTransaction(a.s.kaiDb, txHash)
	return NewPublicTransaction(tx, blockHash, height, index)
}

func (a *PublicTransactionAPI) getReceipts(hash common.Hash) (types.Receipts, error) {
	height := chaindb.ReadHeaderNumber(a.s.kaiDb, hash)
	if height == nil {
		return nil, nil
	}
	return chaindb.ReadReceipts(a.s.kaiDb, hash, *height), nil
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash.
func (a *PublicTransactionAPI) GetTransactionReceipt(ctx context.Context, hash string) (map[string]interface{}, error) {
	txHash := common.HexToHash(hash)
	tx, blockHash, height, index := chaindb.ReadTransaction(a.s.kaiDb, txHash)
	if tx == nil {
		return nil, nil
	}
	receipts, err := a.getReceipts(blockHash)
	if err != nil {
		return nil, err
	}
	if len(receipts) <= int(index) {
		return nil, nil
	}
	receipt := receipts[index]
	from, _ := types.Sender(tx)
	logs := make([]Log, 0)
	if receipt.Logs != nil {
		logs = make([]Log, 0, len(receipt.Logs))
		for _, l := range receipt.Logs {
			topics := make([]string, 0, len(l.Topics))
			for _, topic := range l.Topics {
				topics = append(topics, topic.Hex())
			}
			logs = append(logs, Log{
				Address:     l.Address.Hex(),
				Topics:      topics,
				Data:        hex.EncodeToString(l.Data),
				BlockHeight: l.BlockHeight,
				TxHash:      l.TxHash.Hex(),
				TxIndex:     l.TxIndex,
				BlockHash:   l.BlockHash.Hex(),
				Index:       l.Index,
				Removed:     l.Removed,
			})
		}
	}

	fields := map[string]interface{}{
		"blockHash":         blockHash.Hex(),
		"blockHeight":       uint64(height),
		"transactionHash":   hash,
		"transactionIndex":  uint64(index),
		"from":              from.Hex(),
		"to":                nil,
		"gasUsed":           uint64(receipt.GasUsed),
		"cumulativeGasUsed": uint64(receipt.CumulativeGasUsed),
		"contractAddress":   nil,
		"logs":              logs,
		"logsBloom":         receipt.Bloom.Big().Int64(),
	}

	// To field is nil for contract creation tx.
	if tx.To() != nil {
		fields["to"] = tx.To().Hex()
	}
	// Assign receipt status or post state.
	if len(receipt.PostState) > 0 {
		fields["root"] = common.Bytes(receipt.PostState)
	} else {
		fields["status"] = uint(receipt.Status)
	}
	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if receipt.ContractAddress != (common.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress.Hex()
	}
	return fields, nil
}

// PublicAccountAPI provides APIs support getting account's info
type PublicAccountAPI struct {
	kaiService *Kardia
}

// NewPublicAccountAPI is a constructor that init new PublicAccountAPI
func NewPublicAccountAPI(kaiService *Kardia) *PublicAccountAPI {
	return &PublicAccountAPI{kaiService}
}

// Balance returns address's balance
func (a *PublicAccountAPI) Balance(address string, hash string, height int64) int64 {
	addr := common.HexToAddress(address)
	log.Info("Addr", "addr", addr.Hex())
	block := new(types.Block)
	if len(hash) > 0 && height >= 0 {
		block = a.kaiService.blockchain.GetBlock(common.HexToHash(hash), uint64(height))
	} else if len(hash) > 0 {
		block = a.kaiService.blockchain.GetBlockByHash(common.HexToHash(hash))
	} else if height >= 0 {
		block = a.kaiService.blockchain.GetBlockByHeight(uint64(height))
	} else {
		block = a.kaiService.blockchain.CurrentBlock()
	}
	state, err := a.kaiService.blockchain.StateAt(block.Root())
	if err != nil {
		log.Error("Fail to get state from block", "err", err, "block", block)
		return -1
	}
	return state.GetBalance(addr).Int64()
}

// Nonce return address's nonce
func (a *PublicAccountAPI) Nonce(address string) (uint64, error) {
	addr := common.HexToAddress(address)
	block := a.kaiService.blockchain.CurrentBlock()
	state, err := a.kaiService.blockchain.StateAt(block.Root())
	if err != nil {
		log.Error("Fail to get state from block", "err", err, "block", block)
		return 0, err
	}
	return state.GetNonce(addr), nil
}

// doCall is an interface to make smart contract call against the state of local node
// No tx is generated or submitted to the blockchain
func (s *PublicKaiAPI) doCall(ctx context.Context, args *types.CallArgs, blockNr uint64, vmCfg vm.Config, timeout time.Duration) ([]byte, uint64, bool, error) {
	defer func(start time.Time) { log.Debug("Executing KVM call finished", "runtime", time.Since(start)) }(time.Now())

	var (
		statedb *state.StateDB
		err     error
		header  *types.Header
	)
	// If blockNr is specified, we used the state and header at the block at height blockNr
	// otherwise we use the current state and header
	if blockNr > 0 {
		block := s.kaiService.BlockChain().GetBlockByHeight(blockNr)
		statedb, err = s.kaiService.BlockChain().StateAt(block.Root())
		header = block.Header()
	} else {
		statedb, err = s.kaiService.BlockChain().State()
		header = s.kaiService.BlockChain().CurrentHeader()
	}

	if statedb == nil || err != nil {
		return nil, 0, false, err
	}
	// Set sender address or use a default if none specified
	addr := args.From

	if addr == (common.Address{}) {
		addr = tool.GetRandomGenesisAccount()
	}

	// Set default gas & gas price if none were set
	gas, gasPrice := uint64(args.Gas), args.GasPrice

	if gas == 0 {
		gas = common.MaxInt64 / 2
	}
	if gasPrice.Sign() == 0 {
		gasPrice = new(big.Int).SetUint64(defaultGasPrice)
	}

	// Create new call message
	msg := types.NewMessage(addr, args.To, 0, args.Value, gas, gasPrice, args.Data, false)
	// Setup context so it may be cancelled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	// Create a new context to be used in the KVM environment
	context := blockchain.NewKVMContext(msg, header, s.kaiService.BlockChain())
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	kvm := vm.NewKVM(context, statedb, vmCfg)
	// Wait for the context to be done and cancel the KVM. Even if the
	// KVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		kvm.Cancel()
	}()
	// Apply the transaction to the current state (included in the env)
	gp := new(blockchain.GasPool).AddGas(common.MaxUint64)
	res, gas, failed, err := blockchain.ApplyMessage(kvm, msg, gp)
	if err != nil {
		return nil, 0, false, err
	}
	return res, gas, failed, err
}

// EstimateGas returns an estimate of the amount of gas needed to execute the
// given transaction against the current pending block.
func (s *PublicKaiAPI) EstimateGas(ctx context.Context, call types.CallArgsJSON) (uint64, error) {
	// Binary search the gas requirement, as it may be higher than the amount used
	var (
		lo  uint64 = params.TxGas - 1
		hi  uint64
		cap uint64
	)
	args := types.NewArgs(call)
	if uint64(args.Gas) >= params.TxGas {
		hi = uint64(args.Gas)
	} else {
		// Retrieve the current pending block to act as the gas ceiling
		block := s.kaiService.BlockChain().CurrentBlock()
		hi = block.GasLimit()
	}
	cap = hi

	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(gas uint64) bool {
		args.Gas = gas

		_, _, failed, err := s.doCall(ctx, args, s.BlockNumber(), vm.Config{}, 0)
		if err != nil || failed {
			return false
		}
		return true
	}
	// Execute the binary search and hone in on an executable gas limit
	for lo+1 < hi {
		mid := (hi + lo) / 2
		if !executable(mid) {
			lo = mid
		} else {
			hi = mid
		}
	}
	// Reject the transaction as invalid if it still fails at the highest allowance
	if hi == cap {
		if !executable(hi) {
			return 0, fmt.Errorf("gas required exceeds allowance or always failing transaction")
		}
	}
	return hi, nil
}