// SPDX-License-Identifier: GPL-3.0
pragma solidity 0.8.31;

contract Voting {
    // 候选人 => 票数
    mapping(address => uint256) public votes;

    // 候选人名单
    address[] public candidates;

    // 是否是合法候选人
    mapping(address => bool) public isCandidate;

    // 每人只能投一次（可选）
    mapping(address => bool) public hasVoted;

    //初始化和添加候选人
    function addCandidate(address candidate) public {
        require(!isCandidate[candidate], "Already a candidate");

        candidates.push(candidate);
        isCandidate[candidate] = true;
    }
    //删除
    function deleteCandidate(address candidate) public {
        require(isCandidate[candidate], "Not a candidate");

        isCandidate[candidate] = false;
        delete votes[candidate];

        for (uint256 i = 0; i < candidates.length; i++) {
            if (candidates[i] == candidate) {
                // 用最后一个元素替换
                candidates[i] = candidates[candidates.length - 1];
                // 删除最后一个
                candidates.pop();
                break;
            }
        }
    }

    //投票
    function vote(address candidate) public {
        require(isCandidate[candidate], "Not a valid candidate");
        require(!hasVoted[msg.sender], "Already voted");

        votes[candidate] += 1;
        hasVoted[msg.sender] = true;
    }
    //查询票数
    function getVotes(address candidate) public view returns (uint256) {
        return votes[candidate];
    }

    //重置投票
    function resetVotes() public {
        for (uint256 i = 0; i < candidates.length; i++) {
            votes[candidates[i]] = 0;
        }
    }
}
