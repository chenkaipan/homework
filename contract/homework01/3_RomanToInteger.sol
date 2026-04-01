// SPDX-License-Identifier: MIT
pragma solidity 0.8.31;

contract RomanToInteger {

    function romanToInt(string memory s) public pure returns (uint256) {
        bytes memory str = bytes(s);
        uint256 total = 0;

        for (uint256 i = 0; i < str.length; i++) {
            uint256 value = getValue(str[i]);

            // 如果不是最后一个，并且当前 < 下一个
            if (i < str.length - 1 && value < getValue(str[i + 1])) {
                total -= value;
            } else {
                total += value;
            }
        }

        return total;
    }

    function getValue(bytes1 char) internal pure returns (uint256) {
        if (char == "I") return 1;
        if (char == "V") return 5;
        if (char == "X") return 10;
        if (char == "L") return 50;
        if (char == "C") return 100;
        if (char == "D") return 500;
        if (char == "M") return 1000;
        revert("Invalid Roman numeral");
    }
}