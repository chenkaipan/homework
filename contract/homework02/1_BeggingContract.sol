// SPDX-License-Identifier: MIT
pragma solidity ^0.8.31;

contract BeggingContract {
    //每个捐赠者的捐赠金额
    mapping(address => uint256) public donationAmount;
    //总金额
    uint256 totalAmount;

    //所有者
    address public owner;

    //事件-捐赠信息
    event DonationReceived(address indexed donor, uint256 amount);

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can withdraw");
        _;
    }
    //初始化所有者
    constructor() {
        owner = msg.sender;
    }
    //捐赠函数
    function donate() public payable {
        donationAmount[msg.sender] += msg.value;
        totalAmount += msg.value;
        //记录
        emit DonationReceived(msg.sender, msg.value);
    }

    //查询某个地址的捐赠金额
    function getDonation(address addr) public view returns (uint256) {
        return donationAmount[addr];
    }
    //提取函数
    function withdraw() public onlyOwner {
        uint256 amount = totalAmount;
        totalAmount = 0;

        (bool success, ) = payable(msg.sender).call{value: amount}("");
        require(success, "Transfer failed");
    }
    //直接查合约地址的余额
    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }
}
