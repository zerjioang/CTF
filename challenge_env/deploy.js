// npx hardhat run --network localhost scripts/deploy.js

async function main() {
  // We get the contract to deploy
  const Exploit = await ethers.getContractFactory("Exploit");
  const e = await Exploit.deploy();

  await e.deployed();

  console.log("Greeter deployed to:", e.address);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });

// curl -X POST http://35.160.101.66:1337/get_instance?challenge_number=2