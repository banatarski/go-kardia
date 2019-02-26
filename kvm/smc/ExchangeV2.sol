pragma solidity ^0.4.24;
contract KardiaExchange {
    struct ExchangeOrder {
        string pair;
        string fromAddress;
        string toAddress;
        string originalTxId;
        uint256 sellAmount;
        uint256 receiveAmount;
        string matchedOrderId;
        uint done;
        bool added;
    }

    struct Rate {
        uint sellAmount;
        uint receiveAmount;
    }

    mapping (string => ExchangeOrder) listOrders;
    mapping (string => string[]) orderIDsByPairs;
    mapping (string => Rate) rates;
    mapping (string => string[]) orderIDsByAddress;
    uint256 counter;
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
    function matchOrder(string srcPair, string destPair, string srcAddress, string destAddress, string originalTxId, uint256 amount) public returns (uint256) {
        Rate memory r = getRate(srcPair);
        if (r.receiveAmount == 0 || r.sellAmount == 0 ) {
            return 0;
        }
        ++counter;
        uint256 receiveAmount = amount * r.receiveAmount / r.sellAmount;
        string memory matchOrderId = findMatchingOrder(destPair, receiveAmount);
        ExchangeOrder memory order = ExchangeOrder(srcPair, srcAddress, destAddress, originalTxId, amount, receiveAmount, matchOrderId, 0, true);
        listOrders[originalTxId] = order;
        orderIDsByPairs[srcPair].push(originalTxId);
        orderIDsByAddress[srcAddress].push(originalTxId);
        if (keccak256(abi.encodePacked(order.matchedOrderId)) != keccak256(abi.encodePacked("")) ) {
            listOrders[matchOrderId].matchedOrderId = originalTxId;
            emit Release(destPair, listOrders[matchOrderId].toAddress, matchOrderId, listOrders[matchOrderId].receiveAmount);
        }
        return counter;
    }

    // find matching order for targeted pair with specific amount, return global order ID
    function findMatchingOrder(string destPair, uint256 receiveAmount) internal view returns (string) {
        string[] memory ids = orderIDsByPairs[destPair];
        for (uint256 i = 0; i < ids.length; i++) {
            ExchangeOrder memory order = listOrders[ids[i]];
            if (order.sellAmount == receiveAmount
                && keccak256(abi.encodePacked(order.matchedOrderId)) == keccak256(abi.encodePacked(""))
                && order.done == 0) {
                return ids[i];
            }
        }
        // no matching order ID found
        return "";
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

    // Get matched order detail for a specific order
    function getMatchingOrderInfo(string orderID) public view returns (string matchedOrderID, string destAddress, uint256 sendAmount) {
        if (listOrders[orderID].added == false) {
            return ("", "", 0);
        }

        string memory matchedId = listOrders[orderID].matchedOrderId;
        if (keccak256(abi.encodePacked(matchedId)) != keccak256(abi.encodePacked(""))) {
            ExchangeOrder memory matchedOrder = listOrders[matchedId];
            return (matchedId, matchedOrder.toAddress, matchedOrder.receiveAmount);

        }
        // no matching order found
        return ("", "", 0);
    }

    // Complete a order indicates that the order with orderID has been release successfully
    // TODO(@sontranrad): implement retry logic in case release is failed
    function completeOrder(string orderID, string pair) public returns (uint success){
        ExchangeOrder memory order = listOrders[orderID];
        if (keccak256(abi.encodePacked(order.pair)) != keccak256(abi.encodePacked(pair)) ) {
            return 0;
        }
        if (listOrders[orderID].done != 0) {
            return 0;
        }
        listOrders[orderID].done = 1;
        return 1;
    }

    // Get exchangeable amount of each pair
    function getAvailableAmountByPair(string pair) public view returns (uint256 amount) {
        amount = 0;
        string[] memory ids = orderIDsByPairs[pair];
        bytes32 encodedEmptyString = keccak256(abi.encodePacked(""));
        for (uint i = 0; i < ids.length; i++) {
            if (keccak256(abi.encodePacked(listOrders[ids[i]].matchedOrderId)) == encodedEmptyString)
                amount += listOrders[ids[i]].receiveAmount;
        }
        return amount;
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

    // Get the opposite uncompleted order of a compelete order
    function getUncompletedMatchingOrder(string orderID) public view returns (string matchedOrderID, string destAddress, uint256 sendAmount) {
        string memory oppositeID = listOrders[orderID].matchedOrderId;
        if (keccak256(abi.encodePacked(oppositeID)) == keccak256(abi.encodePacked("")) ) {
            return ("", "", 0);
        }
        ExchangeOrder memory oppositeOrder = listOrders[oppositeID];
        if (oppositeOrder.done == 0) {
            return (oppositeID, oppositeOrder.toAddress, oppositeOrder.receiveAmount);
        }
    }

    function orderToString(ExchangeOrder r ) internal pure returns (string) {
        if (r.added = false) {
            return "";
        }
        string memory sell = uint2str(r.sellAmount);
        string memory receive = uint2str(r.receiveAmount);
        string memory done = uint2str(r.done);
        return string(abi.encodePacked(r.fromAddress, "|", r.toAddress, "|", r.originalTxId,"|", sell, "|", receive,
            "|", r.matchedOrderId, "|", done));
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

    function getOrderBook(string pair) public view returns (string) {
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
    function getOrderHistory(string addr) public view returns (string) {
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
    function getOrderHistoryByPair(string addr, string pair) public view returns (string) {
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
