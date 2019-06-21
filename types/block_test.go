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

package types

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/crypto"
	"github.com/kardiachain/go-kardia/lib/log"
	"github.com/kardiachain/go-kardia/lib/rlp"
)

var testAccounts = map[string]string{
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
	"0x757906bA5023B92e980F61cA9427BFC15f810B6f": "330fcc18cc5c9d438744ad9b0f4274d5e9d34f099b8acf0066e304f4acabdc90",
	"0xf80927B05dc9a25247F3d68B8192cD78361C79f0": "b0de1f5d622fcea073a06386ddffce526165b38a889c04e5088ef88a3d807570",
	"0xA58Cc74A805adb80DA35925F28f7B732Bbf3C9FF": "209f9f8182be0f07f785be76a35c82c7ad2be4359f5c95c610e005dd23e1cfaa",
	"0x9afD6D161E3e19c0191d83D3f48384DdBF24Ad7b": "3b33cfc81ebdbdf699d6f9891f8e4f1f66a8e61372416a4ed068985f7bf2cb8d",
	"0x8A618C045Cb00E8abf6A9cab60edd10dd5572119": "99b860be69e7015b0b6d14b45b3d8534ae8a7b58c673c81369d6d61de511387c",
	"0xA6E9bA68b1C3d9b45c1a07Aeb04169465CCAa69f": "f19d51628268e254db4d7f3a1108a72c36073fe4b4298882cb729b690291ad69",
	"0xAdea9867E143445A92A5DDD11B447de59f0090A6": "761225588f7180bf37fa6d1e2829db3cc4a8972ba6358f151c5249bffa402d1e",
	"0x5360D2e4596CB10211bef5e088A3ED705a6872B1": "5c0eb44a8d1219bd995cfd07245ba3bd6b63a9cc76549e91b0641cfb888d092e",
	"0x74474e320Bf05D9D83B61946faC7C0c49C86e634": "54fd01571a9416fd213c62b8fe62b870e8f1cab5f609ce39e51daa6fdfac2535",
	"0x3c60F22441283FDB7acdE17f1758673e7024a68A": "dee4955aa6bc6d74f910ffd586f2ab16ed8c63b580f49fc6ca219015e0a80ec4",
	"0x8c365De6Cf9644276Ce3410D0BBBB8dE06d30633": "229756158eb85c278fd1c5c909f0fbb0466cec43b65713a036ee2a8feee51360",
	"0x04b7361c47432eA35B4696FA9D93F99fC7C7FeCa": "0155b454947fc2647bf1d500a63e7d57119de8b428cc8f6a0ad329963f8dd36f",
	"0x26fBc0553bF92A1f53b216bBdc9DF93F19209773": "05c82e623edde71f3613a0d6f2885caa00a9b03c10f62a82d2c8753ff8abd380",
	"0xB8258F7EA7B4C0a2Aa13a2F4C092e4a3eCf2a379": "bf0abaef91881b51c4d9141893e9ff8f8782b6155415c6c1c0ba1c83b5e88207",
	"0x1c9bD2E569d990c8183D685C58CFD1af948D2A8a": "3c70d08ea634e7f1e511534ccc94559bf069abe1558c6cff4aefd3413cf725b4",
	"0x90056191A8aEE8c529756bB97DAa0e7524c4b1aA": "8f075735a79b61d1cdcbd24dccf9701e433814825fe5853f7b6ed21884ad306c",
	"0x58Aa25D60FAa22BD29EB986a79931d33A2B9C462": "20d11078782250c7a65ba7d3d304ce703a17322b4279a89463f1c96589b130df",
	"0x6F55C53102e4493CAd5620Ee4ad38a28ed65997C": "0d282ec335ce58f73d69b633a1c0b31fae25edc66b7c755110b9a94575e7f811",
	"0xA3a7523F5183788e1048E24FAb285439a92c3647": "8fe30f1c900334741692ee8f41669d6a8acfe5a29e0f788652f236552bd22c3a",
	"0x96d359D752611255eAd3465Ae4E310b47a5a20b0": "14d86d15da9376f8b0b75b340e213ffd984081d77ecb6d3365dce3a25a3dbb7c",
	"0x737E2d1562b16FA1Bf8C7E510F0be32c0BDA059c": "375eb77a7dcd5ccd6becf6c1ae4a07ae4a2a36d55b52231afd8cd8839085b452",
	"0x291bCb8f5199Bdcd8c611b91209f1a2359AC2FdD": "93c5619c33e0184cb6782551421bb210bc102e185578c687905db0fbd23ee472",
	"0x595e1c030E76eF64633046A7Ca4deCcfce952C73": "7e91210a633c6c59b3e3ee09cb8e13345e14d8723bd6959e7907da2efecb3957",
}

func TestBlockCreation(t *testing.T) {
	block := CreateNewBlock(1)
	if err := block.ValidateBasic(); err != nil {
		t.Fatal("Init block error", err)
	}
}

