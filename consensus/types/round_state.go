package types

import (
	"fmt"
	"math/big"
	"time"

	cmn "github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/types"
)

//-----------------------------------------------------------------------------
// RoundStepType enum type

// RoundStepType enumerates the state of the consensus state machine
type RoundStepType uint8 // These must be numeric, ordered.

// RoundStepType
const (
	RoundStepNewHeight     = RoundStepType(0x01) // Wait til CommitTime + timeoutCommit
	RoundStepNewRound      = RoundStepType(0x02) // Setup new round and go to RoundStepPropose
	RoundStepPropose       = RoundStepType(0x03) // Did propose, gossip proposal
	RoundStepPrevote       = RoundStepType(0x04) // Did prevote, gossip prevotes
	RoundStepPrevoteWait   = RoundStepType(0x05) // Did receive any +2/3 prevotes, start timeout
	RoundStepPrecommit     = RoundStepType(0x06) // Did precommit, gossip precommits
	RoundStepPrecommitWait = RoundStepType(0x07) // Did receive any +2/3 precommits, start timeout
	RoundStepCommit        = RoundStepType(0x08) // Entered commit state machine
	// NOTE: RoundStepNewHeight acts as RoundStepCommitWait.
)

// String returns a string
func (rs RoundStepType) String() string {
	switch rs {
	case RoundStepNewHeight:
		return "RoundStepNewHeight"
	case RoundStepNewRound:
		return "RoundStepNewRound"
	case RoundStepPropose:
		return "RoundStepPropose"
	case RoundStepPrevote:
		return "RoundStepPrevote"
	case RoundStepPrevoteWait:
		return "RoundStepPrevoteWait"
	case RoundStepPrecommit:
		return "RoundStepPrecommit"
	case RoundStepPrecommitWait:
		return "RoundStepPrecommitWait"
	case RoundStepCommit:
		return "RoundStepCommit"
	default:
		return "RoundStepUnknown" // Cannot panic.
	}
}

//-----------------------------------------------------------------------------

// RoundState defines the *cmn.BigInternal consensus state.
// NOTE: Not thread safe. Should only be manipulated by functions downstream
// of the cs.receiveRoutine
type RoundState struct {
	Height         *cmn.BigInt         `json:"height"` // Height we are working on
	Round          *cmn.BigInt         `json:"round"`
	Step           RoundStepType       `json:"step"`
	StartTime      *big.Int            `json:"start_time"`
	CommitTime     *big.Int            `json:"commit_time"` // Subjective time when +2/3 precommits for Block at Round were found
	Validators     *types.ValidatorSet `json:"validators"`  // TODO(huny@): Assume static validator set for now
	Proposal       *types.Proposal     `json:"proposal"`
	ProposalBlock  *types.Block        `json:"proposal_block"`
	LockedRound    *cmn.BigInt         `json:"locked_round"`
	LockedBlock    *types.Block        `json:"locked_block"`
	ValidRound     *cmn.BigInt         `json:"valid_round"` // Last known round with POL for non-nil valid block.
	ValidBlock     *types.Block        `json:"valid_block"` // Last known block of POL mentioned above.
	Votes          *HeightVoteSet      `json:"votes"`
	CommitRound    *cmn.BigInt         `json:"commit_round"` //
	LastCommit     *types.VoteSet      `json:"last_commit"`  // Last precommits at Height-1
	LastValidators *types.ValidatorSet `json:"last_validators"`
}

// RoundStateEvent returns the H/R/S of the RoundState as an event.
func (rs *RoundState) RoundStateEvent() types.EventDataRoundState {
	// XXX: copy the RoundState
	// if we want to avoid this, we may need synchronous events after all
	rsCopy := *rs
	edrs := types.EventDataRoundState{
		Height:     rs.Height,
		Round:      rs.Round,
		Step:       rs.Step.String(),
		RoundState: &rsCopy,
	}
	return edrs
}

func (rs *RoundState) String() string {
	return fmt.Sprintf("RoundState{H:%v R:%v S:%v  StartTime:%v  CommitTime:%v  Validators:%v   Proposal:%v  ProposalBlock:%v  LockedRound:%v  LockedBlock:%v  ValidRound:%v  ValidBlock:%v  Votes:%v  LastCommit:%v  LastValidators:%v}",
		rs.Height, rs.Round, rs.Step,
		time.Unix(rs.StartTime.Int64(), 0),
		time.Unix(rs.CommitTime.Int64(), 0),
		rs.Validators,
		rs.Proposal,
		rs.ProposalBlock,
		rs.LockedRound,
		rs.LockedBlock,
		rs.ValidRound,
		rs.ValidBlock,
		rs.Votes,
		rs.LastCommit,
		rs.LastValidators)
}