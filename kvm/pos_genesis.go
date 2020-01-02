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

package kvm

import (
	"fmt"
	"github.com/kardiachain/go-kardia/kai/base"
	"github.com/kardiachain/go-kardia/kai/pos"
	"github.com/kardiachain/go-kardia/kai/state"
	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
	"math/big"
	"strings"
)

var maximumGasLimit = uint64(8000000)

func newGenesisVM(from common.Address, gasLimit uint64, st base.StateDB) *KVM {
	ctx := NewGenesisKVMContext(from, gasLimit)
	return NewKVM(ctx, st, Config{})
}

func InitGenesisConsensus(st *state.StateDB, gasLimit uint64, consensusInfo pos.ConsensusInfo) error {
	var (
		err error
		masterAbi, dualAbi abi.ABI
	)
	master := consensusInfo.Master
	// get first node owner to be the sender
	sender := master.Nodes.GenesisInfo[0].Owner
	// create master smart contract
	if err = createMaster(gasLimit, st, master, sender); err != nil {
		return err
	}
	if masterAbi, err = abi.JSON(strings.NewReader(consensusInfo.Master.ABI)); err != nil {
		return err
	}
	if dualAbi, err = abi.JSON(strings.NewReader(consensusInfo.DualMaster.ABI)); err != nil {
		return err
	}
	// create nodes
	if err = createGenesisNodes(gasLimit, st, master.Nodes, master.MinimumStakes, master.LockedPeriod, masterAbi, master.Address); err != nil {
		return err
	}
	// create stakers and stake them
	if err = createGenesisStakers(gasLimit, st, master.Stakers, masterAbi, master.Address); err != nil {
		return err
	}
	// create genesis dual
	if len(master.DualGenesis) > 0 && consensusInfo.DualMaster != nil {
		for _, dualGenesis := range master.DualGenesis {
			if err = createDualMaster(gasLimit, st, dualGenesis, consensusInfo.DualMaster, master, sender); err != nil {
				return err
			}
			if err = CollectDualValidators(gasLimit, st, dualGenesis.Address, sender, dualAbi); err != nil {
				return err
			}
		}
	}
	// start collect validators
	return CollectValidators(gasLimit, st, master.Address, sender, masterAbi)
}

func createMaster(gasLimit uint64, st *state.StateDB, master *pos.MasterInfo, sender common.Address) error {
	var (
		masterAbi abi.ABI
		err error
		input []byte
	)
	vm := newGenesisVM(sender, gasLimit, st)
	if masterAbi, err = abi.JSON(strings.NewReader(master.ABI)); err != nil {
		return err
	}
	if input, err = masterAbi.Pack("", master.ConsensusPeriodInBlock, master.MaxValidators, master.MaxViolatePercentageAllowed); err != nil {
		return err
	}
	newCode := append(master.ByteCode, input...)
	if _, _, _, err = InternalCreate(vm, &master.Address, newCode, master.GenesisAmount); err != nil {
		return err
	}
	return err
}

func createDualMaster(gasLimit uint64, st *state.StateDB, dualGenesis pos.DualGenesis, dualMaster *pos.DualMasterInfo, master *pos.MasterInfo, sender common.Address) error {
	var (
		dualMasterAbi abi.ABI
		err error
		input []byte
	)
	vm := newGenesisVM(sender, gasLimit, st)
	if dualMasterAbi, err = abi.JSON(strings.NewReader(dualMaster.ABI)); err != nil {
		return err
	}
	if input, err = dualMasterAbi.Pack("", dualGenesis.Name, dualMaster.ConsensusPeriodInBlock, dualMaster.MaxValidators, dualMaster.MaxViolatePercentageAllowed); err != nil {
		return err
	}
	newCode := append(dualGenesis.ByteCode, input...)
	if _, _, _, err = InternalCreate(vm, &dualGenesis.Address, newCode, big.NewInt(0)); err != nil {
		return err
	}
	// set genesis nodes
	if len(dualGenesis.GenesisNodes) > 0 {
		var (
			err error
			input []byte
			output []byte
			index *big.Int
			nodeInfo availableNode
			masterAbi abi.ABI
		)
		if masterAbi, err = abi.JSON(strings.NewReader(master.ABI)); err != nil {
			return err
		}

		for _, genesisNode := range dualGenesis.GenesisNodes {
			if input, err = masterAbi.Pack(methodGetAvailableNodeIndex, common.HexToAddress(genesisNode)); err != nil {
				return err
			}
			// get nodeIndex
			if output, err = StaticCall(vm, master.Address, input); err != nil {
				return err
			}
			if err = masterAbi.Unpack(&index, methodGetAvailableNodeIndex, output); err != nil {
				return err
			}
			if index.Uint64() == 0 {
				return fmt.Errorf(fmt.Sprintf("cannot find node:%v info", genesisNode))
			}
			if input, err = masterAbi.Pack(methodGetAvailableNode, index); err != nil {
				return err
			}
			if output, err = StaticCall(vm, master.Address, input); err != nil {
				return err
			}
			if err = masterAbi.Unpack(&nodeInfo, methodGetAvailableNode, output); err != nil {
				return err
			}
			// add node to dual genesis
			if input, err = dualMasterAbi.Pack(methodSetGenesis, nodeInfo.NodeAddress, nodeInfo.Owner, nodeInfo.Stakes); err != nil {
				return err
			}
			if _, err = InternalCall(vm, dualGenesis.Address, input, big.NewInt(0)); err != nil {
				return err
			}
		}
	}
	return err
}

