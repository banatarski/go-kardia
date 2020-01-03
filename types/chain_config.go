package types

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/kardiachain/go-kardia/lib/common"
	"time"
)

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	// BaseAccount is used to set default execute account for
	*BaseAccount         `json:"baseAccount,omitempty"`
	Kaicon *ConsensusConfig     `json:"consensusConfig,omitempty"`
}

// BaseAccount defines information for base (root) account that is used to execute internal smart contract
type BaseAccount struct {
	Address common.Address       `json:"address"`
	PrivateKey ecdsa.PrivateKey
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var engine interface{}
	switch {
	case c.Kaicon != nil:
		engine = c.Kaicon
	default:
		engine = "unknown"
	}
	return fmt.Sprintf("{Engine: %v}",
		engine,
	)
}

func (c *ChainConfig) SetBaseAccount(baseAccount *BaseAccount) {
	c.BaseAccount = baseAccount
}

// -------- Consensus Config ---------

// ConsensusConfig defines the configuration for the Kardia consensus service,
// including timeouts and details about the block structure.
type ConsensusConfig struct {
	// All timeouts are in milliseconds
	TimeoutPropose        uint64 `mapstructure:"timeout_propose"`
	TimeoutProposeDelta   uint64 `mapstructure:"timeout_propose_delta"`
	TimeoutPrevote        uint64 `mapstructure:"timeout_prevote"`
	TimeoutPrevoteDelta   uint64 `mapstructure:"timeout_prevote_delta"`
	TimeoutPrecommit      uint64 `mapstructure:"timeout_precommit"`
	TimeoutPrecommitDelta uint64 `mapstructure:"timeout_precommit_delta"`
	TimeoutCommit         uint64 `mapstructure:"timeout_commit"`

	// Make progress as soon as we have all the precommits (as if TimeoutCommit = 0)
	SkipTimeoutCommit bool `mapstructure:"skip_timeout_commit"`

	// EmptyBlocks mode and possible interval between empty blocks in seconds
	CreateEmptyBlocks         bool          `mapstructure:"create_empty_blocks"`
	CreateEmptyBlocksInterval uint64 `mapstructure:"create_empty_blocks_interval"`

	// Reactor sleep duration parameters are in milliseconds
	PeerGossipSleepDuration     uint64 `mapstructure:"peer_gossip_sleep_duration"`
	PeerQueryMaj23SleepDuration uint64 `mapstructure:"peer_query_maj23_sleep_duration"`
}

// WaitForTxs returns true if the consensus should wait for transactions before entering the propose step
func (cfg *ConsensusConfig) WaitForTxs() bool {
	return !cfg.CreateEmptyBlocks || cfg.CreateEmptyBlocksInterval > 0
}

func (cfg *ConsensusConfig) getTimeoutCommit() time.Duration {
	return time.Duration(cfg.TimeoutCommit) * time.Millisecond
}

func (cfg *ConsensusConfig) getTimeoutPropose() time.Duration {
	return time.Duration(cfg.TimeoutPropose) * time.Millisecond
}

func (cfg *ConsensusConfig) getTimeoutProposeDelta() time.Duration {
	return time.Duration(cfg.TimeoutProposeDelta) * time.Millisecond
}

func (cfg *ConsensusConfig) getTimeoutPrevote() time.Duration {
	return time.Duration(cfg.TimeoutPrevote) * time.Millisecond
}

func (cfg *ConsensusConfig) getTimeoutPrecommit() time.Duration {
	return time.Duration(cfg.TimeoutPrecommit) * time.Millisecond
}

func (cfg *ConsensusConfig) getTimeoutPrecommitDelta() time.Duration {
	return time.Duration(cfg.TimeoutPrecommitDelta) * time.Millisecond
}

// Commit returns the amount of time to wait for straggler votes after receiving +2/3 precommits for a single block (ie. a commit).
func (cfg *ConsensusConfig) Commit(t time.Time) time.Time {
	return t.Add(cfg.getTimeoutCommit())
}

// Propose returns the amount of time to wait for a proposal
func (cfg *ConsensusConfig) Propose(round int) time.Duration {
	return time.Duration(
		cfg.getTimeoutPropose().Nanoseconds()+cfg.getTimeoutProposeDelta().Nanoseconds()*int64(round),
	) * time.Nanosecond
}

// Prevote returns the amount of time to wait for straggler votes after receiving any +2/3 prevotes
func (cfg *ConsensusConfig) Prevote(round int) time.Duration {
	return time.Duration(
		cfg.getTimeoutPrevote().Nanoseconds()+cfg.getTimeoutProposeDelta().Nanoseconds()*int64(round),
	) * time.Nanosecond
}

// Precommit returns the amount of time to wait for straggler votes after receiving any +2/3 precommits
func (cfg *ConsensusConfig) Precommit(round int) time.Duration {
	return time.Duration(
		cfg.getTimeoutPrecommit().Nanoseconds()+cfg.getTimeoutPrecommitDelta().Nanoseconds()*int64(round),
	) * time.Nanosecond
}

// PeerGossipSleep returns the amount of time to sleep if there is nothing to send from the ConsensusReactor
func (cfg *ConsensusConfig) PeerGossipSleep() time.Duration {
	return time.Duration(cfg.PeerGossipSleepDuration) * time.Millisecond
}

// PeerQueryMaj23Sleep returns the amount of time to sleep after each VoteSetMaj23Message is sent in the ConsensusReactor
func (cfg *ConsensusConfig) PeerQueryMaj23Sleep() time.Duration {
	return time.Duration(cfg.PeerQueryMaj23SleepDuration) * time.Millisecond
}

func (cfg *ConsensusConfig) GetCreateEmptyBlocksInterval() time.Duration {
	return time.Duration(cfg.CreateEmptyBlocksInterval) * time.Millisecond
}
