pragma solidity ^0.8.7;

contract ERC20 {
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(
        address indexed tokenOwner,
        address indexed spender,
        uint256 value
    );
}
