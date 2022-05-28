require("@nomiclabs/hardhat-ethers");

// task action function receives the Hardhat Runtime Environment as second argument
task(
  "getWETH",
  "downloads WETH contract code",
  async (_, { ethers }) => {
    const wethAddress = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2";
	// Construction to get any contract as an object by its interface and address in blockchain
	// It is necessary to note that you must add an interface to your project
	const WETH = await ethers.getContractAt('IWETH', wethAddress);
	console.log(WETH);
  }
);


module.exports = {};