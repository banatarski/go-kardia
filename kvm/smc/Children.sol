pragma solidity ^0.4.24;

import "github.com/Arachnid/solidity-stringutils/strings.sol";

contract Children {
    using strings for *;

    mapping(uint => bool) isSet;
    mapping(uint => string[]) listAttributes;
    mapping(uint => string[]) listValues;
    function splitString(string str, string delimiter) internal returns (string[]) {
        strings.slice memory s = str.toSlice();
        strings.slice memory delim = delimiter.toSlice();
        string[] memory parts = new string[](s.count(delim) + 1);
        for(uint i = 0; i < parts.length; i++) {
            parts[i] = s.split(delim).toString();
        }
        return parts;
    }

    function addChild(uint id, string attributes, string values) public {
        string[] memory arrAttr = splitString(attributes, "|");
        string[] memory arrValues = splitString(values, "|");
        if (arrAttr.length != arrValues.length) {
            revert();
        }
        listAttributes[id] = arrAttr;
        listValues[id] = arrValues;
        isSet[id] = true;
    }

    function childToString(uint id) internal returns (string) {
        string memory result = string(abi.encodePacked("id:", uint2str(id), "|"));
        if (listAttributes[id].length == 0) {
            return "";
        }
        for (uint i = 0; i < listAttributes[id].length; i++) {
            string memory field = string(abi.encodePacked(listAttributes[id][i], ":",
                listValues[id][i]));
            if (i < listAttributes[id].length - 1) {
                result = string(abi.encodePacked(result, field, "|"));
            } else {
                result = string(abi.encodePacked(result, field));
            }
        }
        return result;
    }


    // getChildInfo returns info of a child by id
    function getChildInfo(uint id) public view returns (string childInfo) {
        return childToString(id);
    }

    // getChildren returns info of children from from_id to to_id, separated by ~
    function getChildren(uint from_id, uint to_id) public view returns (string childrenInfo) {
        if (from_id > to_id) {
            return "";
        }
        string memory result = "";
        for (uint i = from_id; i <= to_id; i++) {
            if (i < to_id) {
                result = string(abi.encodePacked(result, childToString(i), "~"));
            }
            else {
                result = string(abi.encodePacked(result, childToString(i)));
            }
        }
        return result;
    }
    // convert a uint type into string
    function uint2str(uint i) internal pure returns (string) {
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
}