func TestBlockEncodeDecode(t *testing.T) {
	block := CreateNewBlock(1)
	encodedBlock, err := rlp.EncodeToBytes(&block)
	if err != nil {
		t.Fatal("encode error: ", err)
	}
	var decodedBlock Block
	if err := rlp.DecodeBytes(encodedBlock, &decodedBlock); err != nil {
		t.Fatal("decode error: ", err)
	}

	if decodedBlock.Hash() != block.Hash() {
		t.Error("Encode Decode block error")
	}
}

func TestNewDualBlock(t *testing.T) {
	block := CreateNewDualBlock()
	if err := block.ValidateBasic(); err != nil {
		t.Fatal("Error validating New Dual block", err)
	}
}

func TestBlockEncodeDecodeFile(t *testing.T) {
	block := CreateNewBlock(1)
	blockCopy := block.WithBody(block.Body())
	encodeFile, err := os.Create("encodeFile.txt")
	defer encodeFile.Close()
	if err != nil {
		t.Error("Error creating file")
	}

	if err := block.EncodeRLP(encodeFile); err != nil {
		t.Fatal("Error encoding block")
	}

	f, err := os.Open("encodeFile.txt")
	if err != nil {
		t.Error("Error opening file:", err)
	}

	stream := rlp.NewStream(f, 99999)
	if err := block.DecodeRLP(stream); err != nil {
		t.Fatal("Decoding block error:", err)
	}
	if block.Hash() != blockCopy.Hash() {
		t.Fatal("Encode Decode File error")
	}

}

func TestGetDualEvents(t *testing.T) {
	dualBlock := CreateNewDualBlock()
	dualEvents := dualBlock.DualEvents()
	dualEventCopy := NewDualEvent(100, false, "KAI", new(common.Hash), new(EventSummary), new(DualActions))
	if dualEvents[0].Hash() != dualEventCopy.Hash() {
		t.Error("Dual Events hash not equal")
	}
}

func TestBodyCreationAndCopy(t *testing.T) {
	body := CreateNewBlock(1).Body()
	copyBody := body.Copy()
	if rlpHash(body) != rlpHash(copyBody) {
		t.Fatal("Error copy body")
	}
}

func TestBodyEncodeDecodeFile(t *testing.T) {
	body := CreateNewBlock(1).Body()
	bodyCopy := body.Copy()
	encodeFile, err := os.Create("encodeFile.txt")
	if err != nil {
		t.Error("Error creating file")
	}

	if err := body.EncodeRLP(encodeFile); err != nil {
		t.Fatal("Error encoding block")
	}

	encodeFile.Close()

	f, err := os.Open("encodeFile.txt")
	if err != nil {
		t.Error("Error opening file:", err)
	}

	stream := rlp.NewStream(f, 99999)
	if err := body.DecodeRLP(stream); err != nil {
		t.Fatal("Decoding block error:", err)
	}
	defer f.Close()

	if rlpHash(body) != rlpHash(bodyCopy) {
		t.Fatal("Encode Decode from file error")
	}
}

func TestBlockWithBodyFunction(t *testing.T) {
	block := CreateNewBlock(1)
	body := CreateNewDualBlock().Body()

	blockWithBody := block.WithBody(body)
	bwbBody := blockWithBody.Body()
	if blockWithBody.header.Hash() != block.header.Hash() {
		t.Error("BWB Header Error")
	}
	for i := range bwbBody.Transactions {
		if bwbBody.Transactions[i] != body.Transactions[i] {
			t.Error("BWB Transaction Error")
			break
		}
	}
	for i := range bwbBody.DualEvents {
		if bwbBody.DualEvents[i] != body.DualEvents[i] {
			t.Error("BWB Dual Events Error")
			break
		}
	}
	if bwbBody.LastCommit != body.LastCommit {
		t.Error("BWB Last Commit Error")
	}
}

func TestNewZeroBlockID(t *testing.T) {
	blockID := NewZeroBlockID()
	if !blockID.IsZero() {
		t.Fatal("NewZeroBlockID is not empty")
	}
}

func TestBlockSorterSwap(t *testing.T) {
	firstBlock := CreateNewBlock(1)
	secondBlock := CreateNewBlock(3)
	blockSorter := blockSorter{
		blocks: []*Block{firstBlock, secondBlock},
	}
	blockSorter.Swap(0, 1)
	if blockSorter.blocks[0] != secondBlock && blockSorter.blocks[1] != firstBlock {
		t.Fatal("blockSorter Swap error")
	}
}

func TestBlockHeightFunction(t *testing.T) {
	lowerBlock := CreateNewBlock(1)
	higherBlock := CreateNewBlock(2)
	if Height(higherBlock, lowerBlock) {
		t.Fatal("block Height func error")
	} else if !Height(lowerBlock, higherBlock) {
		t.Fatal("Block Height func error")
	}
}

