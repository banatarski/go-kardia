pragma solidity ^0.4.24;
contract KardiaExchange {
    struct ReleaseInfo {
        string matchedOriginalTxId;
        string pair;
        string receiveAddress;
        string releaseTxId;
        uint256 releaseAmount;
        uint status;
    }

    struct ExchangeOrder {
        string pair;
        string fromAddress;
        string toAddress;
        string originalTxId;
        uint256 sellAmount;
        // receiveAmount is total amount to receive, unit is in destination currency
        uint256 receiveAmount;
        // availableAmount is unmatched Amount, unit is in source currency
        uint256 availableAmount;
        // need to add matched and unmatched amount here;
        uint done;
        bool added;
    }

    struct Rate {
        uint sellAmount;
        uint receiveAmount;
    }
    mapping (string => ReleaseInfo[]) releasesByID;
    mapping (string => ExchangeOrder) listOrders;
    mapping (string => string[]) orderIDsByPairs;
    mapping (string => Rate) rates;
    mapping (string => string[]) orderIDsByAddress;
    event Release(
        string indexed pair,
        string indexed addr,
        string matchOrderId,
        uint256 _value
);

    // if 1 eth = 10 neo, then sale_amount = 1, receive_amount = 10
    function addRate(string pair, uint sale_amount, uint receiveAmount) public {
        rates[pair] = Rate(sale_amount, receiveAmount);
    }
    // pair should be "ETH-NEO" for order from ETH to NEO and vice versa
    function getRate(string pair) internal view returns (Rate) {
        return rates[pair];
    }

    // pair should be "ETH-NEO" for order from ETH to NEO and vice versa
    function getRatePublic(string pair) public view returns (uint sale, uint receive) {
        return (rates[pair].sellAmount, rates[pair].receiveAmount);
    }

    // create an order with source - dest address, source - dest pair and amount
    // order will be stored with returned ID in smc
    // order from ETH from NEO should be "ETH-NEO", "NEO-ETH"
    function matchOrder(string srcPair, string destPair, string srcAddress, string destAddress, string originalTxId, uint256 amount) public {
        Rate memory r = getRate(srcPair);
        if (listOrders[originalTxId].added) {
            return;
        }
        if (r.receiveAmount == 0 || r.sellAmount == 0 ) {
            return;
        }
        uint256 receiveAmount = amount * r.receiveAmount / r.sellAmount;
        ReleaseInfo[] releases;
        uint256 totalMatched;
        ExchangeOrder memory order = ExchangeOrder(srcPair, srcAddress, destAddress, originalTxId, amount, receiveAmount, amount , 0, true);
        if (orderIDsByPairs[destPair].length > 0 ) {
            (releases, totalMatched) = findPartialMatchingOrders(originalTxId, srcPair, destPair, orderIDsByPairs[destPair], r, receiveAmount, destAddress);
            if (totalMatched > 0) {
                // calculate totalMatched back to original unit
                order.availableAmount = amount - (totalMatched * r.sellAmount / r.receiveAmount);
            }
        }
        for (uint i = 0 ; i < releases.length; i ++) {
            releasesByID[originalTxId].push(releases[i]);
        }
        // we need to add matched order of original direction here too update available amount of newly added order
        listOrders[originalTxId] = order;
        orderIDsByPairs[srcPair].push(originalTxId);
        orderIDsByAddress[srcAddress].push(originalTxId);
    }

    // find matching order for targeted pair with specific amount, return global order ID and total matched amount in dest currency
    function findPartialMatchingOrders(string originalTxId, string sourcePair, string destPair, string[] ids,  Rate r, uint256 receiveAmount, string receiveAddress) internal returns (ReleaseInfo[] storage releases, uint256 totalMatched ) {
        uint256 i = 0;
        while (i < ids.length && totalMatched < receiveAmount ) {
            // if sourcePair is ETH-NEO, matchableAmount unit is NEO
            uint256 matchableAmount = receiveAmount > listOrders[ids[i]].availableAmount ? listOrders[ids[i]].availableAmount : receiveAmount;
            ReleaseInfo memory r1 = ReleaseInfo(originalTxId, sourcePair, receiveAddress, "", matchableAmount, 0);
            releases.push(r1);
            uint256 matchableAmount2 = matchableAmount * r.sellAmount / r.receiveAmount;
            // update availableAmount of matched order, if matchOrder is NEO-ETH, availableAmount unit is NEO)
            listOrders[ids[i]].availableAmount -= matchableAmount;
            ReleaseInfo memory r2 = ReleaseInfo(ids[i], destPair, listOrders[ids[i]].toAddress, "", matchableAmount2, 0);
            releasesByID[ids[i]].push(r2);
            releases.push(r2);
            totalMatched += matchableAmount;
            i++;
        }
        return (releases, totalMatched);
    }

    // Get order id from its details:
    function getOrderId(string sourcePair, string fromAddress, string toAddress, uint256 amount) internal view returns (string) {
        string[] memory ids = orderIDsByPairs[sourcePair];
        for (uint256 i = 0; i < ids.length; i++) {
            ExchangeOrder memory order = listOrders[ids[i]];
            if (keccak256(abi.encodePacked(order.fromAddress)) == keccak256(abi.encodePacked(fromAddress)) && keccak256(abi.encodePacked(order.toAddress)) == keccak256(abi.encodePacked(toAddress)) && order.sellAmount == amount) {
                return ids[i];
            }
        }
        // no order found
        return "";
    }

    // Complete a order indicates that the order with orderID has been release successfully
    // TODO(@sontranrad): implement retry logic in case release is failed
    function completeOrder(string originalTxId, string releaseTxId, string pair, string receiveAddress, uint256 releaseAmount) public returns (uint success){
        ReleaseInfo[] releases = releasesByID[originalTxId];
        if (releases.length == 0) {
            return 0;
        }
        for (uint i = 0; i < releases.length; i++) {
            if (keccak256(abi.encodePacked(pair)) == keccak256(abi.encodePacked(releases[i].pair)) &&
                keccak256(abi.encodePacked(receiveAddress)) == keccak256(abi.encodePacked(releases[i].receiveAddress)) &&
                releaseAmount == releases[i].releaseAmount) {
                releases[i].releaseTxId = releaseTxId;
                releases[i].status = 1;
                string matchedTxId = releases[i].matchedOriginalTxId;
                // ReleaseInfo[] storage opponents = releasesByID[matchedTxId];
                if (releasesByID[matchedTxId].length > 0) {
                    for (uint j = 0; j < releasesByID[matchedTxId].length; j++) {
                        if (keccak256(abi.encodePacked(releasesByID[matchedTxId][j].receiveAddress))
                            == keccak256(abi.encodePacked(receiveAddress)) && releasesByID[matchedTxId][j].releaseAmount == releaseAmount
                            && releasesByID[matchedTxId][j].status == 0) {
                            releasesByID[matchedTxId][0].releaseTxId = releaseTxId;
                            releasesByID[matchedTxId][0].status = 1;
                        }
                    }
                }
            }
        }

    }

    // Get 10 matchable amount by pair
    function getMatchableAmount(string pair) public view returns (uint256[] amounts) {
        string[] memory ids = orderIDsByPairs[pair];
        amounts = new uint256[](10);
        for (uint i = 0; i < ids.length; i++) {
            if (i < 10) {
                amounts[i] = listOrders[ids[i]].receiveAmount;
            }
        }
        return amounts;
    }

    // Get release by originalTxId and pair
    function getReleaseByTxId(string originalTxId, string pair) public view returns (string releaseInfos) {
        ReleaseInfo[] memory releases = releasesByID[originalTxId];
        if (releases.length == 0) {
            releaseInfos = "";
            return releaseInfos;
        }
        for (uint i = 0; i < releases.length; i++) {
            if (keccak256(abi.encodePacked(releases[i].pair)) == keccak256(abi.encodePacked(pair))) {
                if (keccak256(abi.encodePacked(releaseInfos)) == keccak256(abi.encodePacked(""))) {
                    releaseInfos = releaseInfoToString(releases[i]);
                } else {
                    releaseInfos = string(abi.encodePacked(releaseInfos, "|", releaseInfoToString(releases[i])));
                }
            }
        }
        return releaseInfos;
    }

    function orderToString(ExchangeOrder r ) internal pure returns (string) {
        if (r.added = false) {
            return "";
        }
        string memory sell = uint2str(r.sellAmount);
        string memory receive = uint2str(r.receiveAmount);
        string memory done = uint2str(r.done);
        string memory available = uint2str(r.availableAmount);
        return string(abi.encodePacked(r.pair, ";", r.fromAddress, ";", r.toAddress, ";", r.originalTxId,";", sell, ";", receive, ";", available, ";", done));
    }

    function releaseInfoToString(ReleaseInfo r ) internal pure returns (string) {
        string memory amount = uint2str(r.releaseAmount);
        string memory status = uint2str(r.status);
        return string(abi.encodePacked(r.matchedOriginalTxId, ":", r.pair, ";", r.receiveAddress, ";", r.releaseTxId, ";", amount,";", status));
    }

    function uint2str(uint i) internal pure returns (string){
        if (i == 0) return "0";
        uint j = i;
        uint length;
        while (j != 0){
            length++;
            j /= 10;
        }
        bytes memory bstr = new bytes(length);
        uint k = length - 1;
        while (i != 0){
            bstr[k--] = byte(48 + i % 10);
            i /= 10;
        }
        return string(bstr);
    }

    function getOrderBook(string pair) public view returns (string orderBook) {
        string[] memory ids = orderIDsByPairs[pair];
        uint length = ids.length;
        if (length == 0) {
            return "";
        }
        string memory strResult = "";
        for (uint i = 0; i < length; i ++) {
            if (i == 0) {
                strResult = orderToString(listOrders[ids[i]]);
            }
            else {
                strResult = string(abi.encodePacked(strResult, "|", orderToString(listOrders[ids[i]])));
            }
        }
        return strResult;
    }

    // get all order history of an address
    function getOrderHistory(string addr) public view returns (string orderHistory) {
        string[] memory ids = orderIDsByAddress[addr];
        uint length = ids.length;
        if (length == 0) {
            return "";
        }
        string memory strResult = "";
        for (uint i = 0; i < length; i ++) {
            if (i == 0) {
                strResult = orderToString(listOrders[ids[i]]);
            }
            else {
                strResult = string(abi.encodePacked(strResult, "|", orderToString(listOrders[ids[i]])));
            }
        }
        return strResult;
    }

    // Get all order history of an address and a pair
    function getOrderHistoryByPair(string addr, string pair) public view returns (string orderHistory) {
        string[] memory ids = orderIDsByAddress[addr];
        uint length = ids.length;
        if (length == 0) {
            return "";
        }
        string memory strResult = "";
        bytes32 encodedPair = keccak256(abi.encodePacked(pair));
        for (uint i = 0; i < length; i ++) {
            if (keccak256(abi.encodePacked(listOrders[ids[i]].pair)) == encodedPair) {
                if (i == 0) {
                    strResult = orderToString(listOrders[ids[i]]);
                }
                else {
                    strResult = string(abi.encodePacked(strResult, "|", orderToString(listOrders[ids[i]])));
                }
            }
        }
        return strResult;
    }
}
