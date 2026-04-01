// SPDX-License-Identifier: MIT
pragma solidity ^0.8.31;

contract BinarySearch {

    function search(uint[] memory arr, uint target) public pure returns (int) {
        if (arr.length == 0) return -1;

        uint left = 0;
        uint right = arr.length - 1;

        while (left <= right) {
            uint mid = left + (right - left) / 2;

            if (arr[mid] == target) {
                return int(mid); // 找到返回下标
            } else if (arr[mid] < target) {
                left = mid + 1;
            } else {
                if (mid == 0) break;
                right = mid - 1;
            }
        }

        return -1; // 没找到
    }
}