func TestBlockSortByHeight(t *testing.T) {
	GetBlockByHeight := BlockBy(
		func(b1, b2 *Block) bool {
			return b1.header.Height < b2.header.Height
		})
	b0 := CreateNewBlock(0)
	b1 := CreateNewBlock(1)
	b2 := CreateNewBlock(2)
	b3 := CreateNewBlock(3)
	blocks := []*Block{b3, b2, b1, b0}

	GetBlockByHeight.Sort(blocks)
	if !CheckSortedHeight(blocks) {
		t.Error("Blocks not sorted")
	}
}

func CheckSortedHeight(blocks []*Block) bool {
	prev := blocks[0].header.Height
	for i := range blocks {
		if prev > blocks[i].header.Height {
			return false
		}
		prev = blocks[i].header.Height
	}
	return true
}

func CreateNewBlock(height uint64) *Block {
	header := Header{
		Height: height,
		Time:   big.NewInt(time.Now().Unix()),
	}

	addr := common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
	key, _ := crypto.GenerateKey()
	emptyTx := NewTransaction(
		1,
		addr,
		big.NewInt(99), 1000, big.NewInt(100),
		nil,
	)
	signedTx, _ := SignTx(emptyTx, key)

	txns := []*Transaction{signedTx}

	vote := &Vote{
		ValidatorIndex: common.NewBigInt64(1),
		Height:         common.NewBigInt64(2),
		Round:          common.NewBigInt64(1),
		Timestamp:      big.NewInt(100),
		Type:           VoteTypePrecommit,
	}
	lastCommit := &Commit{
		Precommits: []*Vote{vote, nil},
	}
	return NewBlock(log.New(), &header, txns, nil, lastCommit)
}

func CreateNewDualBlock() *Block {
	header := Header{
		Height: 1,
		Time:   big.NewInt(1),
	}
	vote := &Vote{
		ValidatorIndex: common.NewBigInt64(1),
		Height:         common.NewBigInt64(2),
		Round:          common.NewBigInt64(1),
		Timestamp:      big.NewInt(100),
		Type:           VoteTypePrecommit,
	}
	lastCommit := &Commit{
		Precommits: []*Vote{vote, vote},
	}
	header.LastCommitHash = lastCommit.Hash()
	de := NewDualEvent(100, false, "KAI", new(common.Hash), new(EventSummary), new(DualActions))
	return NewDualBlock(log.New(), &header, []*DualEvent{de, nil}, lastCommit)
}

func CreateNewBlockWithTxs(height uint64) *Block {

	lenTxs := 10000
	txns := make([]*Transaction, lenTxs)
	nonceMap := make(map[string]uint64)
	var keys []*ecdsa.PrivateKey
	var addresses []common.Address

	// Account type
	type Account struct {
		Address    string      `json:"address"`
		PrivateKey string      `json:"privateKey"`
	}

	header := Header{
		Height: height,
		Time:   big.NewInt(time.Now().Unix()),
	}

	// Generate accounts from testAccounts
	accounts := make([]Account, 0)
	for key, value := range testAccounts {
		accounts = append(accounts, Account{
			Address: key,
			PrivateKey: value,
		})
	}

	for _, account := range accounts {
		pkByte, _ := hex.DecodeString(account.PrivateKey)
		keys = append(keys, crypto.ToECDSAUnsafe(pkByte))
		addresses = append(addresses, common.HexToAddress(account.Address))
	}

	addrKeySize := len(addresses)

	for i := 0; i < lenTxs; i++ {
		senderKey := keys[i%addrKeySize]
		toAddr := addresses[(i+1)%addrKeySize]

		senderAddrS := crypto.PubkeyToAddress(senderKey.PublicKey).String()
		nonce := nonceMap[senderAddrS]

		tx, err := SignTx(NewTransaction(
			nonce,
			toAddr,
			big.NewInt(10),
			1000,
			big.NewInt(1),
			nil,
		), senderKey)
		if err != nil {
			panic(fmt.Sprintf("Fail to sign generated tx: %v", err))
		}
		txns[i] = tx
		nonce += 1
		nonceMap[senderAddrS] = nonce
	}

	vote := &Vote{
		ValidatorIndex: common.NewBigInt64(1),
		Height:         common.NewBigInt64(2),
		Round:          common.NewBigInt64(1),
		Timestamp:      big.NewInt(100),
		Type:           VoteTypePrecommit,
	}
	lastCommit := &Commit{
		Precommits: []*Vote{vote, nil},
	}
	return NewBlock(log.New(), &header, txns, nil, lastCommit)
}

func TestBlock_MakePartSet(t *testing.T) {
	block := CreateNewBlockWithTxs(1)

	partSetBytes := make([]byte, 0)
	partSet := block.MakePartSet(4000)
	for i:=uint(0); i<partSet.Total();i++ {
		partSetBytes = append(partSetBytes, partSet.GetPart(int(i)).Bytes.Bytes()...)
	}

	// decode block from partSetBytes
	newBlock := new(Block)
	err := rlp.DecodeBytes(partSetBytes, &newBlock)
	if err != nil {
		t.Error(err)
	}

	if !block.Hash().Equal(newBlock.Hash()) {
		t.Error("newBlock does not match with current block")
	}

	data, err := partSet.MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(data))

	//parts := partSet.BitArray()
}
