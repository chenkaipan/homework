// SPDX-License-Identifier: MIT
pragma solidity 0.8.31;

contract IntToRoman {

    function intToRoman(uint256 num) public pure returns (string memory) {
        require(num >= 1 && num <= 3999, "Out of range");

        uint256[13] memory values = [
            uint256(1000), 900, 500, 400,
            100, 90, 50, 40,
            10, 9, 5, 4, 1
        ];

        string[13] memory symbols = [
            "M", "CM", "D", "CD",
            "C", "XC", "L", "XL",
            "X", "IX", "V", "IV", "I"
        ];

        bytes memory result;

        for (uint256 i = 0; i < values.length; i++) {
            while (num >= values[i]) {
                num -= values[i];
                result = abi.encodePacked(result, symbols[i]);
            }
        }

        return string(result);
    }
}