func createGenesisNodes(gasLimit uint64, st *state.StateDB, nodes pos.Nodes, minimumStakes *big.Int, lockedPeriod uint64, masterAbi abi.ABI, masterAddress common.Address) error {
	nodeAbi, err := abi.JSON(strings.NewReader(nodes.ABI))
	if err != nil {
		return err
	}
	posHandlerVm := newGenesisVM(posHandlerAddress, gasLimit, st)
	for _, n := range nodes.GenesisInfo {
		input, err := nodeAbi.Pack("", masterAddress, n.PubKey, n.Name, n.RewardPercentage, lockedPeriod, minimumStakes)
		if err != nil {
			return err
		}
		newCode := append(nodes.ByteCode, input...)
		vm := newGenesisVM(n.Owner, gasLimit, st)
		if _, _, _, err = InternalCreate(vm, &n.Address, newCode, big.NewInt(0)); err != nil {
			return err
		}
		// add node to master
		if input, err = masterAbi.Pack(methodAddNode, n.Address); err != nil {
			return err
		}
		if _, err = InternalCall(posHandlerVm, masterAddress, input, big.NewInt(0)); err != nil {
			return err
		}
	}
	return nil
}

func createGenesisStakers(gasLimit uint64, st *state.StateDB, stakers pos.Stakers, masterAbi abi.ABI, masterAddress common.Address) error {
	var (
		err error
		stakerAbi abi.ABI
		input []byte
	)
	if stakerAbi, err = abi.JSON(strings.NewReader(stakers.ABI)); err != nil {
		return err
	}
	posHandlerVm := newGenesisVM(posHandlerAddress, gasLimit, st)
	for _, staker := range stakers.GenesisInfo {
		if input, err = stakerAbi.Pack("", masterAddress); err != nil {
			return err
		}
		newStakerCode := append(stakers.ByteCode, input...)
		vm := newGenesisVM(staker.Owner, gasLimit, st)
		if _, _, _, err = InternalCreate(vm, &staker.Address, newStakerCode, big.NewInt(0)); err != nil {
			return err
		}
		// add staker to master
		if input, err = masterAbi.Pack(methodAddStaker, staker.Address); err != nil {
			return err
		}
		if _, err = InternalCall(posHandlerVm, masterAddress, input, big.NewInt(0)); err != nil {
			return err
		}
		// stake to staker
		if input, err = stakerAbi.Pack("stake", staker.StakedNode); err != nil {
			return err
		}
		if _, err = InternalCall(vm, staker.Address, input, staker.StakeAmount); err != nil {
			return err
		}
	}
	return nil
}

func CollectValidators(gasLimit uint64, st *state.StateDB, masterAddress, sender common.Address, masterAbi abi.ABI) error {
	var (
		err error
		input []byte
	)
	vm := newGenesisVM(sender, gasLimit, st)
	if input, err = masterAbi.Pack(methodCollectValidators); err != nil {
		return err
	}
	_, err = InternalCall(vm, masterAddress, input, big.NewInt(0))
	return err
}

func CollectDualValidators(gasLimit uint64, st *state.StateDB, masterAddress, sender common.Address, masterAbi abi.ABI) error {
	var (
		err error
		input []byte
	)
	vm := newGenesisVM(sender, gasLimit, st)
	if input, err = masterAbi.Pack(methodInitValidators); err != nil { // start at block 1
		return err
	}
	_, err = InternalCall(vm, masterAddress, input, big.NewInt(0))
	return err
}
