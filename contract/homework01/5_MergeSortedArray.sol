// SPDX-License-Identifier: MIT
pragma solidity ^0.8.31;

contract MergeSortedArray {

    function merge(uint[] memory a, uint[] memory b) public pure returns (uint[] memory) {
        uint i = 0;
        uint j = 0;
        uint k = 0;

        uint[] memory result = new uint[](a.length + b.length);

        while (i < a.length && j < b.length) {
            if (a[i] <= b[j]) {
                result[k] = a[i];
                i++;
            } else {
                result[k] = b[j];
                j++;
            }
            k++;
        }

        // 把 a 剩余的补进去
        while (i < a.length) {
            result[k] = a[i];
            i++;
            k++;
        }

        // 把 b 剩余的补进去
        while (j < b.length) {
            result[k] = b[j];
            j++;
            k++;
        }

        return result;
    }
}