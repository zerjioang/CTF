// SPDX-License-Identifier: private

pragma solidity 0.8.0;

import "./Setup.sol";

contract FlashTest {

    FlashLoan public flashloanPool;
    Setup public setup;
    WETH9 public weth;
    bool public loanDone;

    // we initialize all previous variables in the constructor
    // based on our development environment data
    constructor(Setup _setup) {
        // 1 extract from existing Setup contract all available data to be in sync
        setup = _setup;
        weth = setup.weth();
        flashloanPool = setup.flashloanPool();
    }

    function attack() public {
        // now request a flash loan an borrow some third party tokens
        flashloanPool.flashLoan(10 ether); // max 1000 ETH
        // after requesting flashloan 1000 eth will be transfered to this exploit balance
        // and callback method receiveFlashLoan(uint256) will be executed
        // ** receiveFlashLoan is called **/
    }

    /**
    receiveFlashLoan is called each time a flashloan is received
     */
    function receiveFlashLoan(uint256 amount) public {
        // finally return the flash loan
        bool success = weth.transfer(address(flashloanPool), amount);
        require(success, "transfer loan back failed");
        loanDone = true;
    }
}