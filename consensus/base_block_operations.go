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

package consensus

import (
	"github.com/kardiachain/go-kardia/kai/state"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/types"
)

type BaseBlockOperations interface {
	Height() uint64
	LoadBlock(height uint64) *types.Block
	LoadBlockMeta(height uint64) *types.BlockMeta
	LoadBlockPart(height uint64, index int) *types.Part
	LoadBlockCommit(height uint64) *types.Commit
	LoadSeenCommit(height uint64) *types.Commit
	CreateProposalBlock(height int64, cs state.LastestBlockState, commit *types.Commit, proposer common.Address) (block *types.Block, blockPart *types.PartSet)
	CommitAndValidateBlockTxs(block *types.Block) error
	CommitBlockTxsIfNotFound(block *types.Block) error
	SaveBlock(block *types.Block, blockParts *types.PartSet, seenCommit *types.Commit)
}
