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

pragma solidity ^0.5.8;

contract DualMaster {

    string constant isMasterGenesisFunc = "isMasterGenesis(address)";
    string constant claimDualRewardFunc = "claimDualReward(address,address,uint64)";

    address constant PoSHandler = 0x0000000000000000000000000000000000000005;
    address constant Master = 0x0000000000000000000000000000000000000009;

    string _name;
    uint64 _consensusPeriod;
    uint64 _maxValidators;
    uint64 _maxViolatePercentage;

    struct NodeInfo {
        address node;
        address owner;
        uint256 stakes;
    }

    struct Validators {
        uint64 totalNodes;
        uint64 startAtBlock;
        uint64 endAtBlock;
        mapping(uint64=>NodeInfo) nodes;
        mapping(address=>uint64) addedNodes;
    }

    struct PendingDeleteInfo {
        NodeInfo node;
        uint64 index;
        uint64 vote;
        mapping(address=>bool) votedAddress;
        bool done;
    }

    struct RejectedVotes {
        uint64 totalVoted;
        bool status;
        mapping(address=>bool) voted;
    }

    struct CollectValidatorRequest {
        uint64 startAtBlock;
        uint64 vote;
        bool status;
        mapping(address=>bool) votedAddress;
    }

    struct ClaimRequest {
        uint64 blockHeight;
        address node;
        uint64 vote;
        bool claimed;
        mapping(address=>bool) votedAddress;
    }

    // _history contains all validators through period.
    Validators[] _history;
    NodeInfo[] _nodes;
    PendingDeleteInfo[] _pendingDeleteNodes;

    // _rejectedVote contains info of voting process of rejecting validation of a node for a specific block.
    // the first key is rejected block height, the second key is node's address.
    mapping(uint64=>mapping(address=>RejectedVotes)) _rejectedVotes;
    mapping(address=>uint64) _addedNodes;
    mapping(address=>address) _ownerNode; // return node's address based on node's owner address
    mapping(address=>bool) _genesises;
    mapping(address=>bool) _genesisOwners;
    mapping(uint64=>bool) _requestedStartBlock; // this variable marks if _startAtBlock has been requested to collect next validators.
    mapping(uint64=>CollectValidatorRequest) _requestCollectValidators;
    mapping(uint64=>bool) _rewardedBlock;
    mapping(uint64=>bool) _claimedRewardBlock;
    mapping(uint64=>ClaimRequest) _requestClaimReward;

    uint64 _totalGenesis = 0;

    // _startAtBlock stores started block in every consensusPeriod
    uint64 _startAtBlock = 1; // block 0 is genesis block

    // _nextBlock stores started block for the next consensusPeriod
    uint64 _nextBlock = 1;

    modifier isPoSHandler {
        require(msg.sender == PoSHandler, "sender is not PoSHandler");
        _;
    }

    modifier isMaster {
        require(msg.sender == Master, "sender is not master");
        _;
    }

    modifier isAvailableNode {
        bool result = false;
        for (uint64 i = 1; i < _nodes.length; i++) {
            if (_nodes[i].node == msg.sender || _nodes[i].owner == msg.sender) {
                result = true;
                break;
            }
        }
        require(result, "sender is not in available nodes");
        _;
    }

    modifier isMasterGenesis {
        (bool success, bytes memory result) = Master.staticcall(abi.encodeWithSignature(isMasterGenesisFunc, msg.sender));
        require(success, "check isMasterGenesis fail");

        bool rs = abi.decode(result, (bool));
        require(rs, "sender is not master's genesis");
        _;
    }

    modifier _isValidator {
        require(isValidator(msg.sender) || isDualGenesis(msg.sender), "sender is not validator");
        _;
    }

    constructor(string memory dualName, uint64 consensusPeriod, uint64 maxValidators, uint64 maxViolatePercentage) public {
        _name = dualName;
        _consensusPeriod = consensusPeriod;
        _maxValidators = maxValidators;
        _maxViolatePercentage = maxViolatePercentage;
        _nodes.push(NodeInfo(address(0x0), address(0x0), 0));
        _pendingDeleteNodes.push(PendingDeleteInfo(_nodes[0], 0, 0, true));
    }

    function setGenesis(address node, address owner, uint256 stakes) public isMasterGenesis {
        addNode(node, owner, stakes);
        _genesises[node] = true;
        _genesisOwners[owner] = true;
        _totalGenesis++;
    }

    function isDualGenesis(address nodeOrSender) public view returns (bool) {
        return _genesises[nodeOrSender] || _genesisOwners[nodeOrSender];
    }

    function addNode(address node, address owner, uint256 stakes) internal {
        _nodes.push(NodeInfo(node, owner, stakes));
        _addedNodes[node] = uint64(_nodes.length-1);
        _ownerNode[owner] = node;
        reIndexBackward(uint64(_nodes.length-1));
    }

    function join(address node, address owner, uint256 amount) public isMaster {
        addNode(node, owner, amount);
    }

    function reIndex(uint64 index) internal {
        require(index > 0, "invalid index");
        reIndexForward(index);
        reIndexBackward(index);
    }

    function reIndexForward(uint64 index) internal {
        while (index < _nodes.length-1) {
            if (_nodes[index].stakes < _nodes[index+1].stakes) {
                // swap 2 values
                NodeInfo memory n = _nodes[index+1];
                _nodes[index+1] = _nodes[index];
                _nodes[index] = n;

                // update _addedNodes
                _addedNodes[_nodes[index].node] = index;
                _addedNodes[_nodes[index+1].node] = index+1;

                index++;
            } else {
                return;
            }
        }
    }

    function reIndexBackward(uint64 index) internal {
        while (index > 1) {
            if (_nodes[index].stakes > _nodes[index-1].stakes) {
                // swap 2 values
                NodeInfo memory n = _nodes[index];
                _nodes[index] = _nodes[index-1];
                _nodes[index-1] = n;

                // update _addedNodes
                _addedNodes[_nodes[index].node] = index;
                _addedNodes[_nodes[index-1].node] = index-1;

                index--;
            } else {
                return;
            }
        }
    }

    function updateStakeAmount(address node, uint256 amount) public isMaster {
        uint64 index = _addedNodes[node];
        require(index > 0, "node is not found");
        NodeInfo storage nodeInfo = _nodes[index];
        nodeInfo.stakes = amount;
        reIndex(index);
    }

    // isQualified checks if vote count is greater than or equal with 2/3 total or not.
    function isQualified(uint64 count, uint64 total) internal pure returns (bool) {
        return count >= (total*2/3) + 1;
    }

    // requestCollectValidators creates by dual validators, if they are qualified (2/3+1 per number of validators), new validators set will be created.
    function requestCollectValidators() public _isValidator {
        uint64 total = _totalGenesis; // assign to _totalGenesis by default if _history is empty.
        NodeInfo memory nodeInfo = _nodes[_addedNodes[_ownerNode[msg.sender]]];

        if (!_requestedStartBlock[_nextBlock]) {
            // new request does not exist
            _requestCollectValidators[_nextBlock] = CollectValidatorRequest(_nextBlock, 0, false);
            _requestedStartBlock[_nextBlock] = true;
        }

        // if node has been voted, then revert.
        require(!_requestCollectValidators[_nextBlock].votedAddress[nodeInfo.node], "this node has voted");

        _requestCollectValidators[_nextBlock].vote++;
        _requestCollectValidators[_nextBlock].votedAddress[nodeInfo.node] = true;

        CollectValidatorRequest storage request = _requestCollectValidators[_nextBlock];

        if (_history.length > 0) {
            total = _history[_history.length-1].totalNodes-1;
        }

        if (isQualified(request.vote, total) && !request.status) {
            _requestCollectValidators[_nextBlock].status = true;
            collectValidators(_nextBlock);
        }
    }

    function initValidators() public _isValidator {
        collectValidators(1); // start at block 1
    }

    // collectValidators base on available nodes, max validators, collect validators and start new consensus period.
    // sometime, tx may be delayed for some blocks due to the traffic.
    // before adding new period, update last end block with current blockHeight
    // update _startAtBlock with current blockHeight + 1

    function collectValidators(uint64 nextStartBlock) internal {
        // update _startAtBlock and _nextBlock
        _startAtBlock = nextStartBlock;
        _nextBlock += _consensusPeriod+1;

        _history.push(Validators(1, _startAtBlock, _nextBlock-1));
        _history[_history.length-1].nodes[0] = _nodes[0];

        // get len based on _totalAvailableNodes and _maxValidators
        uint len = _nodes.length-1;
        if (len > _maxValidators) len = _maxValidators;
        // check valid nodes.
        for (uint64 i=1; i <= len; i++) {
            if (_nodes[i].stakes == 0) continue;
            uint64 currentIndex = _history[_history.length-1].totalNodes;
            _history[_history.length-1].nodes[currentIndex] = _nodes[i];
            _history[_history.length-1].addedNodes[_nodes[i].node] = currentIndex;
            _history[_history.length-1].totalNodes += 1;
        }
    }

    // isValidator checks an address whether it belongs into latest validator.
    function isValidator(address sender) public view returns (bool) {
        if (_history.length == 0) return false;
        Validators memory validators = _history[_history.length-1];
        for (uint64 i=1; i < validators.totalNodes; i++) {
            address owner = _history[_history.length-1].nodes[i].owner;
            address node = _history[_history.length-1].nodes[i].node;
            if (owner == sender || node == sender) {
                return true;
            }
        }
        return false;
    }

    function getLatestValidatorsInfo() public view returns (uint64 totalNodes, uint64 startAtBlock, uint64 endAtBlock) {
        if (_history.length == 0) return (0, 0, 0);
        return (_history[_history.length-1].totalNodes-1, _history[_history.length-1].startAtBlock, _history[_history.length-1].endAtBlock);
    }

    function getLatestValidatorByIndex(uint64 index) public view returns (address node, address owner, uint256 stakes) {
        (uint64 len, , ) = getLatestValidatorsInfo();
        require(index <= len, "invalid index");
        NodeInfo memory validator = _history[_history.length-1].nodes[index];
        return (validator.node, validator.owner, validator.stakes);
    }

    function requestClaimReward(address node, uint64 blockHeight) public _isValidator {
        require(!_rewardedBlock[blockHeight], "block has been rewarded");
        if (!_claimedRewardBlock[blockHeight]) {
            _requestClaimReward[blockHeight] = ClaimRequest(blockHeight, node, 0, false);
            _claimedRewardBlock[blockHeight] = true;
        }
        require(!_requestClaimReward[blockHeight].votedAddress[msg.sender], "sender has voted for this height");
        _requestClaimReward[blockHeight].votedAddress[msg.sender] = true;
        _requestClaimReward[blockHeight].vote++;
        uint64 totalNodes = _history[_history.length-1].totalNodes-1;
        if (isQualified(_requestClaimReward[blockHeight].vote, totalNodes) && !_requestClaimReward[blockHeight].claimed) {
            (bool success,) = PoSHandler.call(abi.encodeWithSignature(claimDualRewardFunc, address(this), node, blockHeight));
            require(success, "claim dual reward to PoSHandler failed");
            _requestClaimReward[blockHeight].claimed = true;
        }
    }

    function setRewarded(uint64 blockHeight) public isPoSHandler {
        _rewardedBlock[blockHeight] = true;
    }

    function isRewarded(uint64 blockHeight) public view returns (bool) {
        return _rewardedBlock[blockHeight];
    }
}
