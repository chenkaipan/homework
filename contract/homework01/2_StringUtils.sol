// SPDX-License-Identifier: MIT
pragma solidity 0.8.31;

contract StringUtils {

    function reverse(string memory str) public pure returns (string memory) {
        bytes memory strBytes = bytes(str);
        uint256 length = strBytes.length;

        bytes memory reversed = new bytes(length);

        for (uint256 i = 0; i < length; i++) {
            reversed[i] = strBytes[length - 1 - i];
        }

        return string(reversed);
    }
}