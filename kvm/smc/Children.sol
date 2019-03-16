pragma solidity ^0.4.24;

contract Children {
    mapping(uint => bool) isSet;
    mapping(uint => string) listAttributes;
    mapping(uint => string) listValues;

    function addChildren(uint id, string atts, string values) public {
        isSet[id] = true;
        listAttributes[id] = atts;
        listValues[id] = values;
    }

    function getChildren(uint _from, uint _to) public view returns (string atts, string values) {
        if (_from < _to) {
            return ("","");
        }
        string memory strAttrs;
        string memory strVals;
        for (uint i = _from; i <= _to; i++) {
            if (isSet[i]) {
                if (i < _to) {
                    strAttrs = string(abi.encodePacked(strAttrs, listAttributes[i],"|"));
                    strVals = string(abi.encodePacked(strAttrs, listValues[i], "|"));
                } else {
                    strAttrs = string(abi.encodePacked(strAttrs, listAttributes[i]));
                    strVals = string(abi.encodePacked(strAttrs, listValues[i]));
                }
            }
        }
        return (strAttrs, strVals);
    }
}
