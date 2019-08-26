package consensus

import (
	"errors"

	amino "github.com/kardiachain/go-kardia/lib/go-amino"
	"github.com/kardiachain/go-kardia/lib/rlp"
	"github.com/kardiachain/go-kardia/types"
)

const (
	MaxLimitBlockStore = 200     // Not use yet
	MaxBlockBytes      = 1048510 // lMB
)

var cdc = amino.NewCodec()

//MakePartSet  block to partset
func MakePartSet(partSize uint, block *types.Block) (*types.PartSet, error) {
	// Prefix the byte length, so that unmarshaling
	// can easily happen via a reader.
	bzs, err := rlp.EncodeToBytes(block)
	if err != nil {
		panic(err)
	}
	bz, err := cdc.MarshalBinary(bzs)
	if err != nil {
		return nil, err
	}
	return types.NewPartSetFromData(bz, partSize), nil
}

//MakeBlockFromPartSet partSet to block
func MakeBlockFromPartSet(reader *types.PartSet) (*types.Block, error) {
	if reader.IsComplete() {
		maxsize := int64(MaxBlockBytes)
		b := make([]byte, maxsize, maxsize)
		_, err := cdc.UnmarshalBinaryReader(reader.GetReader(), &b, maxsize)
		if err != nil {
			return nil, err
		}
		var block types.Block
		if err = rlp.DecodeBytes(b, &block); err != nil {
			return nil, err
		}
		return &block, nil
	}
	return nil, errors.New("Make block from partset not complete")
}
