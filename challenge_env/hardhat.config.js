require("@nomiclabs/hardhat-web3");

const CHAIN_IDS = {
  hardhat: 31337, // chain ID for hardhat testing
};

task("balance", "Prints an account's balance")
  .addParam("account", "The account's address")
  .setAction(async (taskArgs) => {
    const account = web3.utils.toChecksumAddress(taskArgs.account);
    const balance = await web3.eth.getBalance(account);

    console.log(web3.utils.fromWei(balance, "ether"), "ETH");
  });

module.exports = {
  defaultNetwork: "hardhat",
  networks: {
    hardhat: {
      forking: {
        url: "https://mainnet.infura.io/v3/50941666505f4e9cbf6fd10663e4b6f0",      
      }
    },
    local: {
      chainId: CHAIN_IDS.hardhat,
        url: "https://mainnet.infura.io/v3/50941666505f4e9cbf6fd10663e4b6f0",
        blockNumber: 14859860, // a specific block number with which you want to work
    },
    ctf: {
      url: "http://35.160.101.66:8545/a4c483cc-061c-4cc3-a4f2-1b681c5ab7f1",
      accounts: ["0xab52a8fec71c37b80507bc4cdb37989fe56ed4e23d10e064f15e080ae3696fa9"]
    }
  },
  solidity: {
    version: "0.8.0",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200
      }
    }
  },
  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts"
  },
  mocha: {
    timeout: 40000
  }
}

